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
	"github.com/l-orlov/slim-fairy/internal/store"
	lhandlers "github.com/l-orlov/slim-fairy/internal/tg-bot/logic-handlers"
	"github.com/pkg/errors"
)

/* TODO:
- добавить цепочку обновления данных (Добавить метод обновления параметров)

В регистрацию добавить параметры.

Добавить команду: получить меню от диетолога.

Добавить расчет калорий на бэке.
Чтобы в запрос в чат гпт подставилось количество калорий.
Чтобы пользователю отдавалось адекватное меню.

Сделать подпись: получить меню быстро от ИИ
Уведомить, что меню носит рекомендательный характер.
Важно, что это больше поиграть, чтобы не было презентаций.

Функция: получить меню от диетолога.
Название: получить консультацию от диетолога.

Диетологу должно прийти сообщение, что пользователь с ним связался.
Нужно, чтобы ответ от диетолога был профессиональный.
И нужно будет дольше ждать меню от диетолога (2 часа - день).

Диетолог включается и начинает общаться.

После получения меню:
кнопка "сформировать список продуктов для заказа в магазине".
отдает, готовый список для заказа продуктов

появляется кнопка "заказать продукты", которая ведет на сбермаркет.
потому что пока что у нас нет возможности реализовать это.

Спросить у Кати и Даниэля, как они сделали кнопки в телеграм боте.
*/

/* TODO:
делаем идеально функционал с чат гпт.
запрос меню на один день.
должно все хорошо отрабатывать и отдаваться ссылка на сбермаркет.

делаем кнопку: получить диету от диетолога
идет живое общение с диетологом.
скажем, что этот функционал в проработке и разработке.
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
	// command to register user
	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("register", logicHandlers.StartUserRegistration)},
		map[string][]ext.Handler{
			lhandlers.RegisterName:   {handlers.NewMessage(noCommands, logicHandlers.RegisterUserName)},
			lhandlers.RegisterAge:    {handlers.NewMessage(noCommands, logicHandlers.RegisterUserAge)},
			lhandlers.RegisterWeight: {handlers.NewMessage(noCommands, logicHandlers.RegisterUserWeight)},
			lhandlers.RegisterHeight: {handlers.NewMessage(noCommands, logicHandlers.RegisterUserHeight)},
			lhandlers.RegisterGender: {
				handlers.NewCallback(callbackquery.Equal(lhandlers.RegisterGenderCbMale), logicHandlers.RegisterUserGenderMale),
				handlers.NewCallback(callbackquery.Equal(lhandlers.RegisterGenderCbFemale), logicHandlers.RegisterUserGenderFemale),
			},
			lhandlers.RegisterPhysicalActivity: {
				handlers.NewCallback(callbackquery.Equal(lhandlers.RegisterPhysicalActivityCbLow), logicHandlers.RegisterUserPhysicalActivityLow),
				handlers.NewCallback(callbackquery.Equal(lhandlers.RegisterPhysicalActivityCbMedium), logicHandlers.RegisterUserPhysicalActivityMedium),
				handlers.NewCallback(callbackquery.Equal(lhandlers.RegisterPhysicalActivityCbHigh), logicHandlers.RegisterUserPhysicalActivityHigh),
			},
			lhandlers.RegisterConfirm: {
				handlers.NewCallback(callbackquery.Equal(lhandlers.RegisterConfirmCbYes), logicHandlers.ConfirmUserRegistrationYes),
				handlers.NewCallback(callbackquery.Equal(lhandlers.RegisterConfirmCbNo), logicHandlers.ConfirmUserRegistrationNo),
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
			lhandlers.GetDietFromAISelectMeals: {
				handlers.NewCallback(callbackquery.Equal(lhandlers.GetDietFromAICbSelectMeals2), logicHandlers.SelectMeals2),
				handlers.NewCallback(callbackquery.Equal(lhandlers.GetDietFromAICbSelectMeals3), logicHandlers.SelectMeals3),
			},
			lhandlers.GetDietFromAISelectSnacks: {
				handlers.NewCallback(callbackquery.Equal(lhandlers.GetDietFromAICbSelectSnacks0), logicHandlers.SelectSnacks0),
				handlers.NewCallback(callbackquery.Equal(lhandlers.GetDietFromAICbSelectSnacks1), logicHandlers.SelectSnacks1),
				handlers.NewCallback(callbackquery.Equal(lhandlers.GetDietFromAICbSelectSnacks2), logicHandlers.SelectSnacks2),
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
