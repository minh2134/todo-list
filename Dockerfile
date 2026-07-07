FROM golang:1.26-trixie AS build

WORKDIR /app

COPY go.mod go.sum main.go  ./
COPY internal/ ./internal/
RUN go mod tidy

RUN CGO_ENABLED=1 GOOS=linux go build -o todo-list .


FROM debian:trixie-backports AS prod

WORKDIR /opt/todo-list

COPY --from=build /app/todo-list .
COPY static/index.css static/index.html static/
COPY templ/ templ/

EXPOSE 3000

CMD ["./todo-list"]
