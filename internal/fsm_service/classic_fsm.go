package fsmService

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
)

type ctxKey struct{}

var classicFSM = fsm.NewFSM(
	"classicStick",
	fsm.Events{
		{
			Name: "reset",
			Src:  []string{"title_set", "photo_set", "emoji_set", "check_end"},
			Dst:  "classicStick",
		},
		{
			Name: "title",
			Src:  []string{"classicStick"},
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
			Src:  []string{"emoji_set"},
			Dst:  "title_set",
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
	},
	fsm.Callbacks{
		"title_set": func(ctx context.Context, e *fsm.Event) {
			e.FSM.SetMetadata("title", ctx.Value(ctxKey{}))
			fmt.Printf("Set title %v", ctx.Value(ctxKey{}))
		},
		"photo_set": func(ctx context.Context, e *fsm.Event) {
			e.FSM.SetMetadata("photo", ctx.Value(ctxKey{}))
			fmt.Printf("Set photo %v", ctx.Value(ctxKey{}))
		},
		"emoji_set": func(ctx context.Context, e *fsm.Event) {
			e.FSM.SetMetadata("emoji", ctx.Value(ctxKey{}))
			fmt.Printf("Set emoji %v", ctx.Value(ctxKey{}))
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
