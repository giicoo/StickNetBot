package handlers

import (
	"github.com/giicoo/StickAIBot/config"
	fsmService "github.com/giicoo/StickAIBot/internal/fsm_service"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

type CommandHandler struct {
	cfg config.Config
	log *logrus.Logger

	fsm *fsmService.FsmService
}

// Command Start with "Welcome phrase"; inline keyboard from config
func (h *CommandHandler) Start() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// vars
		user_id := tu.ID(update.Message.From.ID)

		// keyboard
		inlineKeyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(
					h.cfg.WelcomeGroup.InlineKeyboard.Row1.Btn1,
				).WithCallbackData("create_sticker_pack"),
			),
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(
					h.cfg.WelcomeGroup.InlineKeyboard.Row2.Btn1,
				).WithCallbackData("manage_sticker_packs"),
			),
		)

		// create msg
		msg_text := h.cfg.WelcomeGroup.Phrase
		msg := tu.Message(
			user_id,
			msg_text,
		).WithReplyMarkup(inlineKeyboard).WithParseMode("HTML")

		// send msg
		_, err := bot.SendMessage(msg)
		if err != nil {
			h.log.Errorf("err send message to %v user: %v", user_id, err)
		}
	}
}
