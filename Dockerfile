FROM golang:1-alpine AS builder

RUN apk add --update --no-cache ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./website /website

ENV GITHUB_TOKEN=""

EXPOSE 8080

CMD ["/website"]
