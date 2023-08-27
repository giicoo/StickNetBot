package handlers

import (
	fsmService "github.com/giicoo/StickAIBot/internal/fsm_service"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type PredicateService struct {
	fsm *fsmService.FsmService
}

func NewPredicateService(fsm *fsmService.FsmService) *PredicateService {
	return &PredicateService{
		fsm: fsm,
	}
}
func (p *PredicateService) ClassicFsmPredicate(step string) th.Predicate {
	return func(update telego.Update) bool {
		user_id := update.Message.From.ID
		step_now, ok := p.fsm.FsmMap[user_id]
		if !ok {
			return false
		}
		return step_now.Fsm.Current() == step
	}
}
