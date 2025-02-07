# Используем базовый образ с поддержкой cgo
FROM golang:1.21-alpine

# Устанавливаем зависимости для cgo
RUN apk add --no-cache gcc musl-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем приложение с включенным cgo
RUN CGO_ENABLED=1 go build -o bot ./cmd/main.go

# Запускаем приложение
CMD ["./bot"]