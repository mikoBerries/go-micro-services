FROM golang:1.20.4-alpine3.18 AS Builder

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

RUN chmod +x authApp


FROM alpine:3.18

RUN mkdir /app

WORKDIR /app

# COPY authApp /app
COPY --from=Builder /app/authApp /app

# RUN chmod +x authApp

EXPOSE 80 

CMD [ "/app/authApp" ]