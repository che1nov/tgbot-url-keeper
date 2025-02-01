package main

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"net/url"
	"os"
	"tgbot-url-keeper/storage"
)

const (
	StateDefault = iota
	StateAwaitingLinkID
)

var userStates = make(map[int64]int) // userID -> state

func main() {
	if err := storage.Init(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Токен бота не указан")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
	})
	if err != nil {
		log.Fatal(err)
	}

	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnStart := menu.Text("🚀 Старт")
	btnSave := menu.Text("💾 Сохранить ссылку")
	btnLinks := menu.Text("📂 Мои ссылки")
	btnDelete := menu.Text("🗑️ Удалить ссылку")
	btnHelp := menu.Text("❓ Помощь")

	menu.Reply(
		menu.Row(btnStart),
		menu.Row(btnSave),
		menu.Row(btnLinks),
		menu.Row(btnDelete),
		menu.Row(btnHelp),
	)

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send(
			`Привет, друг! 👋 Меня зовут Вася, и я — бомж. 
Но не простой бомж, а *бомж-коллекционер*. 
Вместо бутылок я собираю ссылки! 🎯

Отправь мне ссылку, и я её сохраню в свою тележку. 
А если что-то непонятно — нажимай кнопки ниже! 😎`,
			menu,
		)
	})

	bot.Handle(&btnStart, func(c telebot.Context) error {
		return c.Send(
			`Привет, друг! 👋 Меня зовут Вася, и я — бомж. 
Но не простой бомж, а *бомж-коллекционер*. 
Вместо бутылок я собираю ссылки! 🎯

Отправь мне ссылку, и я её сохраню в свою тележку. 
А если что-то непонятно — нажимай кнопки ниже! 😎`,
			menu,
		)
	})

	bot.Handle(&btnSave, func(c telebot.Context) error {
		return c.Send("📥 Отправь мне ссылку, и я её сохраню в свою тележку!")
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID
		text := c.Text()

		if userStates[userID] == StateAwaitingLinkID {
			linkID := text

			log.Printf("Пользователь ввел linkID=%s для удаления", linkID)

			if err := storage.DeleteLink(userID, linkID); err != nil {
				return c.Send("❌ Ой, что-то пошло не так: " + err.Error())
			}

			userStates[userID] = StateDefault
			return c.Send("✅ Ссылка успешно удалена из моей тележки. Теперь её нет даже на свалке! 🗑️")
		}

		if isValidURL(text) {
			if err := storage.SaveLink(userID, text); err != nil {
				log.Printf("Ошибка при сохранении ссылки: %v", err)
				return c.Send("❌ Не удалось сохранить ссылку. Попробуй ещё раз позже.")
			}
			return c.Send("✅ Ссылка успешно сохранена в мою тележку: " + text)
		}

		return c.Send("🤔 Хм.. это не похоже на ссылку. Попробуй отправить что-то, типа хттп или схттп.")
	})

	bot.Handle(&btnLinks, func(c telebot.Context) error {
		userID := c.Sender().ID

		links, err := storage.GetLinks(userID)
		if err != nil {
			return c.Send("❌ Ой, что-то пошло не так: " + err.Error())
		}

		if len(links) == 0 {
			return c.Send("📭 Твоя тележка пуста. Может, добавим туда парочку ссылок? 😉")
		}

		var message string
		for _, link := range links {
			message += fmt.Sprintf("🔗 Ссылка #%d: %s\n", link.ID, link.URL)
		}

		return c.Send("📂 *Твои сохраненные ссылки:*\n" + message)
	})

	bot.Handle(&btnDelete, func(c telebot.Context) error {
		userID := c.Sender().ID
		userStates[userID] = StateAwaitingLinkID
		log.Printf("[DEBUG] Установлено состояние StateAwaitingLinkID для userID=%d", userID)
		return c.Send("🔢 Введи номер (ID) ссылки...")
	})

	bot.Handle(&btnHelp, func(c telebot.Context) error {
		return c.Send(`🤖 *Как использовать этого бомжа:*
- Просто отправь мне ссылку, и я её сохраню. Даже если это ссылка на рецепт борща. 🍲
- Чтобы посмотреть все ссылки, используй кнопку "📂 Мои ссылки".
- Чтобы удалить ссылку, используй кнопку "🗑️ Удалить ссылку".
- Если что-то непонятно, нажимай "❓ Помощь". 😊`)
	})

	bot.Start()
}

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
