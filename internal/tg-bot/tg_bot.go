package tg_bot

import (
	"log"
	"net/http"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/l-orlov/slim-fairy/internal/store"
	lhandlers "github.com/l-orlov/slim-fairy/internal/tg-bot/logic-handlers"
	"github.com/pkg/errors"
)

/* TODO:
- добавить цепочку обновления данных
*/

type (
	Bot struct {
		bot     *gotgbot.Bot
		updater *ext.Updater
	}
)

func New(
	token string,
	aiClient lhandlers.AIClient,
	storage *store.Storage,
) (*Bot, error) {
	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "gotgbot.NewBot")
	}

	// Create logic handlers
	logicHandlers := lhandlers.New(aiClient, storage)

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})
	dispatcher := updater.Dispatcher

	// command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", logicHandlers.Start))
	// command to get diet menu from AI
	dispatcher.AddHandler(handlers.NewCommand("getdietfromai", logicHandlers.GetDietFromAI))
	// command to register user
	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("register", logicHandlers.StartUserRegistration)},
		map[string][]ext.Handler{
			lhandlers.RegisterName:             {handlers.NewMessage(noCommands, logicHandlers.RegisterUserName)},
			lhandlers.RegisterAge:              {handlers.NewMessage(noCommands, logicHandlers.RegisterUserAge)},
			lhandlers.RegisterWeight:           {handlers.NewMessage(noCommands, logicHandlers.RegisterUserWeight)},
			lhandlers.RegisterHeight:           {handlers.NewMessage(noCommands, logicHandlers.RegisterUserHeight)},
			lhandlers.RegisterGender:           {handlers.NewMessage(noCommands, logicHandlers.RegisterUserGender)},
			lhandlers.RegisterPhysicalActivity: {handlers.NewMessage(noCommands, logicHandlers.RegisterUserPhysicalActivity)},
			lhandlers.RegisterConfirm:          {handlers.NewMessage(noCommands, logicHandlers.ConfirmUserRegistration)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancelreg", logicHandlers.CancelUserRegistration)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
		},
	))

	return &Bot{
		bot:     b,
		updater: updater,
	}, nil
}

func (b *Bot) Run() {
	// Start receiving updates.
	err := b.updater.StartPolling(b.bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.bot.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	b.updater.Idle()
}

// Create a matcher which only matches text which is not a command
func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}
