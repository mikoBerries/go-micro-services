FROM golang:1.20.4-alpine3.18 AS Builder

LABEL stage=Builder

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api

RUN chmod +x mailerApp


FROM alpine:3.18

RUN mkdir /app

# COPY brokerApp /app
COPY --from=Builder /app/mailerApp /app

COPY  templates /app/templates

WORKDIR /app

CMD [ "/app/mailerApp" ]