# LoggerGo

![GitHub tag (with filter)](https://img.shields.io/github/v/tag/wasilak/loggergo) ![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/wasilak/loggergo/main) [![Go Reference](https://pkg.go.dev/badge/github.com/wasilak/loggergo.svg)](https://pkg.go.dev/github.com/wasilak/loggergo) [![Maintainability](https://api.codeclimate.com/v1/badges/87dcca9e40f33cf221af/maintainability)](https://codeclimate.com/github/wasilak/loggergo/maintainability)

A lightweight, customizable logging library for Go, designed to provide flexible and simple logging capabilities.

![diagram](diagram.png)

## Features

- Simple API for various log levels (Info, Debug, Error).
- Customizable log formats and outputs.
- Lightweight and easy to integrate.
- Supports OTEL (OpenTelemetry) by injecting `spanID` & `traceID` into logs.
- By default registers as `slog.Default()` but this can be change via option

## Installation

To install LoggerGo, run:

```bash
go get github.com/wasilak/loggergo
```

## Usage

Here's a basic example of how to use LoggerGo:

https://github.com/wasilak/loggergo/blob/main/examples/simple/main.go

## Contributing

Contributions are welcome! Please fork the repository, make changes, and submit a pull request.

---

Feel free to adjust the example to better match the specific functionality of your library.
