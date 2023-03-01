# Клиент серверное приложение Playlist

Клиентская и серверная часть написаны на язке go
Для хранения данных с серверной стороне используется двусвязный список

## Протокол взаимодействия клиентского приложения и серверной части

Взаимодействие просиходит по протоколу gRPC

## Файл proto

```
syntax = "proto3";

option go_package = "playlist/pb";

package playlist;

// Interface exported by the server.
service PlaylistService {
  // A simple RPC.
  rpc SendCommand(Command) returns (Response) {}

}

message Command {
	string com = 1;
	string name = 2;
	int32 duration = 3;
}

message Response {
	string data =1;
}
```
## Формат сообщения
- com - имя команды
- name - название пести
- duration - продолжительность песни в секундах

## Поддерживаемые команды

- add : добавить песню. Обязательно заполнение поля name и duration
- delete : удалить песню. Обязательно заполнение поля name
- play : Запускает исполнение плейлиста на стороне сервера
- next : переход к следующей песне
- previos : переход к предыдущей песне
- print : вывод плейлиста в стандартный вывод на сервере
# Response
Сообщением  Response серверная часть подтверждает успешное получение команды

## Работа приложения
### Серверная часть 
Серверная часть запускается из каталога `server` командой 
`go run main.go`
### Клиентская часть
Клиентская часть запускается из каталога `client` командой 
`go run clent.go`
### Сетевой порт
Приложение для работы использует сетевой порт `9000` и протокол 'tcp'

### Работа с докером
#### докер файл
```
# syntax=docker/dockerfile:1

FROM golang:1.20.1-alpine3.17

WORKDIR /app

COPY server/go.mod ./
COPY server/go.sum ./

RUN go mod download

COPY server/*.go ./
COPY server/pb ./pb

RUN go build -o main

EXPOSE 9000

CMD [ "./main" ]
```
### Запуск серверной части в докере
Запуск лучше производить с подключенной консолью
`docker run -p 9000:9000 --rm -it имя_контенер sh`
## Unit тесты
Unit тесты не реализованы

