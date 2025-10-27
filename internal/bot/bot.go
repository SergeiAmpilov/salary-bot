// internal/bot/bot.go
package bot

import (
	"fmt"
	"log"
	"salary-bot/internal/bot/state"
	"salary-bot/internal/salary/model"
	"salary-bot/internal/salary/service"
	"strconv"
	"time"

	userrepo "salary-bot/internal/user/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramBot struct {
	client         *tgbotapi.BotAPI
	state          *state.Manager
	salaryService  service.Service
	userRepository userrepo.Repository
}

func NewBot(
	token string,
	salarySvc service.Service,
	userRepo userrepo.Repository,
) (Bot, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	client.Debug = false
	log.Printf("Авторизован бот: @%s", client.Self.UserName)

	return &telegramBot{
		client:         client,
		state:          state.NewManager(),
		salaryService:  salarySvc,
		userRepository: userRepo,
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

		user := update.Message.From

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		// Логируем
		log.Printf("[ChatID=%d] %s", chatID, text)

		// Обновляем/создаём пользователя
		_ = b.userRepository.Upsert(user.ID, user.UserName, user.FirstName)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start", "restart":
				b.handleStart(chatID)
			default:
				msg := tgbotapi.NewMessage(chatID, "Неизвестная команда. Напишите /start.")
				b.client.Send(msg)
			}
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
			b.handleExperienceSelection(chatID, text, user.ID)
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

func (b *telegramBot) handleExperienceSelection(chatID int64, expInput string, userID int64) {
	validExp := map[string]bool{
		"0": true, "1": true, "2": true, "3": true,
		"4": true, "5": true, "6": true, "более 6": true,
	}

	if !validExp[expInput] {
		// Повтор вопроса (как раньше)
		expOptions := []string{"0", "1", "2", "3", "4", "5", "6", "более 6"}
		var expButtons []tgbotapi.KeyboardButton
		for _, e := range expOptions {
			expButtons = append(expButtons, tgbotapi.NewKeyboardButton(e))
		}
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(expButtons[0], expButtons[1], expButtons[2], expButtons[3]),
			tgbotapi.NewKeyboardButtonRow(expButtons[4], expButtons[5], expButtons[6], expButtons[7]),
		)
		keyboard.OneTimeKeyboard = true
		keyboard.ResizeKeyboard = true

		msg := tgbotapi.NewMessage(chatID, "Пожалуйста, выберите опыт из списка")
		msg.ReplyMarkup = keyboard
		b.client.Send(msg)
		return
	}

	// Получаем состояние
	userState := b.state.Get(chatID)
	userState.Experience = expInput

	// Преобразуем опыт в число
	var expYears int
	if expInput == "более 6" {
		expYears = 7
	} else {
		expYears, _ = strconv.Atoi(expInput)
	}

	// Формируем фильтр
	filter := &model.FilterDTO{
		Tech: &userState.Tech,
		Type: strPtr("remote"),
	}

	// Опыт: ищем записи, где диапазон покрывает пользователя
	filter.ExperienceMin = &expYears // записи с experience_max >= expYears
	filter.ExperienceMax = &expYears // записи с experience_min <= expYears

	// Дата: последние 30 дней
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02 15:04:05")
	filter.CreatedAtFrom = &thirtyDaysAgo

	// Вызываем сервис напрямую
	salaries, err := b.salaryService.Filter(filter)
	if err != nil {
		log.Printf("Ошибка фильтрации: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные. Попробуйте позже.")
		b.client.Send(msg)
		b.state.Clear(chatID)
		return
	}

	// Расчёт среднего
	avgMin, avgMax := b.calculateAverage(salaries)

	// После расчёта и перед отправкой ответа:
	_ = b.userRepository.IncrementCalculation(userID)

	// Формируем ответ
	var response string
	if avgMin == 0 && avgMax == 0 {
		response = fmt.Sprintf("К сожалению, нет данных о зарплатах для %s с опытом %s год(а/лет) на удалёнке за последние 30 дней.", userState.Tech, expInput)
	} else {
		minStr := "не указано"
		maxStr := "не указано"
		if avgMin > 0 {
			minStr = formatSalary(avgMin)
		}
		if avgMax > 0 {
			maxStr = formatSalary(avgMax)
		}
		response = fmt.Sprintf(
			"Для программиста %s с опытом работы %s год(а/лет) средняя зарплата на удалёнке находится в диапазоне от %s до %s.",
			userState.Tech,
			expInput,
			minStr,
			maxStr,
		)
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/start"),
		),
	)
	keyboard.OneTimeKeyboard = false // оставляем клавиатуру видимой
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, response)
	msg.ReplyMarkup = keyboard
	b.client.Send(msg)
	b.state.Clear(chatID)
}

func strPtr(s string) *string {
	return &s
}

func (b *telegramBot) calculateAverage(salaries []*model.Salary) (avgMin, avgMax int) {
	if len(salaries) == 0 {
		return 0, 0
	}

	var sumMin, sumMax, countMin, countMax int64
	for _, s := range salaries {
		if s.SalaryMin > 0 {
			sumMin += int64(s.SalaryMin)
			countMin++
		}
		if s.SalaryMax > 0 {
			sumMax += int64(s.SalaryMax)
			countMax++
		}
	}

	if countMin > 0 {
		avgMin = int(sumMin / countMin)
	}
	if countMax > 0 {
		avgMax = int(sumMax / countMax)
	}

	return avgMin, avgMax
}

func formatSalary(amount int) string {
	s := fmt.Sprintf("%d", amount)
	if len(s) <= 3 {
		return s
	}
	var result string
	for i := len(s) - 1; i >= 0; i-- {
		if (len(s)-1-i)%3 == 0 && i != len(s)-1 {
			result = " " + result
		}
		result = string(s[i]) + result
	}
	return result + " ₽"
}
