# Introduction
--------------
For a long time, web applications were usually a single application that handled everything—in other words, a monolithic application. This monolith handled user authentication, logging, sending email, and everything else. While this is still a popular (and useful) approach, today, many larger scale applications tend to break things up into microservices. Today, most large organizations are focused on building web applications using this approach, and with good reason.

Microservices, also known as the microservice architecture, are an architectural style which structures an application as a loosely coupled collection of smaller applications. The microservice architecture allows for the rapid and reliable delivery of large, complex applications. Some of the most common features for a microservice are:

it is maintainable and testable;

it is loosely coupled with other parts of the application;

it  can deployed by itself;

it is organized around business capabilities;

it is often owned by a small team.

In this course, we'll develop a number of small, self-contained, loosely coupled microservices that will will communicate with one another and a simple front-end application with a REST API, with RPC, over gRPC, and by sending and consuming messages using AMQP, the Advanced Message Queuing Protocol. The microservices we build will include the following functionality:

A Front End service, that just displays web pages;

An Authentication service, with a Postgres database;

A Logging service, with a MongoDB database;

A Listener service, which receives messages from RabbitMQ and acts upon them;

A Broker service, which is an optional single point of entry into the microservice cluster;

A Mail service, which takes a JSON payload, converts into a formatted email, and send it out.

All of these services will be written in Go, commonly referred to as Golang, a language which is particularly well suited to building distributed web applications.

We'll also learn how to deploy our distributed application to a Docker Swarm and Kubernetes, and how to scale up and down, as necessary, and to update individual microservices with little or no downtime.


* https://www.cyberciti.biz/faq/alpine-linux-install-bash-using-apk-command/

## Docker Swarm
----------------
1. To use Docker swarm we must pushing our images to dockerhub (like github but for docker images)
2. After that make swarm.yml just like docker-compose.yml but images are refering to dockerhub instead of local images
3. Then Do in console to initial swarm will create 1 NODE as Manager
````console
$ docker swarm init
````
4. We can add more Worker / Manager to node (follow this command instruction to get token)
````console
$ docker swarm join-token worker
$ docker swarm join-token manager
````
5. To deploy Docker Swarm we must use Docker-Stack in same level of swarm.yml file (Docker tollBox are not supported to making swarm)
````console
$ docker stack deploy -c swarm.yml someSwarmName
````
6. To check docker service that running
````console
$ docker service ls
````