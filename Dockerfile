FROM harbor.ulbricht.casa/proxy/library/golang:1.21-alpine AS build
WORKDIR /app
COPY . .

RUN go build -o /jinya-releases

FROM harbor.ulbricht.casa/proxy/library/alpine:latest

WORKDIR /app

COPY --from=build /jinya-releases /app/jinya-releases
COPY --from=build /app/static /app/static
COPY --from=build /app/templates /app/templates

EXPOSE 8090

CMD ["/app/jinya-releases"]
