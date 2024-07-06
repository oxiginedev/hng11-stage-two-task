package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oxiginedev/hng11-stage-two-task/api"
	"github.com/oxiginedev/hng11-stage-two-task/config"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore/sqlstore"
)

func Start() error {
	var envFile string

	flag.StringVar(&envFile, "env", ".env", "Path to .env file")
	flag.Parse()

	cfg, err := config.Load(envFile)
	if err != nil {
		return err
	}

	// initialize database
	db, err := sqlstore.New(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		return err
	}

	defer db.Close()

	app := api.New(db, cfg)

	s := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           app.Routes(),
		ReadHeaderTimeout: time.Second * 2,
		ReadTimeout:       time.Second * 15,
		WriteTimeout:      time.Second * 15,
	}

	log.Printf("Server started and listening on port :%d", cfg.Port)

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("cmd: server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// select {
	// case <-quit:

	// }

	if err := s.Shutdown(ctx); err != nil {
		return fmt.Errorf("cmd: server failed to shutdown")
	}

	return nil
}
