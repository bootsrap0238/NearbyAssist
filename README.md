# Project nearbyassist

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## How to use

### Note
In order to build the project you'll need (Go)[https://go.dev/doc/install] to be installed in your machine

### Compile with Make

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```

### Compile with docker
```bash
docker compose --file compiler-docker-compose.yml up --build
```

### Run dev server with docker
```bash
docker compose -f dev-docker-compose.yml up
```

to remove all containers for dev server
```bash
docker compose -f dev-docker-compose down
```

### Run dev server without docker 
1. Run mysql db
```bash 
docker compose up
```

2. Run server
- In windows, either run the server executable in terminal or double click the executable

3. Seed the db
- Run the [seeder](https://github.com/BeepLoop/NearbyAssist_seeder/releases/tag/v1.0)

