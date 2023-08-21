package telegramService

import (
	"github.com/giicoo/StickAIBot/internal/telegram_service/handlers"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (api *APIBot) initHandlers() {
	// TODO: написать условия для хендлеров, что ловить состояние
	// TODO:  дописать фсм и рефакторинг небольшой из-за неправильного начала
	// command start
	api.bh.Handle(handlers.Start(api.cfg, api.log), th.CommandEqual("start"))

	api.bh.Handle(handlers.StartClassicFsm(api.cfg, api.log, api.fsm.ClassicFsm), th.CallbackDataEqual("create_sticker_pack"))
	// echo
	api.bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		msg := tu.Message(tu.ID(message.Chat.ID), message.Text)
		bot.SendMessage(msg)
	})
}
