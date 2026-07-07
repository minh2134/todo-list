FROM golang:1.26-nanoserver AS build

WORKDIR /app

COPY go.mod go.sum main.go internal ./
RUN go mod tidy

RUN CGO_ENABLED=1 GOOS=linux go build -o todo-list .


FROM alpine:3.23.5 AS prod


WORKDIR /opt/todo-list

COPY --from-builder /app/todo-list .
COPY static/index.css static/index.html static/

EXPOSE 3000

CMD ["./server"]
