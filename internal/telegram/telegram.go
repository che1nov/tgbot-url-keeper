package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"net/url"
	"tgbot-url-keeper/internal/repository/storage"
)

const (
	StateDefault = iota
	StateAwaitingLinkID
)

var userStates = make(map[int64]int)

func SetupBot(bot *telebot.Bot) {

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
			`👋 Здравствуйте, уважаемый!
Меня зовут Василий, и я — скромный коллекционер ссылок! 🏚️📚

Я вовсе не бездомный в привычном понимании, просто веду минималистичный образ жизни, собирая исключительно ценные и интересные интернет-ресурсы. 💼✨

📎 Если у вас есть ссылка, которую вы бы хотели сохранить, пожалуйста, отправьте её мне — и я с удовольствием положу её в свою коллекцию.
📂 Желаете ознакомиться с ранее сохранёнными материалами? Конечно, я их аккуратно разложил по полочкам!
🗑️ Хотите что-то удалить? Разумеется, с большим уважением уберу ненужную ссылку.

Благодарю за доверие! Пусть ваш день будет полон хороших открытий и полезных ссылок.`,
			menu,
		)
	})

	bot.Handle(&btnStart, func(c telebot.Context) error {
		return c.Send(
			`👋 Здравствуйте, уважаемый!
Меня зовут Василий, и я — скромный коллекционер ссылок! 🏚️📚

Я вовсе не бездомный в привычном понимании, просто веду минималистичный образ жизни, собирая исключительно ценные и интересные интернет-ресурсы. 💼✨

📎 Если у вас есть ссылка, которую вы бы хотели сохранить, пожалуйста, отправьте её мне — и я с удовольствием положу её в свою коллекцию.
📂 Желаете ознакомиться с ранее сохранёнными материалами? Конечно, я их аккуратно разложил по полочкам!
🗑️ Хотите что-то удалить? Разумеется, с большим уважением уберу ненужную ссылку.

Благодарю за доверие! Пусть ваш день будет полон хороших открытий и полезных ссылок.`,
			menu,
		)
	})

	bot.Handle(&btnSave, func(c telebot.Context) error {
		return c.Send("📥 Отправьте мне ссылку, и я аккуратно добавлю её в вашу коллекцию!")
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
				return c.Send("❌ Такая ссылка уже есть.")
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
		return c.Send(`📜 Дорогие мои, как пользоваться моими скромными услугами?
💾 Есть ценная ссылочка? Смело отправляйте мне! Я сохраню её бережно, словно последний сухарик в кармане.
📂 Хотите посмотреть, что уже накопилось? Нажмите «Мои ссылки», и я, как настоящий архивариус, достану всё из своего надёжного мешка.
🗑️ Надо избавиться от ненужной ссылки? Жмите «Удалить ссылку», и я, с уважением к вашим решениям, отправлю её на покой.
❓ Вопросы, сомнения, экзистенциальный кризис? Нажмите «Помощь», и я помогу, чем смогу.`)
	})

}

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
