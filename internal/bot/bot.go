// internal/bot/bot.go
package bot

import (
	"log"
	"salary-bot/internal/bot/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramBot struct {
	client *tgbotapi.BotAPI
	state  *state.Manager
}

func NewBot(token string) (Bot, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	client.Debug = false
	log.Printf("Авторизован бот: @%s", client.Self.UserName)

	return &telegramBot{
		client: client,
		state:  state.NewManager(),
	}, nil
}

func (b *telegramBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.client.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		// Логируем
		log.Printf("[ChatID=%d] %s", chatID, text)

		// Обработка команды /start
		if text == "/start" {
			b.handleStart(chatID)
			continue
		}

		// Получаем текущее состояние пользователя
		userState := b.state.Get(chatID)

		switch userState.Step {
		case state.StepNone:
			// Если пользователь пишет что-то без /start — напомним
			b.client.Send(tgbotapi.NewMessage(chatID, "Нажмите /start, чтобы начать."))
		case state.StepAwaitingTech:
			b.handleTechSelection(chatID, text)
		case state.StepAwaitingExperience:
			b.handleExperienceSelection(chatID, text)
		}
	}
}

func (b *telegramBot) handleStart(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Хочешь узнать актуальные зарплаты в ИТ? Просто ответь на пару вопросов")
	b.client.Send(msg)

	b.state.Clear(chatID)

	// Гибкий список технологий — легко расширять
	techOptions := []string{"php", "go", "python"}

	var buttonRow []tgbotapi.KeyboardButton
	for _, tech := range techOptions {
		buttonRow = append(buttonRow, tgbotapi.NewKeyboardButton(tech))
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(buttonRow...),
	)
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true

	msg2 := tgbotapi.NewMessage(chatID, "1. Какой ваш технологический стек?")
	msg2.ReplyMarkup = keyboard

	b.state.Set(chatID, &state.UserState{Step: state.StepAwaitingTech})
	b.client.Send(msg2)
}

func (b *telegramBot) handleTechSelection(chatID int64, tech string) {
	validTech := map[string]bool{
		"php":    true,
		"go":     true,
		"python": true,
	}

	if !validTech[tech] {
		// Неверный выбор — повторим вопрос с кнопками
		techOptions := []string{"php", "go", "python"}
		var buttonRow []tgbotapi.KeyboardButton
		for _, t := range techOptions {
			buttonRow = append(buttonRow, tgbotapi.NewKeyboardButton(t))
		}

		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(buttonRow...),
		)
		keyboard.OneTimeKeyboard = true
		keyboard.ResizeKeyboard = true

		msg := tgbotapi.NewMessage(chatID, "Пожалуйста, выберите один из вариантов: php, go, python")
		msg.ReplyMarkup = keyboard
		b.client.Send(msg)
		return
	}

	// Сохраняем выбор и переходим к опыту
	userState := &state.UserState{
		Step:       state.StepAwaitingExperience,
		Tech:       tech,
		Experience: "",
	}
	b.state.Set(chatID, userState)

	// Кнопки опыта
	expOptions := []string{"0", "1", "2", "3", "4", "5", "6", "более 6"}
	var expButtons []tgbotapi.KeyboardButton
	for _, exp := range expOptions {
		expButtons = append(expButtons, tgbotapi.NewKeyboardButton(exp))
	}

	// Разбиваем на две строки по 4 кнопки
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(expButtons[0], expButtons[1], expButtons[2], expButtons[3]),
		tgbotapi.NewKeyboardButtonRow(expButtons[4], expButtons[5], expButtons[6], expButtons[7]),
	)
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, "2. Сколько лет у вас опыта?")
	msg.ReplyMarkup = keyboard
	b.client.Send(msg)
}

func (b *telegramBot) handleExperienceSelection(chatID int64, exp string) {
	validExp := map[string]bool{
		"0": true, "1": true, "2": true, "3": true,
		"4": true, "5": true, "6": true, "более 6": true,
	}

	if !validExp[exp] {
		// Повторяем вопрос
		expButtons := make([]tgbotapi.KeyboardButton, 0)
		for _, e := range []string{"0", "1", "2", "3", "4", "5", "6", "более 6"} {
			expButtons = append(expButtons, tgbotapi.NewKeyboardButton(e))
		}
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(expButtons[:4]...),
			tgbotapi.NewKeyboardButtonRow(expButtons[4:]...),
		)
		keyboard.OneTimeKeyboard = true
		keyboard.ResizeKeyboard = true

		msg := tgbotapi.NewMessage(chatID, "Пожалуйста, выберите опыт из списка")
		msg.ReplyMarkup = keyboard
		b.client.Send(msg)
		return
	}

	// Получаем полное состояние
	userState := b.state.Get(chatID)
	userState.Experience = exp

	// Формируем ответ
	response := "Вас интересует зарплата для IT для работы с " + userState.Tech + " и опыт работы " + exp
	if exp == "0" {
		response += " лет"
	} else if exp == "1" {
		response += " год"
	} else if exp == "более 6" {
		response += " лет"
	} else {
		// 2–6 → "года" или "лет"?
		// Для простоты — "лет"
		response += " лет"
	}

	// Отправляем итог
	msg := tgbotapi.NewMessage(chatID, response)
	b.client.Send(msg)

	// Сбрасываем состояние (диалог завершён)
	b.state.Clear(chatID)
}
