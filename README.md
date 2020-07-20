# rv

A simple ranked voting app.

## layout

Go + postgresql web application deployed with Heroku. Models are generated code created with [sqlboiler](https://github.com/volatiletech/sqlboiler).

Purposely avoids any concept of user or authentication; may later add sessions if necessary.

## local dev

Local dev configuration is located within `.env` (consider using [direnv](https://direnv.net/) -- the [.envrc](.envrc) exists to immediately load .env). Most local dev commands will expect these in the environment -- especially, `DATABASE_URL`.

The `DEBUG` env var enables local dev niceties (e.g. hot reload of templates).

To run locally, tools expect:

- docker
- docker-compose

The makefile defines most relevant commands

- `make test` runs tests
- `make dbup` starts a database container on port 5432, creating a database and user within it.
- `make migrate` applies all migrations to the database container
- `make gen` regenerates the database models.

### common tasks

#### creating a new model

- Write a new migration in [migrations](migrations).
- Generate updated models with `make gen`.
- Program away.

#### adding a new route / view

- Register the route in [app.go](app.go).
- Define the method on handling method on handler.
- Define the templates in [templates](templates).
