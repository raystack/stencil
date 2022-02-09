FROM alpine:3.13

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY stencil .

EXPOSE 8080
ENTRYPOINT ["./stencil"]
