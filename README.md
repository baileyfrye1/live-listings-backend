# Live Listings Server

Go server to live listings dashboard. This is a project to practice building larger APIs in Go and to learn WebSockets. Real estate agents will be able to log in and post any active listings and receive real-time feedback from uesr interaction such as favoriting, views, and messages

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make dev
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
