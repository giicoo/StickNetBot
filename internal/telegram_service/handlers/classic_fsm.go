package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/giicoo/StickAIBot/config"
	fsmService "github.com/giicoo/StickAIBot/internal/fsm_service"
	resizeService "github.com/giicoo/StickAIBot/internal/resize_service"
	"github.com/lovelydeng/gomoji"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

type ClassicFsmHandler struct {
	cfg config.Config
	log *logrus.Logger

	fsm *fsmService.FsmService
}

func MessageError(user_id telego.ChatID, cfg config.Config, fsm *fsmService.FsmService, text string, delete_fsm bool) *telego.SendMessageParams {
	// inline_keyboard
	inline_keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				cfg.ClassicFsmGroup.Start.InlineKeyboard.Row1.Btn1,
			).WithCallbackData("cancel_start"),
		),
	)
	msg_text := fmt.Sprintf("Извините, произошла ошибка\nПопробуйте позже\n%v", text)
	msg := tu.Message(
		user_id,
		msg_text,
	).WithReplyMarkup(inline_keyboard).WithParseMode(telego.ModeHTML)
	if delete_fsm {
		// delete fsm not to lose memory
		fsm.DeleteFsm(user_id.ID)
	}
	return msg
}

// Start msg with create fsm
func (h *ClassicFsmHandler) StartClassicFsm() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// vars
		user_id := tu.ID(update.CallbackQuery.From.ID)
		callback_id := update.CallbackQuery.ID

		// inline_keyboard
		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(
					h.cfg.ClassicFsmGroup.Start.InlineKeyboard.Row1.Btn1,
				).WithCallbackData("cancel_start"),
			),
		)

		// create msg
		msg_text := h.cfg.ClassicFsmGroup.Start.Phrase
		msg := tu.Message(
			user_id,
			msg_text,
		).WithReplyMarkup(inline_keyboard).WithParseMode(telego.ModeHTML)

		// send msg
		_, err := bot.SendMessage(msg)
		if err != nil {
			h.log.Errorf("send message to %v chat: %v", user_id, err)
		}

		// init fsm
		h.fsm.NewClassicFsm(user_id.ID)

		// answer callback query
		callback := tu.CallbackQuery(callback_id)
		err = bot.AnswerCallbackQuery(callback)
		if err != nil {
			h.log.Errorf("send answer callback to %v callback: %v", callback_id, err)
		}
	}
}

// Start msg with delete fsm
func (h *ClassicFsmHandler) CancelStartCallBack() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// vars
		callback_id := update.CallbackQuery.ID
		user_id := tu.ID(update.CallbackQuery.From.ID)

		// start message
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
		).WithReplyMarkup(inlineKeyboard).WithParseMode(telego.ModeHTML)

		// send msg
		_, err := bot.SendMessage(msg)
		if err != nil {
			h.log.Errorf("Err send message to %v chat: %v", user_id, err)
		}

		// delete fsm
		h.fsm.DeleteFsm(user_id.ID)

		// answer callback
		callback := tu.CallbackQuery(callback_id)
		err = bot.AnswerCallbackQuery(callback)
		if err != nil {
			h.log.Errorf("answer callback to %v: %v", callback_id, err)
		}
	}
}

// Set title
func (h *ClassicFsmHandler) TitleClassicFsm() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// init vars
		user_id := tu.ID(update.Message.From.ID)
		title := update.Message.Text
		h.log.Info(update.Message.Text)
		// inline_keyboard
		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(h.cfg.ClassicFsmGroup.Title.InlineKeyboard.Row1.Btn1).WithCallbackData("cancel_start"),
			),
		)

		// create msg
		msg_text := fmt.Sprintf(h.cfg.ClassicFsmGroup.Title.Phrase, title)
		msg := tu.Message(
			user_id,
			msg_text,
		).WithReplyMarkup(inline_keyboard).WithParseMode(telego.ModeHTML)

		// check if not text and exit
		if title == "" {
			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "Напишите текстом название", false))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		// next fsm
		fsm_elm := h.fsm.FsmMap[user_id.ID]
		err := fsm_elm.Fsm.Event(context.WithValue(context.Background(), struct{}{}, title), "title")
		fsm_elm.StickerPack.Title = title

		if err != nil {
			h.log.Errorf("fsm err to title: %v", err)

			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "Или начать еще раз", true))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		// send msg
		_, err = bot.SendMessage(msg)
		if err != nil {
			h.log.Errorf("send message to %v user: %v", user_id, err)
		}

	}
}

