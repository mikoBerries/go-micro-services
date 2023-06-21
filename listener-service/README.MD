# Listener-service
------------------

* listener service serve as listener to rabbitmq container
* Rabbitmq images (https://hub.docker.com/_/rabbitmq)
* Rabbitmq Driver:
    - github.com/rabbitmq/amqp091-go
* RabbitMq Documentation - https://www.rabbitmq.com/tutorials/tutorial-one-go.html
* RabbitMq are publish/subscribe pattern.
* RabbitMQ is designed as a dedicated message broker, where Redis is a database that you can use as a message broker. (https://www.trustradius.com/compare-products/rabbitmq-vs-redis#trustradius-insights)
* For advance case we can use redis instead rabbitmq
* Message durability : We have learned how to make sure that even if the consumer dies, the task isn't lost. But our tasks will still be lost if RabbitMQ server stops.
* The core idea in the messaging model in RabbitMQ is that the producer never sends any messages directly to a queue. Actually, quite often the producer doesn't even know if a message will be delivered to any queue at all.
* There are a few exchange types available: direct, topic, headers and fanout.