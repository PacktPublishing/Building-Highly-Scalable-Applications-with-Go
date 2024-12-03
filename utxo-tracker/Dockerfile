FROM golang:1.22-alpine AS builder

WORKDIR /parent/utxo-tracker

RUN apk add --no-cache git tzdata build-base

RUN go install github.com/magefile/mage@latest

COPY ./utxo-tracker/go.mod ./utxo-tracker/go.sum ./
RUN go mod download

COPY .git ../.git
COPY ./utxo-tracker/. .
RUN mage build

FROM gcr.io/distroless/static:nonroot
COPY --chown=nonroot:nonroot --from=builder /parent/utxo-tracker/account-service /home/nonroot/

USER nonroot:nonroot
EXPOSE 8080
EXPOSE 8081

CMD ["/home/nonroot/account-service"]