// Download file -> resize file -> send new_file -> save id new_file
func (h *ClassicFsmHandler) PhotoClassicFsm(resize *resizeService.ResizeService) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// init vars
		user_id := tu.ID(update.Message.From.ID)
		photo := update.Message.Document

		if photo == nil || strings.Split(photo.MimeType, "/")[0] != "image" {

			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "<strong>Отправьте фото как файл</strong>", false))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		// Download file from Telegram using FileDownloadURL helper func to get full URL
		file, _ := bot.GetFile(&telego.GetFileParams{
			FileID: photo.FileID,
		})
		fileData, err := tu.DownloadFile(bot.FileDownloadURL(file.FilePath))
		if err != nil {
			h.log.Errorf("download file: %v", err)

			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "", true))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		err = resize.ResizeImage(fileData, fmt.Sprintf("%v_%v", user_id.ID, photo.FileID), photo.MimeType)
		if err != nil {
			h.log.Errorf("resize img: %v", err)

			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "", true))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		// resize photo
		resize_photo, err := os.Open(fmt.Sprintf("%v/new_%v_%v.png", resize.Path, user_id.ID, photo.FileID))
		if err != nil {
			h.log.Errorf("open resize img: %v", err)

			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "", true))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		// inline_keyboard
		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(
					h.cfg.ClassicFsmGroup.Photo.InlineKeyboard.Row1.Btn1,
				).WithCallbackData("cancel_start"),
			),
		)

		// create msg
		msg_text := fmt.Sprintf(h.cfg.ClassicFsmGroup.Photo.Phrase, photo)
		document := tu.Document(
			user_id,
			tu.File(resize_photo),
		).WithCaption(msg_text).WithReplyMarkup(inline_keyboard).WithParseMode(telego.ModeHTML)

		// send msg
		msg, err := bot.SendDocument(document)
		if err != nil {
			h.log.Errorf("send document to %v user: %v", user_id, err)
		}

		// next fsm
		new_file := tu.FileFromID(msg.Document.FileID)
		stick := telego.InputSticker{
			Sticker: new_file,
		}

		fsm_elm := h.fsm.FsmMap[user_id.ID]
		err = fsm_elm.Fsm.Event(context.WithValue(context.Background(), struct{}{}, photo), "photo")
		fsm_elm.StickerPack.Sticks = append(
			fsm_elm.StickerPack.Sticks,
			stick,
		)

		if err != nil {
			h.log.Errorf("fsm err to title: %v", err)

			// send msg about err
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "", true))
			if err != nil {
				h.log.Errorf("send msg to %v user: %v", user_id, err)
			}
		}
	}
}

