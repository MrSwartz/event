## О сервисе
Сервис поднимает HTTP сервер, в котором который реализует эндпоинт обрабатывающий запросы позволяет сохранять тело запроса в ClickHouse.
## Тесты
Есть небольшое тестовое покрытие, в основном это /internal..., pkg/eventservice/service/data и функции не работающие с бд из /pkg...

в идеале нужно написать интеграционные тесты на эндпоинты, но <s>у меня сейчас нет времени даже на сон</s> я проверил через постман и ab

## Пример использования
Для старта и тестов нужно выпонить:
```
export ENV=dev/test
export CLICKHOUSE_NAME=default            
export CLICKHOUSE_HOST=127.0.0.1
export CLICKHOUSE_PASSWORD=qwerty123
export CLICKHOUSE_PORT=9000
export CLICKHOUSE_USER=default
```
## Результаты тестов
При тестировании была поднята база в докере, а сервис был запущен нативно на MacOS 12.x.x M1+16GB Ram
в момент тестирования ос слегка перегружена и 2.5GB лежит в свопе
конфиг буфера:
```
LoopTimeout=0
Size=0
```
запуск нагрузки:
```
ab -n 10000 -p payload.txt http://localhost:8080/v1/events
```
результат нагрузочного тестирования:
```
Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /v1/events
Document Length:        30 bytes

Concurrency Level:      1
Time taken for tests:   80.375 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1470000 bytes
Total body sent:        68010000
HTML transferred:       300000 bytes
Requests per second:    124.42 [#/sec] (mean)
Time per request:       8.038 [ms] (mean)
Time per request:       8.038 [ms] (mean, across all concurrent requests)
Transfer rate:          17.86 [Kbytes/sec] received
                        826.33 kb/s sent
                        844.19 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       2
Processing:     4    8   6.5      6     144
Waiting:        4    8   6.5      6     144
Total:          4    8   6.5      6     145

Percentage of the requests served within a certain time (ms)
  50%      6
  66%      7
  75%      8
  80%      9
  90%     12
  95%     15
  98%     23
  99%     39
 100%    145 (longest request)
```

С включеным буфером удалось получить результат получше

конфиг буфера:
```
[Buffer]
LoopTimeout=6
Size=10000
```
запуск нагрузки:
```
ab -n 15000 -p payload.txt http://localhost:8080/v1/events
```
результат нагрузочного тестирования:
```
Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /v1/events
Document Length:        30 bytes

Concurrency Level:      1
Time taken for tests:   8.564 seconds
Complete requests:      15000
Failed requests:        0
Total transferred:      2205000 bytes
Total body sent:        102015000
HTML transferred:       450000 bytes
Requests per second:    1751.49 [#/sec] (mean)
Time per request:       0.571 [ms] (mean)
Time per request:       0.571 [ms] (mean, across all concurrent requests)
Transfer rate:          251.44 [Kbytes/sec] received
                        11632.73 kb/s sent
                        11884.16 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.3      0      25
Processing:     0    0   1.6      0     114
Waiting:        0    0   1.5      0     114
Total:          0    1   1.6      0     114

Percentage of the requests served within a certain time (ms)
  50%      0
  66%      0
  75%      1
  80%      1
  90%      1
  95%      1
  98%      2
  99%      3
 100%    114 (longest request)
```

Играясь с настройками буфера можно подобрать более эффективные значения/нагрузку и получить прирост к скорости на несколько(а может и пару десятков) процентов