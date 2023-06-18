FROM golang:1.20-alpine AS Builder

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x brokerApp


FROM alpine:latest

RUN mkdir /app

COPY --from=Builder /app/brokerApp /app

EXPOSE 80 

CMD [ "/app/brokerApp" ]