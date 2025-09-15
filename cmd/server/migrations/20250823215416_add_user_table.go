package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up, down)
}

func up(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		password_salt VARCHAR(255) NOT NULL,
	);	
`)
	if err != nil {
		return err
	}
	return nil
}

func down(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE users")
	if err != nil {
		return err
	}
	return nil
}

//TODO: авозможно добавить ключи шифрования (если шифруем на клиенте, что убдет оптимальным решением  !!! )
// -- Ключи шифрования (если храним на сервере в завёрнутом виде)
// create table user_keys (
//   user_id           uuid primary key references users(id) on delete cascade,
//   wrapped_master_key bytea not null,  -- зашифровано KMS/Vault или клиентом
//   kms_key_id        text,             -- метка ключа KMS, если используем
//   updated_at        timestamptz not null default now()
// );
