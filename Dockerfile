FROM golang:1.20.5-alphine3.18

WORKDIR /manager_api

COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./src ./src

CMD ["go", "run", "./main.go"]