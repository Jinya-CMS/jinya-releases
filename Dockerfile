FROM library/golang:1.22-alpine AS build
WORKDIR /app
COPY . .

RUN go build -o /jinya-releases

FROM library/alpine:latest

WORKDIR /app

COPY --from=build /jinya-releases /app/jinya-releases

EXPOSE 8090

CMD ["/app/jinya-releases"]
