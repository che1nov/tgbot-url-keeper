package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3" // Импорт драйвера SQLite
)

var db *sql.DB

// Init инициализирует базу данных
func Init() error {
	var err error
	db, err = sql.Open("sqlite3", "./links.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Создание таблицы, если она не существует
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			url TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

// SaveLink сохраняет ссылку в базе данных
func SaveLink(userID int64, url string) error {
	_, err := db.Exec("INSERT INTO links (user_id, url) VALUES (?, ?)", userID, url)
	return err
}

// GetLinks возвращает все ссылки пользователя
func GetLinks(userID int64) ([]Link, error) {
	rows, err := db.Query("SELECT id, url FROM links WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		if err := rows.Scan(&link.ID, &link.URL); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func DeleteLink(userID int64, linkID string) error {
	id, err := strconv.Atoi(linkID)
	if err != nil {
		log.Printf("[ERROR] Некорректный ID ссылки: %s (userID=%d)", linkID, userID)
		return fmt.Errorf("неверный ID ссылки: %v", err)
	}

	log.Printf("[DEBUG] Запрос DELETE: id=%d, userID=%d", id, userID)
	result, err := db.Exec("DELETE FROM links WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		log.Printf("[ERROR] Ошибка SQL-запроса: %v", err)
		return fmt.Errorf("ошибка при удалении: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("[DEBUG] Удалено строк: %d", rowsAffected)

	if rowsAffected == 0 {
		log.Printf("[WARN] Ссылка не найдена: id=%d, userID=%d", id, userID)
		return fmt.Errorf("ссылка с ID %d не найдена", id)
	}

	return nil
}

// Link представляет структуру ссылки
type Link struct {
	ID  int
	URL string
}
