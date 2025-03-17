package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"net/url"
	"strconv"
	"tgbot-url-keeper/internal/repository/storage"
)

const (
	StateDefault = iota
	StateAwaitingLinkID
)

var userStates = make(map[int64]int)

// Обработчик команды /start
func handleStart(c telebot.Context, menu *telebot.ReplyMarkup) error {
	welcomeText := `👋 Здравствуйте, уважаемый!
Меня зовут Василий, и я — скромный коллекционер ссылок! 🏚️📚

📎 Отправьте мне ссылку, и я сохраню её.
📂 Хотите посмотреть сохраненные ссылки? Нажмите «Мои ссылки».
🗑️ Нужно удалить ссылку? Жмите «Удалить ссылку».`
	return c.Send(welcomeText, menu)
}

// Обработчик сохранения ссылки
func handleSaveLink(c telebot.Context) error {
	return c.Send("📥 Отправьте мне ссылку, и я аккуратно добавлю её в вашу коллекцию!")
}

// Обработчик списка ссылок
func handleGetLinks(c telebot.Context) error {
	userID := c.Sender().ID
	links, err := storage.GetLinks(userID)
	if err != nil {
		log.Println("Ошибка получения ссылок:", err)
		return c.Send("❌ Ой, что-то пошло не так.")
	}

	if len(links) == 0 {
		return c.Send("📭 Ваша коллекция пока пуста.")
	}

	var message string
	for i, link := range links { // Используем i + 1 для правильной нумерации
		message += fmt.Sprintf("🔗 %d: %s\n", i+1, link.URL)
	}

	return c.Send("📂 *Ваши сохраненные ссылки:*\n" + message)
}

// Обработчик удаления ссылки
func handleDeleteLink(c telebot.Context) error {
	userID := c.Sender().ID
	userStates[userID] = StateAwaitingLinkID
	return c.Send("🔢 Введите номер (ID) ссылки для удаления...")
}

// Обработчик текста (сохранение или удаление)
func handleText(c telebot.Context) error {
	userID := c.Sender().ID
	text := c.Text()

	if userStates[userID] == StateAwaitingLinkID {
		linkID, err := strconv.Atoi(text)
		if err != nil {
			return c.Send("❌ Пожалуйста, введите корректный числовой ID.")
		}

		if err := storage.DeleteLink(userID, strconv.Itoa(linkID)); err != nil {
			log.Println("Ошибка при удалении ссылки:", err)
			return c.Send("❌ Ошибка удаления ссылки: " + err.Error())
		}

		userStates[userID] = StateDefault
		return c.Send("✅ Ссылка удалена!")
	}

	if isValidURL(text) {
		if err := storage.SaveLink(userID, text); err != nil {
			log.Println("Ошибка при сохранении ссылки:", err)
			return c.Send("❌ Такая ссылка уже есть.")
		}
		return c.Send("✅ Ссылка сохранена: " + text)
	}

	return c.Send("🤔 Это не похоже на ссылку. Попробуйте снова.")
}

func handleHelp(c telebot.Context) error {
	return c.Send(`📝 *Как пользоваться этим чудесным ботом?*

💾 *Добавить ссылку:* Просто отправь мне её, и я её сохраню навеки!
📂 *Просмотреть ссылки:* Жми «Открыть сундук» — и перед тобой раскроется древний архив.
🗑️ *Удалить ссылку:* Нажми «Избавиться от хлама», укажи ID — и я сотру её с лица земли.

💡 Совет: Используй меня как цифровую записную книжку, и твои ссылки никогда не затеряются!`)
}

// Функция проверки ссылки
func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
