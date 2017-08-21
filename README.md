# Luna

Luna is a command line tool for microservices in monorepos, trying to make developers life easier.

In case when you have your services in a monorepo, there are a some obstacles needs to solve.

During development: it would be great to build and run your services concurrently, maybe even restart them when files changed.

On CI: detect which services has been changed, and package only those.

## Features

- Build and run multiple services
- Watch for changes and restart services
- Clean up on services stopped
- Detect changes in monorepos
- Package services for deployment

## Install or update

```shell
$ go get -u github.com/moonwalker/luna
```

## Development

```shell
$ go run luna.go
```

## License

MIT
