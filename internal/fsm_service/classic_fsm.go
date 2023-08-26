package fsmService

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
	"github.com/mymmrac/telego"
)

type sticker_pack struct {
	Title  string
	Sticks []telego.InputSticker
}

type classicFsm *fsm.FSM

func newClassicFsm() *fsm.FSM {
	fs := fsm.NewFSM(
		"start",
		fsm.Events{
			{
				Name: "title",
				Src:  []string{"start", "emoji_set"},
				Dst:  "title_set",
			},
			{
				Name: "photo",
				Src:  []string{"title_set"},
				Dst:  "photo_set",
			},
			{
				Name: "emoji",
				Src:  []string{"photo_set"},
				Dst:  "emoji_set",
			},
			{
				Name: "more",
				Src:  []string{"title_set"},
				Dst:  "more_set",
			},
		},
		fsm.Callbacks{
			"title": func(ctx context.Context, e *fsm.Event) {
				e.FSM.SetMetadata("title", ctx.Value(struct{}{}))
				fmt.Printf("\nSet title %v\n", ctx.Value(struct{}{}))
			},
			"photo": func(ctx context.Context, e *fsm.Event) {
				e.FSM.SetMetadata("photo", ctx.Value(struct{}{}))
				fmt.Printf("\nSet photo %v\n", ctx.Value(struct{}{}))
			},
			"emoji": func(ctx context.Context, e *fsm.Event) {
				e.FSM.SetMetadata("emoji", ctx.Value(struct{}{}))
				fmt.Printf("\nSet emoji %v\n", ctx.Value(struct{}{}))
			},
		},
	)
	return fs
}
