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
# COPY .env .

# RUN tree -a

# RUN go build -c ./cmd/api/ -v -o ./main
# # RUN tree --gitignore
#
# RUN swag init -d cmd/api/,api/resource/system/,api/resource/books/
#
# CMD ["cmd/api/main"]

# EXPOSE 8009

RUN git config --global --add safe.directory /app

CMD swag init -d cmd/api/,api/resource/system/,api/resource/books/ && CompileDaemon --exclude-dir="docs" --build="./build.sh" --command="./main" --color

# ENTRYPOINT CompileDaemon --build="go build -a -installsuffix cgo -o main ." --command=./main
