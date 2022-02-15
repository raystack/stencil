FROM alpine:3.13

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY stencil /usr/bin/stencil

EXPOSE 8080
ENTRYPOINT ["stencil"]
