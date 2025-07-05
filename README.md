# Users Microservice

![Go Badge](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=fff&style=for-the-badge)
![PostgreSQL Badge](https://img.shields.io/badge/PostgreSQL-4169E1?logo=postgresql&logoColor=fff&style=for-the-badge)

## Development Environment Setup

Prerequisites :

- [Nix](https://github.com/DeterminateSystems/nix-installer)
- [Direnv](https://direnv.net/)

Once you have them installed :

1. Allow `direnv` to automatically land you into the **Nix development shell environment**, whenever you cd into this directory :
    ```shell script
    direnv allow
    ```

## Null safety

We can analyze the sourcecode locally for potential nil panics using [nilaway](https://github.com/uber-go/nilaway), by running :
```shell script
nilaway ./...
```

However I'm not doing it in the CI, since a lot of false positives are reported.

## REFERENCEs

**GoLang** :

- [Go Microservice with Clean Architecture: Application Logging](https://medium.com/@jfeng45/go-microservice-with-clean-architecture-application-logging-b43dc5839bce)

- [How to get stacktraces from errors in Golang with go-errors/errors](https://www.bugsnag.com/blog/go-errors/)

- [A Beautiful Way To Deal With ERRORS in Golang HTTP Handlers](https://www.youtube.com/watch?v=aS1cJfQ-LrQ)

- [Setup pgx query logger](https://gist.github.com/zaydek/91f27cdd35c6240701f81415c3ba7c07)

- [Type conversions in GoLang](https://go.dev/ref/spec#Conversions)

- [Golden config for golangci-lint](https://gist.github.com/maratori/47a4d00457a92aa426dbd48a18776322)

**Security** :

- [How to Correctly Validate Passwords - Most Websites Do It Wrong](https://blog.boot.dev/open-source/how-to-validate-passwords/)
