FROM golang:1.20.4-alpine3.18 AS Builder

LABEL stage=Builder

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 go build -o loggerServiceApp ./cmd/api

RUN chmod +x loggerServiceApp


FROM alpine:3.18

RUN mkdir /app

WORKDIR /app

COPY --from=Builder /app/loggerServiceApp /app

CMD [ "/app/loggerServiceApp" ]