package telegramService

import (
	"github.com/giicoo/StickAIBot/internal/telegram_service/handlers"
	th "github.com/mymmrac/telego/telegohandler"
)

func (api *APIBot) initHandlers() {
	// TODO: запусти докер наконец и собери там резайз сервис

	// ___ start group ____
	// command start
	api.bh.Handle(handlers.Start(api.cfg, api.log), th.CommandEqual("start"))
	// cancel start fsm
	api.bh.Handle(handlers.CancelStart(api.cfg, api.log, api.fsm), th.CallbackDataEqual("cancel_start"))
	// ____________________

	// ___ classicFsm group ___
	api.bh.Handle(handlers.StartClassicFsm(api.cfg, api.log, api.fsm), th.CallbackDataEqual("create_sticker_pack"))
	api.bh.Handle(handlers.MoreClassicFsm(api.cfg, api.log, api.fsm), th.CallbackDataEqual("more_classicFsm"))
	api.bh.Handle(handlers.TitleClassicFsm(api.cfg, api.log, api.fsm), handlers.ClassicFsmPredicate("start", api.fsm), th.AnyMessage())
	api.bh.Handle(handlers.PhotoClassicFsm(api.cfg, api.log, api.fsm, api.resizeService), handlers.ClassicFsmPredicate("title_set", api.fsm), th.AnyMessage())
	api.bh.Handle(handlers.EmojiClassicFsm(api.cfg, api.log, api.fsm), handlers.ClassicFsmPredicate("photo_set", api.fsm), th.AnyMessage())
	// ________________________

}
