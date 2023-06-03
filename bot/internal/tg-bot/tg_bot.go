package tg_bot

import (
	"log"
	"net/http"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/l-orlov/slim-fairy/bot/internal/store"
	"github.com/l-orlov/slim-fairy/bot/internal/tg-bot/logic-handlers"
	"github.com/pkg/errors"
)

/* TODO:
делаем идеально функционал с чат гпт.
запрос меню на один день.
должно все хорошо отрабатывать и отдаваться ссылка на сбермаркет.

делаем кнопку: получить диету от диетолога
идет живое общение с диетологом.
скажем, что этот функционал в проработке и разработке.

добавить вопросы при регистрации.

убрать reply, а просто обычные сообщения отправлять

удалить бд, чтобы можно было заново регу пройти.
*/

type (
	Bot struct {
		bot     *gotgbot.Bot
		updater *ext.Updater
	}
)

func New(
	token string,
	aiClient logic_handlers.AIClient,
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
	logicHandlers := logic_handlers.New(aiClient, storage)

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
	// command to register user
	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("register", logicHandlers.StartUserRegistration)},
		map[string][]ext.Handler{
			logic_handlers.RegisterName:   {handlers.NewMessage(noCommands, logicHandlers.RegisterUserName)},
			logic_handlers.RegisterAge:    {handlers.NewMessage(noCommands, logicHandlers.RegisterUserAge)},
			logic_handlers.RegisterWeight: {handlers.NewMessage(noCommands, logicHandlers.RegisterUserWeight)},
			logic_handlers.RegisterHeight: {handlers.NewMessage(noCommands, logicHandlers.RegisterUserHeight)},
			logic_handlers.RegisterGender: {
				handlers.NewCallback(callbackquery.Equal(logic_handlers.RegisterGenderCbMale), logicHandlers.RegisterUserGenderMale),
				handlers.NewCallback(callbackquery.Equal(logic_handlers.RegisterGenderCbFemale), logicHandlers.RegisterUserGenderFemale),
			},
			logic_handlers.RegisterPhysicalActivity: {
				handlers.NewCallback(callbackquery.Equal(logic_handlers.RegisterPhysicalActivityCbLow), logicHandlers.RegisterUserPhysicalActivityLow),
				handlers.NewCallback(callbackquery.Equal(logic_handlers.RegisterPhysicalActivityCbMedium), logicHandlers.RegisterUserPhysicalActivityMedium),
				handlers.NewCallback(callbackquery.Equal(logic_handlers.RegisterPhysicalActivityCbHigh), logicHandlers.RegisterUserPhysicalActivityHigh),
			},
			logic_handlers.RegisterConfirm: {
				handlers.NewCallback(callbackquery.Equal(logic_handlers.RegisterConfirmCbYes), logicHandlers.ConfirmUserRegistrationYes),
				handlers.NewCallback(callbackquery.Equal(logic_handlers.RegisterConfirmCbNo), logicHandlers.ConfirmUserRegistrationNo),
			},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", logicHandlers.CancelUserRegistration)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
		},
	))

	// command to get diet menu from AI
	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("getdietfromai", logicHandlers.StartGettingDiet)},
		map[string][]ext.Handler{
			logic_handlers.GetDietFromAISelectMeals: {
				handlers.NewCallback(callbackquery.Equal(logic_handlers.GetDietFromAICbSelectMeals2), logicHandlers.SelectMeals2),
				handlers.NewCallback(callbackquery.Equal(logic_handlers.GetDietFromAICbSelectMeals3), logicHandlers.SelectMeals3),
			},
			logic_handlers.GetDietFromAISelectSnacks: {
				handlers.NewCallback(callbackquery.Equal(logic_handlers.GetDietFromAICbSelectSnacks0), logicHandlers.SelectSnacks0),
				handlers.NewCallback(callbackquery.Equal(logic_handlers.GetDietFromAICbSelectSnacks1), logicHandlers.SelectSnacks1),
				handlers.NewCallback(callbackquery.Equal(logic_handlers.GetDietFromAICbSelectSnacks2), logicHandlers.SelectSnacks2),
			},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", logicHandlers.CancelGettingDiet)},
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
