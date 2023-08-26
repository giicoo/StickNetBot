package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/giicoo/StickAIBot/config"
	fsmService "github.com/giicoo/StickAIBot/internal/fsm_service"
	resizeService "github.com/giicoo/StickAIBot/internal/resize_service"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func ClassicFsmPredicate(step string, fsm *fsmService.FsmService) th.Predicate {
	return func(update telego.Update) bool {
		user_id := update.Message.From.ID
		step_now, ok := fsm.FsmMap[user_id]
		if !ok {
			return false
		}
		return step_now.Fsm.Current() == step
	}
}
func StartClassicFsm(cfg config.Config, log *logrus.Logger, fsm *fsmService.FsmService) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		log.Info("start")

		// init vars
		callback_id := update.CallbackQuery.ID

		user_id := tu.ID(update.CallbackQuery.From.ID)
		msg_text := cfg.ClassicFsmGroup.Start.Phrase

		// init fsm
		fsm.NewClassicFsm(user_id.ID)

		// message with inline_keyboard
		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cfg.ClassicFsmGroup.Start.InlineKeyboard.Row1.Btn1).WithCallbackData("cancel_start"),
			),
		)

		msg := tu.Message(user_id, msg_text).WithReplyMarkup(inline_keyboard)
		msg.ParseMode = telego.ModeHTML

		_, err := bot.SendMessage(msg)
		if err != nil {
			log.Errorf("send message to %v chat: %v", user_id, err)
		}

		// answer callback query
		call := tu.CallbackQuery(callback_id)

		err = bot.AnswerCallbackQuery(call)
		if err != nil {
			log.Errorf("send answer callback to %v callback: %v", callback_id, err)
		}
	}
}

func TitleClassicFsm(cfg config.Config, log *logrus.Logger, fsm *fsmService.FsmService) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// init vars
		user_id := tu.ID(update.Message.From.ID)
		title := update.Message.Text
		msg_text := fmt.Sprintf(cfg.ClassicFsmGroup.Title.Phrase, title)

		// next fsm -> set_title
		fsm.FsmMap[user_id.ID].Fsm.Event(context.WithValue(context.Background(), struct{}{}, title), "title")
		fsm.FsmMap[user_id.ID].StickerPack.Title = title
		// message with inline_keyboard
		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cfg.ClassicFsmGroup.Title.InlineKeyboard.Row1.Btn1).WithCallbackData("cancel_start"),
			),
		)

		msg := tu.Message(user_id, msg_text).WithReplyMarkup(inline_keyboard)
		msg.ParseMode = telego.ModeHTML

		_, err := bot.SendMessage(msg)
		if err != nil {
			log.Errorf("send message to %v chat: %v", user_id, err)
		}

	}
}

func PhotoClassicFsm(cfg config.Config, log *logrus.Logger, fsm *fsmService.FsmService, resize *resizeService.ResizeService) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// init vars
		user_id := tu.ID(update.Message.From.ID)
		photo := update.Message.Document
		msg_text := fmt.Sprintf(cfg.ClassicFsmGroup.Photo.Phrase, photo)

		file, _ := bot.GetFile(&telego.GetFileParams{
			FileID: photo.FileID,
		})

		// Download file from Telegram using FileDownloadURL helper func to get full URL
		fileData, err := tu.DownloadFile(bot.FileDownloadURL(file.FilePath))
		err = resize.ResizeImage(fileData, fmt.Sprintf("%v_%v", user_id.ID, photo.FileID), photo.MimeType)
		if err != nil {
			log.Errorf("resize img: %v", err)
		}

		resize_photo, err := os.Open(fmt.Sprintf("%v/new_%v_%v.png", resize.Path, user_id.ID, photo.FileID))
		if err != nil {
			log.Errorf("open resize img: %v", err)
		}

		// message with inline_keyboard
		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cfg.ClassicFsmGroup.Photo.InlineKeyboard.Row1.Btn1).WithCallbackData("cancel_start"),
			),
		)

		document := tu.Document(
			user_id,
			tu.File(resize_photo),
		).WithCaption(msg_text).WithReplyMarkup(inline_keyboard)
		document.ParseMode = telego.ModeHTML

		msg, err := bot.SendDocument(document)
		if err != nil {
			log.Errorf("send document to %v chat: %v", user_id, err)
		}

		// next fsm -> set+photo
		fsm.FsmMap[user_id.ID].Fsm.Event(context.WithValue(context.Background(), struct{}{}, photo), "photo")
		fsm.FsmMap[user_id.ID].StickerPack.Sticks = append(fsm.FsmMap[user_id.ID].StickerPack.Sticks, telego.InputSticker{Sticker: telego.InputFile{FileID: msg.Document.FileID}})

	}
}

