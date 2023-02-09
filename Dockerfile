FROM quay.imanuel.dev/dockerhub/library---golang:1.20-alpine
WORKDIR /app
COPY . .

RUN go build -o /app/jinya-releases

EXPOSE 8090

CMD ["/app/jinya-releases"]