FROM quay.imanuel.dev/dockerhub/library---golang:1.20-alpine AS build
WORKDIR /app
COPY . .

RUN go build -o /jinya-releases

FROM quay.imanuel.dev/dockerhub/library---alpine:latest

COPY --from=build /jinya-releases /app/jinya-releases
COPY --from=build /app/static /app/static
COPY --from=build /app/templates /app/templates

EXPOSE 8090

CMD ["/jinya-releases"]