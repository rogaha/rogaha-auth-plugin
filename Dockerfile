FROM alpine:3.4

RUN mkdir -p /run/docker/plugins

COPY rogaha-auth-plugin rogaha-auth-plugin

CMD ["rogaha-auth-plugin"]
