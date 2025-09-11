package app

import (
	"database/sql"

	"go.uber.org/zap"
)

var (
	Sugar   zap.SugaredLogger
	DB      *sql.DB
	ConnStr string
)
