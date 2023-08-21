package handlers

import (
	"github.com/giicoo/StickAIBot/config"
	"github.com/looplab/fsm"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func StartClassicFsm(cfg config.Config, log *logrus.Logger, classicFsm *fsm.FSM) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		callback := update.CallbackQuery

		callback_id := callback.ID

		chat_id := tu.ID(callback.From.ID)
		msg_text := cfg.ClassicFsmGroup.Start.Phrase + "\n" + classicFsm.Current()

		msg := tu.Message(chat_id, msg_text)
		msg.ParseMode = telego.ModeHTML

		_, err := bot.SendMessage(msg)
		if err != nil {
			log.Errorf("send message to %v chat: %v", chat_id, err)
		}

		call := tu.CallbackQuery(callback_id)

		err = bot.AnswerCallbackQuery(call)
		if err != nil {
			log.Errorf("send answer callback to %v callback: %v", callback_id, err)
		}
	}
}
