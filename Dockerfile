FROM golang:latest

RUN apt update && apt install -y git && apt install -y tree

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go get github.com/githubnemo/CompileDaemon

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/githubnemo/CompileDaemon

ENV PATH="${PATH}:$HOME/go/bin"
ENV PATH="${PATH}:/usr/local/go/bin"

COPY . .

RUN git config --global --add safe.directory /app

CMD swag init -d cmd/api/,api/resource/system/,api/resource/books/ && CompileDaemon --exclude-dir="docs" --build="./build.sh" --command="./main" --color
