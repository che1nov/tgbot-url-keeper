# Используем официальный образ Go
FROM golang:1.21-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Скачиваем зависимости
RUN go mod download

# Установите CGO_ENABLED=1
ENV CGO_ENABLED=1

# Собираем приложение
RUN go build -o bot .

# Запускаем приложение
CMD ["./bot"]