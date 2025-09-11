package services

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// NormalizeFilename приводит имя файла к безопасному виду
func NormalizeFilename(original string) string {
	// Оставляем только имя файла (без путей)
	base := filepath.Base(original)

	// Заменяем пробелы на "_"
	base = strings.ReplaceAll(base, " ", "_")

	// Регулярка: оставляем только буквы/цифры/._-
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	base = re.ReplaceAllString(base, "_")

	// Ограничиваем длину (например, 100 символов)
	if len(base) > 100 {
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		if len(name) > 90 {
			name = name[:90]
		}
		base = name + ext
	}

	// Добавляем UUID, чтобы избежать коллизий
	return uuid.New().String() + "_" + base
}
