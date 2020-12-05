FROM golang:1.15-buster
WORKDIR /app
COPY . .

RUN go build -o /app/jinya-releases

EXPOSE 8090

CMD ["/app/jinya-releases"]