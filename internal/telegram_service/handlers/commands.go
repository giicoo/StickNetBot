package handlers

import (
	"github.com/giicoo/StickAIBot/config"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

// Command Start with "Welcome phrase"; inline keyboard from config
func Start(cfg config.Config, log *logrus.Logger) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {

		chat_id := tu.ID(update.Message.Chat.ID)
		msg_text := cfg.WelcomeGroup.Phrase

		inlineKeyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cfg.WelcomeGroup.InlineKeyboard.Row1.Btn1).WithCallbackData("create_sticker_pack"),
			),
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cfg.WelcomeGroup.InlineKeyboard.Row2.Btn1).WithCallbackData("manage_sticker_packs"),
			),
		)

		msg := tu.Message(chat_id, msg_text).WithReplyMarkup(inlineKeyboard)
		msg.ParseMode = telego.ModeHTML

		_, err := bot.SendMessage(msg)
		if err != nil {
			log.Errorf("Err send message to %v chat: %v", chat_id, err)
		}
	}
}
