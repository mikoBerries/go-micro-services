FROM golang:1.20.4-alpine3.18 AS Builder

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x brokerApp


FROM alpine:3.18

RUN mkdir /app

WORKDIR /app

# COPY brokerApp /app
COPY --from=Builder /app/brokerApp /app

# RUN chmod +x brokerApp

EXPOSE 80 

CMD [ "/app/brokerApp" ]