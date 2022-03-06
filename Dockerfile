FROM golang:1.17.8-alpine

RUN apk update && apk add git

WORKDIR /app

COPY . .

RUN go mod tidy -v

CMD ["go", "run" ,"cmd/api/main.go"]
