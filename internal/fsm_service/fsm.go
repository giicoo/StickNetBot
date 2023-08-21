package fsmService

import "github.com/looplab/fsm"

type FsmService struct {
	ClassicFSM *fsm.FSM
}

func NewFsmService() *FsmService {
	return &FsmService{
		ClassicFSM: classicFSM,
	}
}
