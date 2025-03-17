# Telegram URL Keeper Bot

**Telegram URL Keeper Bot** is a simple Telegram bot written in Go using the [Telebot v3](https://github.com/tucnak/telebot) library. The bot allows users to save, view, and delete URLs. It implements the following features:

- **Saving URLs:**  
  Users can send a URL in the chat, and the bot will check its validity. If the URL is valid, it gets saved.

- **Viewing saved URLs:**  
  By pressing the "📂 My Links" button, the bot sends a list of previously saved URLs.

- **Deleting a URL:**  
  By pressing the "🗑️ Delete Link" button, the bot asks the user for the ID of the URL to delete and removes the selected URL.

- **Help:**  
  The "❓ Help" button displays instructions on how to use the bot.

- **Start:**  
  The `/start` command or pressing the "🚀 Start" button displays a welcome message with a description of the bot's functionality.

## Project Structure

The project is organized according to the principles of modularity. The main components are:

- **package telegram**  
  Contains the setup and command handling of the bot, defines user states for various operations (e.g., deleting a URL).

- **package storage (internal/repository/storage)**  
  Implements operations for saving, retrieving, and deleting URLs in the storage (e.g., in a database or in-memory).

- **Used libraries:**
  - [Telebot v3](https://github.com/tucnak/telebot) for interacting with the Telegram API.
  - The standard `net/url` package for checking URL validity.
  - The `log` package for logging operations.