func EmojiClassicFsm(cfg config.Config, log *logrus.Logger, fsm *fsmService.FsmService) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// init vars
		user_id := tu.ID(update.Message.From.ID)
		emoji := update.Message.Text

		// next fsm -> set_emoji
		fsm.FsmMap[user_id.ID].Fsm.Event(context.WithValue(context.Background(), struct{}{}, emoji), "emoji")
		fsm.FsmMap[user_id.ID].StickerPack.Sticks[len(fsm.FsmMap[user_id.ID].StickerPack.Sticks)-1].EmojiList = append(fsm.FsmMap[user_id.ID].StickerPack.Sticks[len(fsm.FsmMap[user_id.ID].StickerPack.Sticks)-1].EmojiList, emoji)
		// message with inline_keyboard
		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cfg.ClassicFsmGroup.Emoji.InlineKeyboard.Row1.Btn1).WithCallbackData("more_classicFsm"),
				tu.InlineKeyboardButton(cfg.ClassicFsmGroup.Emoji.InlineKeyboard.Row1.Btn2).WithCallbackData("create_classicFsm"),
			),
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cfg.ClassicFsmGroup.Emoji.InlineKeyboard.Row2.Btn1).WithCallbackData("cancel_start"),
			),
		)

		msg_text := fmt.Sprintf(cfg.ClassicFsmGroup.Emoji.Phrase)
		msg := tu.Message(user_id, msg_text).WithReplyMarkup(inline_keyboard)
		msg.ParseMode = telego.ModeHTML

		_, err := bot.SendMessage(msg)
		if err != nil {
			log.Errorf("send message to %v chat: %v", user_id, err)
		}
	}
}

func MoreClassicFsm(cfg config.Config, log *logrus.Logger, fsm *fsmService.FsmService) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		callback := update.CallbackQuery
		callback_id := callback.ID
		user_id := tu.ID(callback.From.ID)

		log.Info(fsm.FsmMap[user_id.ID].Fsm.Current())
		fsm.FsmMap[user_id.ID].Fsm.Event(context.Background(), "title")
		log.Info(fsm.FsmMap[user_id.ID].Fsm.Current())
		log.Info(fsm.FsmMap[user_id.ID].StickerPack)
		msg_text := fmt.Sprintf(cfg.ClassicFsmGroup.Title.Phrase)

		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cfg.ClassicFsmGroup.Title.InlineKeyboard.Row1.Btn1).WithCallbackData("cancel_start"),
			),
		)

		msg := tu.Message(user_id, msg_text).WithReplyMarkup(inline_keyboard)
		msg.ParseMode = telego.ModeHTML

		_, err := bot.SendMessage(msg)
		if err != nil {
			log.Errorf("send message to %v chat: %v", user_id, err)
		}

		// answer callback query
		call := tu.CallbackQuery(callback_id)

		err = bot.AnswerCallbackQuery(call)
		if err != nil {
			log.Errorf("send answer callback to %v callback: %v", callback_id, err)
		}
	}
}

func CreateClassicFsm(cfg config.Config, log *logrus.Logger, fsm *fsmService.FsmService) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		user_id := tu.ID(update.CallbackQuery.From.ID)
		callback_id := update.CallbackQuery.ID

		stick_pack := fsm.FsmMap[user_id.ID].StickerPack

		_, _ = bot.GetMyName(&telego.GetMyNameParams{})
		log.Info(stick_pack.Sticks)
		err := bot.CreateNewStickerSet(&telego.CreateNewStickerSetParams{
			UserID:        user_id.ID,
			Title:         fmt.Sprintf("%v || @StickNetBot", stick_pack.Title),
			Name:          fmt.Sprintf("%v_%v_by_StickNetBot", stick_pack.Title, user_id.ID),
			Stickers:      stick_pack.Sticks,
			StickerFormat: "static",
		})

		if err != nil {
			log.Errorf("create sticker pack from %v: %v", user_id, err)
		}

		msg := tu.Message(user_id, fmt.Sprintf("https://t.me/addstickers/%v_%v_by_StickNetBot", stick_pack.Title, user_id.ID))
		msg.ParseMode = telego.ModeHTML

		_, err = bot.SendMessage(msg)
		if err != nil {
			log.Errorf("send message to %v chat: %v", user_id, err)
		}

		// answer callback query
		call := tu.CallbackQuery(callback_id)

		err = bot.AnswerCallbackQuery(call)
		if err != nil {
			log.Errorf("send answer callback to %v callback: %v", callback_id, err)
		}
	}
}
