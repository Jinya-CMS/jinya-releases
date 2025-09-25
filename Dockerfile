ARG  CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX

FROM $CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX/library/alpine:latest

COPY jinya-releases /app/jinya-releases

CMD ["/app/jinya-releases"]