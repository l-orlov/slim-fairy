package logic_handlers

import (
	"io"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/pkg/errors"
)

var (
	defaultSendMessageOpts = &gotgbot.SendMessageOpts{
		ParseMode: "html",
	}
)

func replyWithStatesInConversation(
	b *gotgbot.Bot, ctx *ext.Context,
	msg string, success bool,
	failureStateName, successStateName string,
	opts *gotgbot.SendMessageOpts,
) (nextState error) {
	// If not success -> try again
	if !success {
		return replyInConversation(b, ctx, msg, failureStateName, opts)
	}

	// If success -> go to next state
	return replyInConversation(b, ctx, msg, successStateName, opts)
}

func replyInConversation(
	b *gotgbot.Bot, ctx *ext.Context,
	msg string, nextStateName string,
	opts *gotgbot.SendMessageOpts,
) (nextState error) {
	reply(b, ctx, msg, opts)
	return handlers.NextConversationState(nextStateName)
}

func endConversation(b *gotgbot.Bot, ctx *ext.Context, msg string) (nextState error) {
	reply(b, ctx, msg, nil)
	return handlers.EndConversation()
}

func reply(b *gotgbot.Bot, ctx *ext.Context, msg string, opts *gotgbot.SendMessageOpts) {
	sendOpts := defaultSendMessageOpts
	if opts != nil {
		sendOpts = opts
	}

	_, err := ctx.EffectiveMessage.Reply(b, msg, sendOpts)
	if err != nil {
		log.Printf("ctx.EffectiveMessage.Reply: %v", err)

		// Fallback with usual sending msg
		sendMessage(b, ctx, msg, opts)
	}
}

func sendMessage(b *gotgbot.Bot, ctx *ext.Context, msg string, opts *gotgbot.SendMessageOpts) {
	sendOpts := defaultSendMessageOpts
	if opts != nil {
		sendOpts = opts
	}

	_, err := b.SendMessage(ctx.EffectiveChat.Id, msg, sendOpts)
	if err != nil {
		log.Printf("b.SendMessage: %v", err)
	}
}

func sendDocument(
	b *gotgbot.Bot, ctx *ext.Context,
	file io.Reader, fileName, caption string,
) error {
	_, err := b.SendDocument(ctx.EffectiveChat.Id, gotgbot.NamedFile{
		File:     file,
		FileName: fileName,
	}, &gotgbot.SendDocumentOpts{
		Caption:          caption,
		ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	})
	if err != nil {
		return errors.Wrap(err, "b.SendDocument")
	}

	return nil
}
