package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/amiosamu/vk-internship/config"
	"github.com/amiosamu/vk-internship/internal/repo"
	"github.com/amiosamu/vk-internship/internal/service"
	"github.com/amiosamu/vk-internship/pkg/httpserver"
	"github.com/amiosamu/vk-internship/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/amiosamu/vk-internship/internal/controller/http/v1"
)

func Run(confPath string) {

	cfg, err := config.NewConfig(confPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	SetLogrus(cfg.Log.Level)

	log.Info("Init database...")

	db, err := postgres.InitDB()
	if err != nil {
		log.Fatalf("unable to init database: %v\n", err)
	}

	defer db.Close()

	log.Info("Init repositories...")

	repository := repo.NewRepos(db.DB)

	log.Info("Init services...")
	dependencies := service.ServiceDependencies{
		Repos: repository,

		Signkey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}

	services := service.NewServices(dependencies)

	log.Info("Init handlers and routes...")

	handler := gin.New()

	v1.NewRouter(handler, services)

	log.Info("Starting http server...")
	log.Debugf("Servce port: %s", cfg.HTTP.Port)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))


	log.Info("Configuring graceful shutdown...")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)


	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	log.Info("Shutting down...")

	err = httpServer.Shutdown()

	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
