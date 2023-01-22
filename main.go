package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/router"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func SetupDB(db *sql.DB) {
	createTb := `
	CREATE TABLE IF NOT EXISTS pockets (id SERIAL PRIMARY KEY, name TEXT, category TEXT, amount TEXT, goal TEXT, currency TEXT, createdAt TIMESTAMP, updatedAt TIMESTAMP, deletedAt TIMESTAMP);
	CREATE TABLE IF NOT EXISTS transactions (id SERIAL PRIMARY KEY, type TEXT, status TEXT, sourcePocketId INT, destinationPocketID INT, description TEXT, amount TEXT, currency TEXT, createdAt TIMESTAMP);
	`
	_, err := db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}
}

func SeedData(db *sql.DB) {
	insertTb := `INSERT INTO transactions (type, status, sourcePocketId, destinationPocketID, description, amount, currency, createdAt) SELECT 'deposit', 'success', 1, 2, '', '10.00', 'THB', '2021-09-01T00:00:00Z' WHERE NOT EXISTS (SELECT * FROM transactions)`

	_, err := db.Exec(insertTb)

	if err != nil {
		log.Fatal("can't insert table", err)
	}
}

func main() {
	cfg := config.New().All()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		logger.Fatal("unable to configure database", zap.Error(err))
	}

	SetupDB(sql)
	SeedData(sql)

	e := router.RegRoute(cfg, logger, sql)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Hostname, cfg.Server.Port)

	go func() {
		err := e.Start(addr)
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("unexpected shutdown the server", zap.Error(err))
		}
		logger.Info("gracefully shutdown the server")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	gCtx := context.Background()
	ctx, cancel := context.WithTimeout(gCtx, 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("unexpected shutdown the server", zap.Error(err))
	}
}
