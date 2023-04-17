package main

import (
	"database/sql"
	"go.uber.org/zap"
	"log"
)

func initLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func initDatabase(dsn string, logger *zap.Logger) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err = db.Ping(); err != nil {
		logger.Fatal(err.Error())
	}
	return db
}
