package fsmService

import "github.com/looplab/fsm"

type FsmService struct {
	FsmMap map[int64]*fsm.FSM
}

func NewFsmService() *FsmService {
	return &FsmService{}
}

func (f *FsmService) NewClassicFsm(chat_id int64) {
	fs := newClassicFsm()
	f.FsmMap[chat_id] = fs
}

func (f *FsmService) DeleteFsm(chat_id int64) {
	delete(f.FsmMap, chat_id)
}
