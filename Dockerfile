FROM golang:1.20-bullseye

WORKDIR /manager_api

COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./src ./src

RUN go mod download

EXPOSE 8080

CMD ["go", "run", "./main.go"]