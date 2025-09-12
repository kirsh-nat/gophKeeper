package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/shlex"
)

// ------------------- Version Info -------------------
var (
	Version   = "dev"
	BuildDate = "unknown"
)

// ------------------- Token Storage -------------------
type Config struct {
	Token string `json:"token"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gophkeer-cli", "config.json"), nil
}

func saveToken(token string) error {
	cfg := Config{Token: token}
	b, _ := json.MarshalIndent(cfg, "", "  ")

	path, err := configPath()
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(path), 0700)
	return os.WriteFile(path, b, 0600)
}

func loadToken() (string, error) {
	path, err := configPath()
	if err != nil {
		return "", err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return "", err
	}

	return cfg.Token, nil
}

// ------------------- HTTP Client -------------------
type Client struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// ------------------- API Methods -------------------
func (c *Client) Register(login, password string) error {
	body := map[string]string{"login": login, "password": password}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", c.BaseURL+"/api/user/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("register failed: %s %s", resp.Status, string(b))
	}

	fmt.Println("Registered:", login)
	return nil
}

func (c *Client) Login(login, password string) error {
	body := map[string]string{"login": login, "password": password}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", c.BaseURL+"/api/user/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	// читаем токен из Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			c.Token = cookie.Value
		}
	}

	if c.Token == "" {
		return fmt.Errorf("no token received")
	}

	fmt.Println("Login OK, token saved")
	return saveToken(c.Token)
}

func (c *Client) IsTokenValid() bool {
	if c.Token == "" {
		return false
	}
	// Используем простой GET для проверки токена
	req, _ := http.NewRequest("GET", c.BaseURL+"/ping", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: c.Token})
	resp, err := c.Client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (c *Client) CreateInfoItem(dataType, info, loginStr, passwordStr string) (int, error) {
	body := map[string]interface{}{
		"data_type": dataType,
		"data": map[string]string{
			"info":     info,
			"login":    loginStr,
			"password": passwordStr,
		},
	}

	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", c.BaseURL+"/createItem", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: c.Token})

	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// читаем тело как текст
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("createItem failed: %s %s", resp.Status, string(respBody))
	}

	// пробуем распарсить JSON с item_id
	var res struct {
		ItemID int `json:"item_id"`
	}
	if err := json.Unmarshal(respBody, &res); err != nil {
		// если не получилось, просто выводим тело
		fmt.Println("Response:", string(respBody))
		return 0, nil
	}

	return res.ItemID, nil
}

func (c *Client) UploadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}
	writer.Close()

	req, _ := http.NewRequest("POST", c.BaseURL+"/uploadFile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "token", Value: c.Token})

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println("Upload response:", string(b))

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("upload failed: %s", resp.Status)
	}
	return nil
}

func (c *Client) DownloadFile(itemID int, dest string) error {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/downloadFile?item_id=%d", c.BaseURL, itemID), nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: c.Token})

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed: %s %s", resp.Status, string(b))
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("File downloaded to:", dest)
	return nil
}

// ------------------- MAIN REPL -------------------
func main() {
	client := NewClient("http://localhost:8080")

	// Автологин по токену (как раньше)
	token, _ := loadToken()
	client.Token = token
	if !client.IsTokenValid() {
		fmt.Println("Token expired or missing. Please login.")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Username: ")
		user, _ := reader.ReadString('\n')
		user = strings.TrimSpace(user)

		fmt.Print("Password: ")
		pass, _ := reader.ReadString('\n')
		pass = strings.TrimSpace(pass)

		if err := client.Login(user, pass); err != nil {
			fmt.Println("Login failed:", err)
			return
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to Gophkeer CLI! Type 'help' for commands, 'exit' to quit.")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		// Разбор строки с поддержкой кавычек
		args, err := shlex.Split(line)
		if err != nil {
			fmt.Println("Invalid input:", err)
			continue
		}
		if len(args) == 0 {
			continue
		}

		cmd := args[0]

		switch cmd {
		case "exit", "quit":
			fmt.Println("Bye!")
			return
		case "help":
			fmt.Println("Available commands: register, login, create, upload, download, version, exit")
		case "version":
			fmt.Printf("Version: %s\nBuild Date: %s\n", Version, BuildDate)
		case "register":
			if len(args) < 3 {
				fmt.Println("Usage: register <username> <password>")
				continue
			}
			if err := client.Register(args[1], args[2]); err != nil {
				fmt.Println("Register error:", err)
			}
		case "login":
			if len(args) < 3 {
				fmt.Println("Usage: login <username> <password>")
				continue
			}
			if err := client.Login(args[1], args[2]); err != nil {
				fmt.Println("Login error:", err)
			}
		case "create":
			if len(args) < 5 {
				fmt.Println("Usage: create <data_type> <info> <login> <password>")
				continue
			}
			id, err := client.CreateInfoItem(args[1], args[2], args[3], args[4])
			if err != nil {
				fmt.Println("Create error:", err)
				continue
			}
			fmt.Println("Created item ID:", id)
		case "upload":
			if len(args) < 2 {
				fmt.Println("Usage: upload <file_path>")
				continue
			}
			if err := client.UploadFile(args[1]); err != nil {
				fmt.Println("Upload error:", err)
			}
		case "download":
			if len(args) < 3 {
				fmt.Println("Usage: download <item_id> <dest_file>")
				continue
			}
			itemID := 0
			fmt.Sscanf(args[1], "%d", &itemID)
			if err := client.DownloadFile(itemID, args[2]); err != nil {
				fmt.Println("Download error:", err)
			}
		default:
			fmt.Println("Unknown command:", cmd)
		}
	}
}