func (h *ClassicFsmHandler) EmojiClassicFsm() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// init vars
		user_id := tu.ID(update.Message.From.ID)
		emoji := update.Message.Text
		_, check := gomoji.GetInfo(emoji)
		if emoji == "" || check == gomoji.ErrStrNotEmoji {
			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "<strong>Отправьте одно эмоджи</strong>", false))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		// message with inline_keyboard
		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(h.cfg.ClassicFsmGroup.Emoji.InlineKeyboard.Row1.Btn1).WithCallbackData("more_classicFsm"),
				tu.InlineKeyboardButton(h.cfg.ClassicFsmGroup.Emoji.InlineKeyboard.Row1.Btn2).WithCallbackData("create_classicFsm"),
			),
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(h.cfg.ClassicFsmGroup.Emoji.InlineKeyboard.Row2.Btn1).WithCallbackData("cancel_start"),
			),
		)

		msg_text := fmt.Sprintf(h.cfg.ClassicFsmGroup.Emoji.Phrase)
		msg := tu.Message(user_id, msg_text).WithReplyMarkup(inline_keyboard)
		msg.ParseMode = telego.ModeHTML

		// next fsm -> set_emoji
		fsm_elm := h.fsm.FsmMap[user_id.ID]
		err := fsm_elm.Fsm.Event(context.WithValue(context.Background(), struct{}{}, emoji), "emoji")
		fsm_elm.StickerPack.Sticks[len(fsm_elm.StickerPack.Sticks)-1].EmojiList = append(
			fsm_elm.StickerPack.Sticks[len(fsm_elm.StickerPack.Sticks)-1].EmojiList,
			emoji,
		)
		if err != nil {
			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "", true))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		_, err = bot.SendMessage(msg)
		if err != nil {
			h.log.Errorf("send message to %v chat: %v", user_id, err)
		}
	}
}

func (h *ClassicFsmHandler) MoreClassicFsmCallBack() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		callback_id := update.CallbackQuery.ID
		user_id := tu.ID(update.CallbackQuery.From.ID)

		fsm_elm := h.fsm.FsmMap[user_id.ID]
		err := fsm_elm.Fsm.Event(context.Background(), "title")
		if err != nil {
			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "", true))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(h.cfg.ClassicFsmGroup.Title.InlineKeyboard.Row1.Btn1).WithCallbackData("cancel_start"),
			),
		)

		msg_text := fmt.Sprintf(h.cfg.ClassicFsmGroup.Title.Phrase)
		msg := tu.Message(user_id, msg_text).WithReplyMarkup(inline_keyboard)
		msg.ParseMode = telego.ModeHTML

		_, err = bot.SendMessage(msg)
		if err != nil {
			h.log.Errorf("send message to %v chat: %v", user_id, err)
		}

		// answer callback query
		call := tu.CallbackQuery(callback_id)

		err = bot.AnswerCallbackQuery(call)
		if err != nil {
			h.log.Errorf("send answer callback to %v callback: %v", callback_id, err)
		}
	}
}

func (h *ClassicFsmHandler) CreateClassicFsm() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		callback_id := update.CallbackQuery.ID
		user_id := tu.ID(update.CallbackQuery.From.ID)

		stick_pack := h.fsm.FsmMap[user_id.ID].StickerPack

		err := bot.CreateNewStickerSet(&telego.CreateNewStickerSetParams{
			UserID:        user_id.ID,
			Title:         fmt.Sprintf("%v || @StickNetBot", stick_pack.Title),
			Name:          fmt.Sprintf("%v_%v_by_StickNetBot", stick_pack.Title, user_id.ID),
			Stickers:      stick_pack.Sticks,
			StickerFormat: "static",
		})

		if err != nil {
			h.log.Errorf("create sticker pack from %v: %v", user_id, err)

			// send msg
			_, err := bot.SendMessage(MessageError(user_id, h.cfg, h.fsm, "", true))
			if err != nil {
				h.log.Errorf("send message to %v user: %v", user_id, err)
			}
			return
		}

		msg := tu.Message(user_id, fmt.Sprintf("https://t.me/addstickers/%v_%v_by_StickNetBot", stick_pack.Title, user_id.ID))
		msg.ParseMode = telego.ModeHTML

		_, err = bot.SendMessage(msg)
		if err != nil {
			h.log.Errorf("send message to %v chat: %v", user_id, err)
		}

		// answer callback query
		call := tu.CallbackQuery(callback_id)

		err = bot.AnswerCallbackQuery(call)
		if err != nil {
			h.log.Errorf("send answer callback to %v callback: %v", callback_id, err)
		}
	}
}
