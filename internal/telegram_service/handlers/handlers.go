package handlers

import (
	"github.com/giicoo/StickAIBot/config"
	fsmService "github.com/giicoo/StickAIBot/internal/fsm_service"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	cfg config.Config
	log *logrus.Logger

	fsm *fsmService.FsmService

	ClassicFsmService *ClassicFsmHandler
	CommandHandler    *CommandHandler
}

func NewHandlers(cfg config.Config, log *logrus.Logger, fsm *fsmService.FsmService) *Handlers {
	return &Handlers{
		ClassicFsmService: &ClassicFsmHandler{
			cfg: cfg,
			log: log,
			fsm: fsm,
		},
		CommandHandler: &CommandHandler{
			cfg: cfg,
			log: log,
			fsm: fsm,
		},
	}
}
