FROM golang:1.16-alpine3.13 AS builder
WORKDIR /go/src/github.com/odpf/stencil
COPY . .
RUN apk add make bash
RUN make dist

FROM alpine:3.13
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
EXPOSE 8080
COPY ./server/migrations /root/migrations
COPY --from=builder /go/src/github.com/odpf/stencil/dist/linux-amd64/stencil .
ENTRYPOINT ["./stencil"]
