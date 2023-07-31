FROM quay.imanuel.dev/dockerhub/library---golang:1.20-alpine
WORKDIR /app
COPY . .

RUN go build -o /jinya-releases

FROM quay.imanuel.dev/dockerhub/library---alpine:latest

COPY --from=build /jinya-releases /jinya-releases

EXPOSE 8090

CMD ["/jinya-releases"]