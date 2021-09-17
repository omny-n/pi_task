## Тестовое задание для https://productivityinside.com
- [Go](https://golang.org)
- [gRPC](https://grpc.io)
- [gRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [Docker](https://www.docker.com)

## Описание

CRUD-сервис для работы с пользователями, предоставляющий API, используя HTTP и gRPC. В качестве хранилища данных используется MongoDB

## Установка

```shell
docker-compose -f docker-compose.yml up
```

## Доступные команды

- Список всех пользователей

```shell
curl -X GET 'http://localhost:8080/user'
```

- Добавить пользователя

```shell

curl -X POST 'http://localhost:8080/users' \
-d '{
        "firstname": "Oleg",
        "lastname": "Oleg",
        "age": 42,
        "email": "olegd@oleg.com"
    }
```

- Получить данные пользователе

```shell
curl -X GET 'http://localhost:8080/users/6144871b0fbadfb354eb55aa'
```

- Обновить данные пользователя 

```shell
curl -X PUT 'http://localhost:8080/users/6144871b0fbadfb354eb55aa' -d \
'{
            "firstname": "Ne Oleg",
            "lastname": "Ne Oleg",
            "age": 42,
            "email": "neoleg@neoleg.com"
}'
```

- Удалить пользователя

```shell
curl -X DELETE "http://localhost:8080/users/6144871b0fbadfb354eb55aa"
```
