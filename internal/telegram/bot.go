package telegram

import "gopkg.in/telebot.v3"

func SetupBot(bot *telebot.Bot) {
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnStart := menu.Text("🚀 Старт")
	btnSave := menu.Text("💾 Сохранить ссылку")
	btnLinks := menu.Text("📂 Мои ссылки")
	btnDelete := menu.Text("🗑️ Удалить ссылку")

	menu.Reply(
		menu.Row(btnStart),
		menu.Row(btnSave),
		menu.Row(btnLinks),
		menu.Row(btnDelete),
	)

	// Подключаем обработчики
	bot.Handle("/start", func(c telebot.Context) error { return handleStart(c, menu) })
	bot.Handle(&btnStart, func(c telebot.Context) error { return handleStart(c, menu) })
	bot.Handle(&btnSave, handleSaveLink)
	bot.Handle(&btnLinks, handleGetLinks)
	bot.Handle(&btnDelete, handleDeleteLink)
	bot.Handle(telebot.OnText, handleText)
}
