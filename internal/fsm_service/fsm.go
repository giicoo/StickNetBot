package fsmService

import "github.com/looplab/fsm"

type FsmElement struct {
	Fsm         *fsm.FSM
	StickerPack *sticker_pack
}
type FsmService struct {
	FsmMap map[int64]*FsmElement
}

func NewFsmService() *FsmService {
	return &FsmService{
		FsmMap: map[int64]*FsmElement{},
	}
}

func (f *FsmService) NewClassicFsm(user_id int64) {
	fs := newClassicFsm()
	f.FsmMap[user_id] = &FsmElement{
		Fsm:         fs,
		StickerPack: &sticker_pack{},
	}
}

func (f *FsmService) DeleteFsm(chat_id int64) {
	delete(f.FsmMap, chat_id)
}
