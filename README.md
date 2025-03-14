## Сервис домов

# Запуск   
БД запускается в docker-compose оболочке, настройки парсятся из .env файла.
Запуск докер-образа, обновление пакетов, миграции в одной команде 

`
make run-all
`

Тесты 

`
make test
`

Другие команды можно найти в Makefile.

# Настройки   
Настройки сервера (хост, порт), а также БД хранятся в .env файле в корне каталога. 

# Создание таблиц БД   
Для создания таблиц в БД используются миграции (db/migrations)   

Применение   
`
make migrateup
`

Отклонение   
`
make migratedown
`

# Идентификаторы квартир   
Не до конца разобрался в условии, поэтому в качестве id в запросах с квартирами используются сгенерированные id из БД.

# Изменение времени обновления дома   
Изменение updated_at для дома происходит с помощью триггера при добавлении или изменении квартиры, связанной с домой внешним ключем (flats.house_id = houses.id). Также при получении списка квартир для дома по его id в структуру вывода были добавлены поля дома (id, address, developer, year, created_at, updated_at), чтобы было видно, как изменяется поле updated_at.

# Обновление дома   
В случае, если модератором не был передан новый статус дома (такая возможность, в соответствии с API, существует) по умолчанию ставится статус `approved`.

# Регистрация   
По ручке /dummyLogin выдаются 2 токена по умолчанию, которые прописаны руками в коде.
Реализовано дополнительное задание по регистрации и логину новых пользователей, используются дополнительные таблицы users и tokens. В качестве токенов используются JWT-токены (генерация в пакете pkg/jwt), ключ лежит в файле .secret. Введенные пароли хэшируются с помощью md5 функции (pkg/crypt).

# Подписка   
Подписка на уведомления по добавлению квартир в доме реализовано в структуре SubscriberService (internal/service), вызов функции происходит асинхронно после успешного добавления квартиры в существующий дом. Для хранения списка подписчиков и соответствующих ему домов создана дополнительная таблица subscribers.

# Логирование   
В сервисе производится логирование запросов, функций и т.д. с помощью логгера slog, добавленного в ядро Go.
