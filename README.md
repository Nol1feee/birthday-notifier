[![Golang](https://img.shields.io/badge/Go-v1.21-EEEEEE?logo=go&logoColor=white&labelColor=00ADD8)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

<div align="center">
    <h1>Birthday Greeting Service</h1>
    <h5>
        A service written in Golang for congratulating birthday celebrants and notifying about upcoming birthdays.
    </h5>
    <p>
        English | <a href="README.ru.md">Russian</a>
    </p>
</div>

- - -

## Technologies Used:
- [Golang](https://go.dev), [PostgreSQL](https://www.postgresql.org/), [Docker](https://www.docker.com/), [REST](https://en.wikipedia.org/wiki/Representational_state_transfer)

- - -
## Navigation

- - -
## Installation
```shell
git clone git@github.com/Nol1feee/birthday-notifier.git
```
- - - 
## Getting Started
1. **Set up the .env file and make necessary adjustments to config.yaml**
### .env
```shell
DB_HOST=localhost
DB_PORT=5432
DB_NAME=birthdayNotifier
DB_USER=rutube
DB_PASSWORD=qwerty123
DB_SSLMODE=disable

DB_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# Email password for SMTP server connection
# You can find instructions for generating a password at the link below:
# https://mailmeteor.com/blog/gmail-smtp-settings#how-to-use-the-gmail-smtp-settings
EMAIL_PASS= 
```

### config.yaml
```shell
http:
  host: "localhost"
  port: 8080

email:
  host: "smtp.gmail.com" # Standard host if using Google
  port: 587              # Standard port for SMTP
  from:                  # Email address from which notifications will be sent

log_mode: "dev"          # dev || prod
```
2.**Run Docker daemon** <br>
3.**Execute the `make` command**
- - -
## Example Requests
Detailed documentation is available at the following [link](https://documenter.getpostman.com/view/27531074/2sA3XLDiRv)
