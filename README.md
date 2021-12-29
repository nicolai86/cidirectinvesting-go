# cidirectinvesting-go

[![Go Reference](https://pkg.go.dev/badge/github.com/nicolai86/cidirectinvesting-go.svg)](https://pkg.go.dev/github.com/nicolai86/cidirectinvesting-go)
[![go](https://github.com/nicolai86/cidirectinvesting-go/actions/workflows/go.yml/badge.svg)](https://github.com/nicolai86/cidirectinvesting-go/actions/workflows/go.yml)

This is a Go SDK for the [CI Direct Investing](https://cidirectinvesting.com/) application.
It's a tiny wrapper around the rest API exposed by CI Direct Investing. 
I couldn't find an official SDK, so here's just what I need to integrate it into my home system.

## Running all tests locally

```bash
CDI_KEY_ID=<access key> CDI_SECRET_KEY=<secret key> go test ./... -v
```
