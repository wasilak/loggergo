# LoggerGo Project Overview

## Purpose
LoggerGo is a lightweight, customizable logging library for Go that provides flexible and simple logging capabilities with OpenTelemetry integration.

## Key Features
- Simple API for various log levels (Info, Debug, Error)
- Customizable log formats and outputs
- OpenTelemetry support (spanID & traceID injection)
- Built on top of Go's standard log/slog package
- Multiple output modes (console, OTEL, fanout)
- Context-aware logging with automatic value extraction
- Dynamic log level adjustment

## Tech Stack
- Go 1.24.0+
- log/slog (Go standard library)
- OpenTelemetry (go.opentelemetry.io)
- Property-based testing with gopter
- Multiple dev flavors: tint, slogor, devslog

## Project Structure
```
loggergo/
├── lib/                    # Internal implementation
│   ├── config.go          # Configuration management
│   ├── modes/             # Output mode implementations
│   ├── outputs/           # Format implementations
│   └── types/             # Type definitions
├── examples/              # Usage examples
├── logger.go              # Main API
├── context_handler.go     # Context extraction
└── wrapper.go             # Type exports
```

## Current Status
The project is undergoing a comprehensive audit and improvement process to transform it into a production-ready, enterprise-grade logging solution. Tasks 1-10 have been completed, focusing on configuration validation, thread safety, error handling, and property-based testing.