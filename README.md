# commercial-studio-photos-generator

Simple yet Production Ready MLOps example of using any generating model

https://github.com/Dominux/commercial-studio-photos-generator/assets/55978340/31ca8f52-8830-4eb2-b079-d29cb77fa546

## Introduction

Nowadays, with generational neural networks increasingly becoming a significant part of our life, many devs got to intergrate them into production ready systems, no matter whether the system is a big entreprise or just a wrapper for such a model with some web, mobile or bot user interface. Programmatic madia creating, such as generating social media content like posts, videos, music and etc isn't already some new thing for now, but with neural networks becoming so popular and moreover affordable for ordinar people, makes such models so desirable to use.

You can get some model weights, it's architecture implementaion and just generate some content with it by hand - it's not that complex. But when it comes to integrating such a model into production ready system, even as simple as just that model wrapper, it's no longer that straighforward. You gotta think about load balancing, queued tasks, perfomance control, durability, rodustness and resource and health monitoring. CI/CD becomes a thing too since deploying process becomes complex and you don't wanna do that with your bare hands everytime.

This project uses [Stable Diffusion](https://huggingface.co/spaces/stabilityai/stable-diffusion) txt2img model via using [diffussers library](https://github.com/huggingface/diffusers) to demonstrate such an example of how to implement simple production ready model wrapper, but with the idea you can use any other generational model like:

- txt2img
- img2img
- txt2txt
- tts
- stt
- txt2video
- img2video

It also provides a simple demo web interface shown above to interact with, but you can use the arcitectual approache to implement any other way of interating via, maybe you would like to create a telegram bot or native mobile app, or use it as an external API via some [RESTful API](https://en.wikipedia.org/wiki/REST), [RPC](https://en.wikipedia.org/wiki/Remote_procedure_call), [gRPC](https://en.wikipedia.org/wiki/GRPC) or [message queues](https://en.wikipedia.org/wiki/Message_queue) - feel free to pick any way that suits your case best

## Installation and usage

Since the meaning of the project is to implement as easy to use production ready model wrapper as possible, installation process is simple too. It uses [Docker](https://www.docker.com/) so it's required to be installed to build and run the project:

1. Clone the project:

```sh
git clone https://github.com/Dominux/commercial-studio-photos-generator.git
```

2. Copy `.env.example` to `.env`:

```sh
cp ./.env.example ./.env
```

and you can optionally edit it as you wish

3. Build and run the project:

```sh
make up
```

Once it ran you can access [localhost:8000](http://localhost:8000) to get the webpage.

If something went wrong or you're just willing to check the logs, you can do that via looking into logs of `cspg-server` and `cspg-worker` services:

```sh
docker logs -f cspg-server
```

and

```sh
docker logs -f cspg-worker
```

## Architecture

It uses [Web-Queue-Worker](https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/web-queue-worker) pattern by obvious reasons. The whole architecture is presented below:

![image](https://github.com/Dominux/commercial-studio-photos-generator/assets/55978340/c5d5f601-c609-480e-b345-6608aecfa86c)

With scalability being extremelly necessery in case of such applications, this architecture brings ability to simply scale:

- inner traffic by increasing Server workers amount
- amount of messages to be queued by... increasing Queue capacity
- model perfomance by increasing Worker workers amount

This architecture also allows to distribute the whole system between many machines, like taking all the services into different machines or even increase generational perfomance by using different GPUs by using a separate instance of Worker worker on each GPU unit.

To make able to use such a distributed system, the app provides a centralized service to store file objects ([MinIO](https://min.io/)) and a cache system ([Redis](https://redis.io/)) to speed up checking for done results that also increases the system's overall perfomance.

### Generation acceleration

As was said before, to speed up generation process, you need to examine the GPU usage first. If it doesn't use it on 100%, increasing amount of parallel generation processes can increase total perfomance. Nevertheless better to increase usage of a single GPU up to 100% because it won't change overall perfomance, but it will speed up generating process when you have only one user at a time, at least. In this case increasing workers amount on a single GPU won't result in any acceleration since you already uses it on its maximum, but will drammatically multiply RAM and VRAM usage, I'm pretty sure no one interested in such a waste of resources for nothing.

To demonstrate how the perfomance stays the same with increasing workers amount, I performed tests on running the same process of parallel 10 requests on generating the same product. Even though workers were really taking multiple messages at the same time, their speed was multiply slowed down. So here's the chart of running this test for 1-4 workers (since I got only 16GB RAM and 12GB of VRAM on my home machine and renting VDS even with same GPU is so expensive):

![image](https://github.com/Dominux/commercial-studio-photos-generator/assets/55978340/b4d588ee-a712-4286-a9b0-1158b5e5bd61)

So, as I said before, in this case increasing perfomance should be performed by increasing amount of GPUs. [Direct RabbitMQ Exchange](https://www.rabbitmq.com/tutorials/amqp-concepts.html#exchange-direct) I use provides ability to fan out messages from a queue to multiple consumers so it allows such a way of generation acceleration.

### Tech stack

- Frontend
  - [HTMX](https://htmx.org/) - I gave it a try and it perfoms well, especially in such a simple cases
  - [Materialize](https://materializecss.com/) - CSS library cause I try to avoid using JS with HTMX

- Backend
  - [Go](https://go.dev/) - its simple intuitive syntax and perfomance make it the best choice when it comes to create a fast microservices fast, with multiple connections to different systems
  - [Python](https://www.python.org/) - the language with the easiest prototyping [ML](https://en.wikipedia.org/wiki/Machine_learning) way
    - [FastAPI](https://fastapi.tiangolo.com/) - async backend web framework for fast prototying. TBH, in this case I use it to provide healthcheck functionallity and many other stuff in a simple way
    - [Gunicorn](https://gunicorn.org/) - [WSGI](https://ru.wikipedia.org/wiki/WSGI) HTTP server, I use it for a simple workers amount control since it uses multiprocess based workers allocation
  - [RabbitMQ](https://www.rabbitmq.com/) - message broker store messages and to push them to their consumers
  - [MinIO](https://min.io/) - an object storage. I use it for centralized storing generation results
  - [Redis](https://redis.io/) - an in-memory database. I use it for storing results IDs for the server to check if the result is already done

### Client-server communication

I would use [websockets](https://en.wikipedia.org/wiki/WebSocket) or [SSE](https://en.wikipedia.org/wiki/Server-sent_events) for that, but I used a good ol' polling (even not a long polling) and here's why:

- Generation process takes not nanoseconds or even milliseconds to perform. It can take seconds or even minutes for just a single task to be done. So realtime communication isn't a thing here at all
- WS and SSE require to keep alive stateful connections. And with increasing parallel users and the slowing the process down it starts to take RAM and CPU usage so wild. So in terms of handling as much as possible users at the same time, simple polling strategy is the best choice cause it's easy to adjust the interval and the user won't need to keep polling the server once they will get the result. WS and SSE is rather for complext cases afterall

Even though it's only a choice of mine and you can feel free to pick the way you wish. Afterall it wasn't the project's main point.
