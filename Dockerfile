FROM golang:1.20-bullseye

WORKDIR /manager_api

ENV LOGGER_LOG_LEVEL ''
ENV LOGGER_FILE_PATH ''
ENV GIN_PORT ''
ENV GIN_MODE ''
ENV AZURE_AUTH_TYPE ''
ENV AZURE_AUTH_CLIENT_ID ''
ENV AZURE_AUTH_CLIENT_SECRET ''
ENV AZURE_AUTH_TENANT_ID ''

COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./src ./src

RUN go mod download

EXPOSE 8080

CMD ["go", "run", "./main.go"]