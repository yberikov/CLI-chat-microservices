
# CLI Chat

## How to get started
1. Run docker-compose
```bash
docker-compose up --build
```
2. Run migrations
```bash
cd ./storage-microservice
make migrate-up
```
3. Run client
```bash
cd ./chat-microservice
make runClient
```

## Application structure

![img.png](img.png)
