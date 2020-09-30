FROM golang:1.15.2-buster AS builder

WORKDIR /src
COPY . .
RUN go build cmd/souq/souq.go

FROM debian:buster-slim

ENV APPS_DIR=/apps
COPY --from=builder /src/souq /bin/souq
COPY --from=builder /src/apps /apps
RUN chmod +x /bin/souq

ENTRYPOINT ["souq"]