# Project Template

## Dependencies

- GO
- NODE 21
- PNPM

## Setup

### Env variables

Copy `.env.example` and update the variables to your desire

```bash
cp .env.example .env
```

### Docker compose (optional)

If you want to use docker to run the database you can use the provided
docker-compose file.

```bash
cp docker-compose.example.yml docker-compose.yml
```

### Install server dependencies

```bash
make build
```

### Update the database

Run the migrations to prepare the database

```bash
make migrations
```

### Run server

```bash
make server
```

### Run app

```bash
make frontend
```
