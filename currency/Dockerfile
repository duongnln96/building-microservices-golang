FROM golang:1.16-alpine

WORKDIR /app

COPY . ./

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

RUN go build .

EXPOSE 9092

CMD ["go", "run", "." ]
