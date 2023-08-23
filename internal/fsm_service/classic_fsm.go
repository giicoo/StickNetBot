package fsmService

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
)

type ctxKey struct{}

type classicFsm *fsm.FSM

func newClassicFsm() *fsm.FSM {
	fs := fsm.NewFSM(
		"start",
		fsm.Events{
			{
				Name: "title",
				Src:  []string{"start"},
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
				Name: "end",
				Src:  []string{"emoji_set"},
				Dst:  "check_end",
			},
			{
				Name: "create",
				Src:  []string{"check_end"},
				Dst:  "create_done",
			},
			{
				Name: "more",
				Src:  []string{"title_set"},
				Dst:  "photo_set",
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
			"create_done": func(ctx context.Context, e *fsm.Event) {
				title, ok := e.FSM.Metadata("title")
				if !ok {
					fmt.Println("Dont have title")
				}
				photo, ok := e.FSM.Metadata("photo")
				if !ok {
					fmt.Println("Dont have photo")
				}
				emoji, ok := e.FSM.Metadata("emoji")
				if !ok {
					fmt.Println("Dont have emoji")
				}
				fmt.Printf("Title: %v\nPhoto: %v\nEmoji: %v\n", title, photo, emoji)
			},
		},
	)
	return fs
}
