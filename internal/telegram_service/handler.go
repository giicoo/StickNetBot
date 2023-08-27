package telegramService

import (
	"github.com/giicoo/StickAIBot/internal/telegram_service/handlers"
	th "github.com/mymmrac/telego/telegohandler"
)

func (api *APIBot) initHandlers() {
	// TODO: запусти докер наконец и собери там резайз сервис
	handler := handlers.NewHandlers(api.cfg, api.log, api.fsm)
	predicate := handlers.NewPredicateService(api.fsm)
	// ___ start group ____
	// command start
	api.bh.Handle(handler.CommandHandler.Start(), th.CommandEqual("start"))
	// cancel start fsm

	// ____________________

	// ___ classicFsm group ___
	api.bh.Handle(handler.ClassicFsmService.StartClassicFsm(), th.CallbackDataEqual("create_sticker_pack"))
	api.bh.Handle(handler.ClassicFsmService.MoreClassicFsmCallBack(), th.CallbackDataEqual("more_classicFsm"))
	api.bh.Handle(handler.ClassicFsmService.CancelStartCallBack(), th.CallbackDataEqual("cancel_start"))
	api.bh.Handle(handler.ClassicFsmService.CreateClassicFsm(), th.CallbackDataEqual("create_classicFsm"))
	api.bh.Handle(handler.ClassicFsmService.TitleClassicFsm(), predicate.ClassicFsmPredicate("start"), th.AnyMessage())
	api.bh.Handle(handler.ClassicFsmService.PhotoClassicFsm(api.resizeService), predicate.ClassicFsmPredicate("title_set"), th.AnyMessage())
	api.bh.Handle(handler.ClassicFsmService.EmojiClassicFsm(), predicate.ClassicFsmPredicate("photo_set"), th.AnyMessage())
	// ________________________
}
