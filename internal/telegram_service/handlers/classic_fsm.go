package handlers

import (
	"context"
	"fmt"

	"github.com/giicoo/StickAIBot/config"
	fsmService "github.com/giicoo/StickAIBot/internal/fsm_service"
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
		return step_now.Current() == step
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
		log.Info("title")
		// init vars
		user_id := tu.ID(update.Message.From.ID)
		title := update.Message.Text
		msg_text := fmt.Sprintf(cfg.ClassicFsmGroup.Title.Phrase, title)

		// next fsm -> set_title
		fsm.FsmMap[user_id.ID].Event(context.WithValue(context.Background(), struct{}{}, title), "title")

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

func PhotoClassicFsm(cfg config.Config, log *logrus.Logger, fsm *fsmService.FsmService) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// init vars
		user_id := tu.ID(update.Message.From.ID)
		photo := update.Message.Text
		msg_text := fmt.Sprintf(cfg.ClassicFsmGroup.Photo.Phrase, photo)

		// next fsm -> set+photo
		fsm.FsmMap[user_id.ID].Event(context.WithValue(context.Background(), struct{}{}, photo), "photo")

		// message with inline_keyboard
		inline_keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cfg.ClassicFsmGroup.Photo.InlineKeyboard.Row1.Btn1).WithCallbackData("cancel_start"),
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

func EmojiClassicFsm(cfg config.Config, log *logrus.Logger, fsm *fsmService.FsmService) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		// init vars
		user_id := tu.ID(update.Message.From.ID)
		emoji := update.Message.Text

		// next fsm -> set_emoji
		fsm.FsmMap[user_id.ID].Event(context.WithValue(context.Background(), struct{}{}, emoji), "emoji")

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
		title, ok := fsm.FsmMap[user_id.ID].Metadata("title")
		if !ok {
			log.Errorf("err get title for classicFsm")
			title = "Не записалось, попробуйте позже"
		}

		msg_text := fmt.Sprintf(cfg.ClassicFsmGroup.Emoji.Phrase, title)
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
		log.Info("1")
		callback := update.CallbackQuery
		callback_id := callback.ID
		user_id := tu.ID(callback.From.ID)

		fsm.FsmMap[user_id.ID].Event(context.TODO(), "title")
		title, _ := fsm.FsmMap[user_id.ID].Metadata("title")
		msg_text := fmt.Sprintf(cfg.ClassicFsmGroup.Title.Phrase, title)

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
