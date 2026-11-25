FROM golang:1.20-alpine

WORKDIR /app

COPY . .

# DOWNLOAD AND INSTALL DEPENDENCIES:
RUN go get -d -v ./...

# BUILD THE APP
RUN go build -o api .

EXPOSE 8080

CMD ["./api"]