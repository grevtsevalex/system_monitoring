# Сервис мониторинга системы.

![example workflow](https://github.com/grevtsevalex/system_monitoring/actions/workflows/tests.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/grevtsevalex/system_monitoring)](https://goreportcard.com/report/github.com/grevtsevalex/system_monitoring)

Команда запуска:
```
./progName --config ./configs/config.toml --port 50001
```
или
```
make run
```

### Описание работы.
Сервис роботает в режиме демона на хосте. В конфигурационном файле можно выбрать метрики, которые нужно собирать.

Сбор метрик осуществляется скаутами, каждый из которых собирает свой тип метрик (`loadAverages`, `cpu`, `tps`) и кладет их в свое хранилище.

При запросе пользователя открывается постоянное соединение с клиентом, в которое раз в `N` секунд отправляется статистика за `M` секунд. При этом коллектор метрик проходится по всем хранилищам скаутов и собирает их в единый снэпшот, который отправляет клиенту.

### Запуск простого клиента
Реализован клиент, для примера работы сервера. Он получает снэпшоты от сервера и выводит отформатированный вывод в `stdout`.
Запустить клиента можно командой:

```
go run ./cmd/client/main.go -port 55555
```
или
```
make run-client
```

### Запуск интеграционных тестов
Интеграционный тест запускает на хостовой машине сервер и клиент в докер-контейнере. Клиент получает данные и проверяет, что они изменяются.
Запустить тесты можно командой:

```
make integration-test
```

[Схема архитектуры](https://drive.google.com/file/d/1g72OyR0tcWNLNYvNxVvma_0FzSxUfRPl/view?usp=sharing)
