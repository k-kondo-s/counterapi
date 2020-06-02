# Counter API

# Prepare Environment

These tools are required:

* docker (19.03.8 or later)
* docker-compose (1.25.5 or later)
* jq

First, build containers. Please run the script `scripts/build.sh`,

```
scripts/build.sh
```

and check these 2 containers are built.

```
docker images  | grep -E "ansible|counterapi"
```

Next, set the environment variable `NGINX_IP` that is the endpoint of this application. If you run the above script on your desctop PC, `NGINX_IP` might be `127.0.0.1`.

```
export NGINX_IP=127.0.0.1
```

That's all preparation. Let's run applications as follow.

# Task 1

Run

```
bash -x task1.sh
```

# Task 2

Run

```
bash -x task2.sh
```

# Task 3

Run

```
bash -x task3.sh
```

# Task 4

Run

```
bash -x task4.sh
```

Note: Run `while :; do curl $NGINX_IP; echo; done`, while change the apps counts as you like, and you can see there is no downtime.

# Task 5

Run

```
bash -x task5.sh
```

# Clean up the environment

```
scripts/setup_api.sh stop
```

# Architecture and Design

* I chose docker-compose as a core environment. This is because, as I see Task 4, I have thought that I am required to let whole system so called declarative behavior. docker-compose fits into this requirement and is able to be run on a desktop PC.
  * Actually, first I had used `minikube` but changed my mind. If people implements this applications on k8s, then they would use `type: Service` in terms of k8s, while, in this case, deploying a bare Nginx container onto k8s seems ridiculous.
* I used Redis as a datastore. I have regarded that the higher performance of request/response are required.
* I used Ansible to dynamically change the proxy passes on Nginx settings.
* I let all responses be JSON formatted according to popular API requirements of public services in the world.
* I implemented a counter calculating `nowTimestamp - startTimestamp + 1`, and Redis is in charge of the deleting counter management using its functionality.
  * This architecture can make counter API applications immutable, which means that these applications can be stateless, so the whole system can be scalable.

# TODO and Bugs

* I couldn't care about logger enough, so the log format of the application is so ugly.
* (I've just found a bug, but I no longer have enough time to deal with it. gave up)
