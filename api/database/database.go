package database

import (
	"booksapi/config"
	"booksapi/logger"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init() {
	db := config.GetAppsettings().Database
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", db.User, db.Pass, db.Host, db.Port, db.Db)

	var err error
	Pool, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		logger.Info(fmt.Sprintf("ERROR connecting pgx -> %s", err.Error()))
		os.Exit(1)
	}

	logger.Info("database pool initialized")
}

func Close() {
	Pool.Close()
}
