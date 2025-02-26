FROM registry.fluidcloud.bskyb.com/digital-trading/golang-arm:1.24

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o clipboard-app

EXPOSE 8088

CMD ["./clipboard-app"]