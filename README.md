# EffectiveMobile
 
 ![Build Status](https://github.com/QuickSilver-1/EffectiveMobile/actions/workflows/go.yml/badge.svg)

 <h3>API для работы с музыкальной библиотекой в рамках технического задания для отбора на стажировку</h3>

 <h2>Структура проекта</h2>

<code>
/EffectiveMobile
│
├── /api
│   └── /swagger.yml 
├── /cmd
│   └── /music
│       └── main.go - ОСНОВНОЙ ФАЙЛ
├── /internal
│   └── /config 
│   |   └── config.go - файл с конфигами
│   └── /crud  - пакет с функциями для взаимодействия с PostgreSQL
│   |   └── collectCRUD.go - сбор CRUD в единую коллекцию
│   |   └── musicCRUD.go
│   |   └── struct.go
│   |   └── validation.go
│   └── /handlers - пакет с обработчиками входящий запросов
│   |   └── handler.go
│   └── /migrations - миграции БД
│   |   └── ...
│   └── /serverbuilder - подключение хэндлеров к веб-серверу
│       └── serverbuilder.go
├── /log
│   └── log.log - файл с логами
├── /pkg
│   └── /db - пакет с подключением к СУБД
|   |   └── connect.go - подключение к основной базе (Postgre)
|   |   └── migrate.go - запуск миграций при старте сервера
|   |   └── redis.go - подключение к базе для кэширования (Redis)
│   └── /log - пакет с конфигурацией логгера
│   |   └── logger.go
|   └── /server - пакет с конфигурацей веб-сервера
│       └── httpserver.go - сервер
|       └── middleware.go - мидлвары
├── config.env - КОНФИГУРАЦИОННЫЙ ФАЙЛ приложения
├── Dockerfile - файл создания образа докер
├── go.mod
└── go.sum
</code>

<h2>Запуск</h2>

Запустить сервер можно одним из двух способов
<ul>
<li>Gпрописать команду <code>go run cmd/music/main.go</code>, находясь в главной директории</li>
<li>Собрать и запустить докер контейнер из Dockerfile</li>
<li>Также можно потестировать сервис через <code>http://89.46.131.181:8080</code> (но здесь нет вашего обработчика на /info)</li>
</ul>

<h2>Общее описание</h2>
Спасибо за интересную задачу, было интересно делать. Реализовал все необходимые функции, а также дополнительно сделал кэширование через Redis. В качестве основной базы данных использовался PostgreSQL. В пакет <code>pkg</code> я складывал модули, которые в будущем потенциально могу использовать повторно. В <code>internal</code> - всё остальное. В <code>config.env</code> лежат конфиги, которые необходимо поменять на ваши
Написал небольшую защиту от брут форса - ограничение на кол-во запросов в секунду.
<h3>Спецификацию api можно посмотреть в папке api</h3>
