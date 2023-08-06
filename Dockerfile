FROM golang:1.20-alpine

WORKDIR /go-kafka

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN build

CMD ["./bin/app"]