FROM golang:1.22.5

RUN apt-get install curl && \
    apt-get install git

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
    | tar -xz -C /tmp/ && \
    mv /tmp/migrate /usr/local/bin/ && \
    chmod +x /usr/local/bin/migrate

WORKDIR /greetbot

RUN git clone https://github.com/m1hunter/greetings_bot.git . && \
    go mod tidy && \
    go build -o greetingsbotik ./cmd/main.go && \
    #echo "TELEGRAM_BOT_TOKEN=7090855604:AAE9Ui7J34yFjnQ7Be8EyRtomfkxTVFm7CU" > .env

CMD ["./greetingsbotik"]