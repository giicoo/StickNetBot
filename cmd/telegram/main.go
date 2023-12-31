package main

import (
	"time"

	"github.com/giicoo/StickAIBot/config"
	fsmService "github.com/giicoo/StickAIBot/internal/fsm_service"
	resizeService "github.com/giicoo/StickAIBot/internal/resize_service"
	telegramService "github.com/giicoo/StickAIBot/internal/telegram_service"
	"github.com/giicoo/StickAIBot/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// init logger
	log := logger.InitLogger()

	// init env var
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("load .env: %v", err)
	}

	// init config
	cfg, err := config.InitConfig("./config/config.yaml")
	if err != nil {
		log.Panicf("config init: %v", err)
	}

	// init fsm
	rootFSM := fsmService.NewFsmService()

	// init resize
	resize := resizeService.NewResizeService(1 * time.Minute)
	// create bot
	api, err := telegramService.CreateBot(log, cfg, rootFSM, resize)
	if err != nil {
		log.Panicf("create bot: %v", err)
	}

	// init handlers -> start polling
	api.Start()
}
