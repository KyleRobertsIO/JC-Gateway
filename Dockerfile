FROM golang:1.20-bullseye

WORKDIR /manager_api

COPY ./entrypoint.sh ./entrypoint.sh
COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./src ./src

RUN chmod +x ./entrypoint.sh
RUN go mod download

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
