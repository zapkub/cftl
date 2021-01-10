# cftl
Code First Talk Later 👨🏽‍💻


## Development
- Run preparation script first from root workspace 
```
$ ./scripts/prepare_dev.sh
```

then run 
```
$ go run cmd/frontend/main.go
```



## Migrations

Migrations are managed using
[github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate),
with the [CLI tool](https://github.com/golang-migrate/migrate/tree/master/cli).

If this is your first time using golang-migrate, check out the
[Getting Started guide](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md).

To install the golang-migrate CLI, follow the instructions in the
[migrate CLI README](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md).