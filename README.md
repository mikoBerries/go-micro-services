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
```console
// Login to dockerhub
$ docker login -u your_user_name - The -u option allows us to pass our user name.
// Build images with to push to docker hub
$ docker build -f some-service.dockerfile -t some_service_name:tag .
// Push images to loged user in dockerhub
$ docker push some_service_name:tag
```
2. After that make swarm.yml just like docker-compose.yml but images are refering to dockerhub instead of local images
3. Then Do in console to initial swarm will create 1 NODE as Manager
```console
$ docker swarm init
```
4. We can add more Worker / Manager to node (follow this command instruction to get token)
````console
$ docker swarm join-token worker
$ docker swarm join-token manager
````
5. To deploy Docker Swarm we must use Docker-Stack in same level of swarm.yml file (Docker tollBox are not supported to making swarm)
````console
// swarm_name as name of node/swarm we created
$ docker stack deploy -c swarm.yml swarm_name
````
6. To check docker service that running
````console
$ docker service ls
````
7. We can scaling service by using (Images mode must Be Replicated instead global)
````console
// Scaling up some_service_name to 3 
$ docker service scale some_service_name = 3
// Scaling Down some_service_name to 3 
$ docker service scale some_service_name = 2
````
8. Updating 1 some-service from new images in dockerhub that created at project code
````console
// Build images with to push to docker hub
$ docker build -f some-service.dockerfile -t some_service_name:tag .
// Push images to loged user in dockerhub
$ docker push some_service_name:tag
// updating service with new images (some_service_name) from dockerhub
// We can update it to new tags or rollback using tag images we need
$ docker service update --image some_service_name:tag swarm_name_some_service_name
````
9. Stopping swarm Service  & removing swarm
````console
// Stopping swarm service 
$ docker service scale swarm_name=0
// Removing swarm
$ docker stack rm swram_name
// using --force is we leave Manager swarm not Worker
$ docker swarm leave --force
````

## Caddy server
-------
1. Caddy 2 is a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go-lang (https://hub.docker.com/_/caddy)
2. oficial website (https://caddyserver.com/business)
3. how to modify host file in windows (https://www.thewindowsclub.com/hosts-file-in-windows)
4. Cloud service Cost efective for experimental:
    - https://www.linode.com/
    - https://www.digitalocean.com/
    - https://www.vultr.com/
5. connecting to cloud server we make and setting some for usage
    - Adding user and giving sudo prefilege beside root user:
    ```console
    // adding user in ubuntu
    $ add user some-user-name
    // giving user previlege
    $ usermod -aG sudo some-user-name
    ```
    - Setting firewall in ubuntu: 
    ```console
    $ ufw allow ssh
    $ ufw allow http
    $ ufw allow https
    // port for docker
    $ ufw allow 2377/tcp
    $ ufw allow 7946/tcp
    $ ufw allow 7946/udp
    $ ufw allow 4789/udp
    // mailhog port
    $ ufw allow 8025/udp
    $ ufw enable
    $ ufw status
    ```
6. After all set up ubuntu are ready to build docker machine (https://docs.docker.com/engine/install/ubuntu/)
    - Follow instruction in secdtion Install using the apt repository (there are few method to instal docker)
    - After finish we can check our docker are instaled
    ```console
    $ where docker
    ```
7. Change ubuntu host name to appropied host name
```console
$ sudo hostnamectl set-hostname node-1
```
8. update host and write our node ip in the host-list
```console
$ sudo vi /etc/hosts
```
9. Setting DNS using godaddy
10. Setting up docker swarm inside ubuntu server with few node server
```console
// this will make a Swarm manager in node_ip server
$ sudo docker swarm init --advertise-addr node_ip

// change to other node of server and execute docker swarm join --token to make wokrker listed to node Manager server
```
11. udpating caddy file and caddy docker file so caddy file with fetch porxy named our domain
12. dont forget to create folder inside the ubuntu server 
13. starting docker swarm like in docker swarm section no 5 ~
14. adding user to docker gruop
```console
$ sudo usermod -aG docker some-user-name
```
15. Checking all container in node using 
```console
$ docker node ps
```
16. Docker volume need placement to refering where is volume located?
```console
    ## this is example when volume path exist in node-1
    placement:
        constraints:
            - node.hostname == node-1
```