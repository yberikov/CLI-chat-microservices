[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/kF9GUL1O)
# 2024-spring-AB-Go-HW-3-

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