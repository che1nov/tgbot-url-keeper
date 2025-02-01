# Используем официальный образ Go
FROM golang:1.21-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN go build -o bot .

# Запускаем приложение
CMD ["./bot"]