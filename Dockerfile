FROM golang:1.22.5

RUN apt-get install curl && \
    apt-get install git

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
    | tar -xz -C /tmp/ && \
    mv /tmp/migrate /usr/local/bin/ && \
    chmod +x /usr/local/bin/migrate

WORKDIR /greetbot

RUN git clone https://github.com/m1hunter/greetings_bot.git .

RUN go mod tidy

RUN go build -o greetingsbotik ./cmd/main.go

RUN echo "TELEGRAM_BOT_TOKEN=7090855604:AAE9Ui7J34yFjnQ7Be8EyRtomfkxTVFm7CU\nDB_HOST=localhost\nDB_PORT=5891\nDB_USER=postgres\nDB_PASSWORD=postgres\nDB_NAME=postgres" > .env

CMD ["./greetingsbotik"]