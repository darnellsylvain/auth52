FROM golang:1.24

WORKDIR /api

RUN go install github.com/air-verse/air@latest

# Copy go.mod first for caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/api .

CMD ["air", "-c", ".air.toml"]
