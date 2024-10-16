FROM golang:1.22.5

RUN apt-get update && apt-get install -y curl git

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
    | tar -xz -C /tmp/ && \
    mv /tmp/migrate /usr/local/bin/ && \
    chmod +x /usr/local/bin/migrate

WORKDIR /greetbot

# Исправленный путь для git clone
RUN git clone https://github.com/m1hunter/greetings_bot.git . && \
    go mod tidy && \
    go build -o greetingsbotik ./cmd/main.go


CMD ["./greetingsbotik"]