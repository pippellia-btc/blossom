# Blossom

This is a minimal, zero-dependency library that provides the core data structures and utility functions for working with [Blossom](https://github.com/hzrd149/blossom). It intentionally does not include HTTP clients, servers, or storage backends, just the building blocks so others can build higher-level libraries on top without inheriting unwanted dependencies.


[![Go Report Card](https://goreportcard.com/badge/github.com/pippellia-btc/blossom)](https://goreportcard.com/report/github.com/pippellia-btc/blossom)
[![Go Reference](https://pkg.go.dev/badge/github.com/pippellia-btc/blossom.svg)](https://pkg.go.dev/github.com/pippellia-btc/blossom)

## Installation

```bash
go get github.com/pippellia-btc/blossom
```

## Features

- Zero dependencies: standard library only
- `Blob`: interface for files, byte slices, and streams with automatic MIME type detection
- `BlobDescriptor`: blob metadata with automatic handling of extra fields
- `Hash`: 32-byte SHA-256 with JSON and SQL serialization built-in
- `Error`: protocol-compliant HTTP errors with `X-Reason` header support