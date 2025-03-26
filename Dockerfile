FROM golang:1.24

WORKDIR /api

# Copy go.mod first for caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/api .

CMD ["api"]