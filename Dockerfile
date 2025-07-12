FROM docker.io/library/alpine:latest

COPY jinya-releases /app/jinya-releases

CMD ["/app/jinya-releases"]