FROM golang:1.21
WORKDIR /var/server/trade-bot

COPY go.mod go.sum ./
COPY .deploy/local/.env .
RUN go mod download && go mod verify

COPY . .

CMD ["go","run","main.go"]
