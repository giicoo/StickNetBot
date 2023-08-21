package telegramService

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/giicoo/StickAIBot/config"
	fsmService "github.com/giicoo/StickAIBot/internal/fsm_service"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/sirupsen/logrus"
)

type APIBot struct {
	log *logrus.Logger
	cfg config.Config

	b  *telego.Bot
	bh *th.BotHandler

	fsm *fsmService.FsmService

	stop chan struct{}
	done chan struct{}
}

func CreateBot(log *logrus.Logger, cfg config.Config, fsm *fsmService.FsmService) (*APIBot, error) {
	// initialize bot
	b, err := telego.NewBot(cfg.TOKEN, telego.WithDefaultDebugLogger())
	if err != nil {
		return nil, fmt.Errorf("init bot: %v", err)
	}

	// Initialize done and stop chan
	done := make(chan struct{}, 1)
	stop := make(chan struct{}, 1)

	// Get updates
	updates, err := b.UpdatesViaLongPolling(nil)
	if err != nil {
		return nil, fmt.Errorf("updates bot: %v", err)
	}

	// Create bot handler with stop timeout
	bh, err := th.NewBotHandler(b, updates, th.WithStopTimeout(time.Second*10))
	if err != nil {
		return nil, fmt.Errorf("handler bot: %v", err)
	}

	api := &APIBot{
		log: log,
		cfg: cfg,

		b:  b,
		bh: bh,

		fsm: fsm,

		stop: stop,
		done: done,
	}
	return api, nil
}

func (api *APIBot) Start() {
	// handle stop signal
	api.handleStopSignal()

	//init handlers
	api.initHandlers()

	// Initialize signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Handle Ctrl+C stop signal
	go func() {
		<-sigs
		api.Stop()
	}()

	// Start handling in goroutine
	go api.bh.Start()
	api.log.Info("Handling updates...")

	<-api.done
	api.log.Info("Done")
}

// Split "sigs" and "stop" channels, so that it is possible to safely terminate the bot under a different condition
func (api *APIBot) Stop() {
	go func() {
		api.stop <- struct{}{}
	}()
}

// Handle any stop signal
func (api *APIBot) handleStopSignal() {
	go func() {
		// Wait for stop signal
		<-api.stop

		api.log.Info("Stopping...")

		api.b.StopLongPolling()
		api.log.Info("Long polling done")

		api.bh.Stop()
		api.log.Info("Bot handler done")

		// Notify that stop is done
		api.done <- struct{}{}
	}()
}
