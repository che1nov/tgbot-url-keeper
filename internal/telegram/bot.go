package telegram

import "gopkg.in/telebot.v3"

func SetupBot(bot *telebot.Bot) {
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnStart := menu.Text("🚀 Поехали")
	btnSave := menu.Text("💾 Сохранить ссылку")
	btnLinks := menu.Text("📂 Мои ссылки")
	btnDelete := menu.Text("🗑️ Удалить ссылку")
	btnHelp := menu.Text("❓ Хэлп")

	menu.Reply(
		menu.Row(btnStart),
		menu.Row(btnSave),
		menu.Row(btnLinks),
		menu.Row(btnDelete),
		menu.Row(btnHelp),
	)

	// Передаем `menu` в обработчик /start
	bot.Handle("/start", func(c telebot.Context) error { return handleStart(c, menu) })
	bot.Handle(&btnStart, func(c telebot.Context) error { return handleStart(c, menu) })
	bot.Handle(&btnSave, handleSaveLink)
	bot.Handle(&btnLinks, handleGetLinks)
	bot.Handle(&btnDelete, handleDeleteLink)
	bot.Handle(&btnHelp, handleHelp)
	bot.Handle(telebot.OnText, handleText)
}
