FROM golang:1.27.3-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o checklinks

ENTRYPOINT ["/app/checklinks"]
