[![Golang](https://img.shields.io/badge/Go-v1.21-EEEEEE?logo=go&logoColor=white&labelColor=00ADD8)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

<div align="center">
    <h1>Сервис для поздравления с днем рождения</h1>
    <h5>
Сервис, написанный на golang для поздравления именинников и уведомлении о предстоящих днях рождениях. 
    </h5>
    <p>
        <a href="README.md">English</a> | Russian
    </p>
</div>

- - -

## Используемые технологии:
- [Golang](https://go.dev), [PostgreSQL](https://www.postgresql.org/), [Docker](https://www.docker.com/), [REST](https://ru.wikipedia.org/wiki/REST)

- - -
## Навигация

- - -
## Загрузка
```shell
git clone git@github.com:Nol1feee/birthday-notifier.git
```

## Начало работы
1.**Настройка .env файла и, при необходимости, внесение правок в config.yaml**

### .env
```shell
DB_HOST=localhost
DB_PORT=5432
DB_NAME=birthdayNotifier
DB_USER=rutube
DB_PASSWORD=qwerty123
DB_SSLMODE=disable

DB_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

#email pass - пароль для подключение к smtp серверу
#по ссылке ниже вы можете найти инстукцию для генерации пароля
#https://mailmeteor.com/blog/gmail-smtp-settings#how-to-use-the-gmail-smtp-settings
EMAIL_PASS= 
````

### config.yaml
```shell
http:
  host: "localhost"
  port: 8080

email:
  host: "smtp.gmail.com" #стандартный хост, если используете google
  port: 587              #стандартный порт для smtp
  from:                  #email адрес, с которого будут отправлены уведомления

log_mode: "dev"          #dev || prod
```
2.**запустите docker daemon** <br>
3.**выполните команду 'make'**
## Примеры запросов
> подробная документация доступна по следующей [ссылке](https://documenter.getpostman.com/view/27531074/2sA3XLDiRv)
