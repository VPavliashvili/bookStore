FROM golang:latest

RUN apt update && apt install -y git && apt install -y tree

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="${PATH}:$HOME/go/bin"
ENV PATH="${PATH}:/usr/local/go/bin"

COPY . .
# COPY .env .

# RUN tree -a

RUN go build -C ./cmd/api/ -v -o ./main
# RUN tree --gitignore

CMD ["cmd/api/main"]

# RUN swag init -d cmd/api/,api/resource/system/,api/resource/books/

# EXPOSE 8009
