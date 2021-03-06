# goth

A tiny go auth micro-framework, with support for JWT access and refresh tokens. 🔐

## Available HTTP APIs

**auth endpoints:**

- POST `/api/v1/auth/signup` : signup
- POST `/api/v1/auth/login` : login
- POST `/api/v1/auth/logout` : logout
- POST `/api/v1/auth/refresh` : refresh tokens
- POST `/api/v1/auth/reset/init` : init password reset
- POST `/api/v1/auth/reset/verify` : verify password reset
- DELETE `/api/v1/auth/delete` : 🛡 delete account

**user endpoints:**

- GET `/api/v1/user` : 🛡 get user info
- PUT `/api/v1/user` : 🛡 update user info

🛡: protected route i.e. requires valid bearer token `Authorization` header

## Getting started

### requirements

- golang 1.17+
- docker 20.10+
- docker compose v2+

### setting up env

and a `.env` with the following contents:

```sh
GOTH_LISTEN_HOST=0.0.0.0
GOTH_LISTEN_PORT=3333

REDIS_HOST=localhost
REDIS_PORT=6379

DB_HOST=localhost
DB_PORT=27017
DB_NAME=gothDb
DB_USER=root
DB_PASSWORD=secret

ACCESS_TOKEN_MAX_AGE_IN_SECONDS=3600
REFRESH_TOKEN_MAX_AGE_IN_SECONDS=1296000
ACCESS_TOKEN_SIGNING_KEY=your-access-signing-key
REFRESH_TOKEN_SIGNING_KEY=your-refresh-signing-key

SENDGRID_API_KEY=your-sendgrid-api-key
FROM_EMAIL_ADDRESS=verified-sendgrid-sender@example.com
```

copy the same file and name that `.prod.env`, with these two variables updated:

```sh
REDIS_HOST=redis
DB_HOST=mongodb
```

add a `.mongo.env` with the following contents:

```sh
MONGO_INITDB_DATABASE=gothDb
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=secret
```

NOTE:

- update the `DB_PASSWORD` (& `MONGO_INITDB_ROOT_PASSWORD`), `ACCESS_TOKEN_SIGNING_KEY` and `REFRESH_TOKEN_SIGNING_KEY` to something more secure
- when db credentials are updated, make sure you sync them across all the `*.env` files
- to be able to send emails for password reset successfully, you need to
  - [create a sendgrid account](https://sendgrid.com/)
  - [create an API key](https://app.sendgrid.com/settings/api_keys) and update env `SENDGRID_API_KEY`
  - [create a sender](https://app.sendgrid.com/settings/sender_auth/senders/new) and update env `FROM_EMAIL_ADDRESS`

### Usage in Docker

**to spin up the whole project, just run:**

```sh
docker compose up -d
```

it might take a few minutes to download and build the images the first time. grab a cup of tea perhaps... ☕️

**to bring down the project, run:**

```sh
docker compose down
```

add ` -v` after `down ` if you also want to remove the volumes associated with mongodb and redis.

**to rebuild the image after you've made any changes to the code, run:**

```sh
docker compose build
```

and then run using `up -d` as usual to run the project.

### Usage in Development

**to run the goth project directly on your terminal using `go` and not using docker, you need to run mongodb and redis locally.**

that can be done very easily with docker:

```sh
docker compose up -d mongodb redis
```

**then, to run the project, from the root execute:**

```sh
go run .
```

## License

MIT.
