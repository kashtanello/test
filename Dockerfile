FROM golang:1.23.2 AS builder

# Устанавливаем рабочий каталог внутри контейнера
WORKDIR /src/

# Копируем исходники в контейнер
COPY ./main.go /src/

# Инициализируем модуль Go
RUN go mod init mymodule

# Скачиваем все зависимости
RUN go mod download
RUN go get ./...

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

# Используем образ alpine:latest как базовый
FROM alpine:latest 

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарный файл из этапа builder
COPY --from=builder /src/app .
EXPOSE 8080
# Запускаем приложение 
ENTRYPOINT ["./app"]
