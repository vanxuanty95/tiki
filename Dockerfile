FROM golang:1.16-alpine

WORKDIR /go/src/tiki/cmd/api
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api
RUN ["chmod", "+x", "main"]
EXPOSE 8080

CMD ["./main"]