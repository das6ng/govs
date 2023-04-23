# Introduction

This is a simple go version switcher.

# Installation

Just run: `go install github.com/dashengyeah/govs@latest`

# Usage

1. First, you should specify your `GOROOT`'s parent directory in a environment variable `GO_ROOT_PARENT_DIR`

for example, `export GO_ROOT_PARENT_DIR="$HOME/.local/"`

2. And then, prepare a your go sdk in `GO_ROOT_PARENT_DIR` MANUALLY.

your `GO_ROOT_PARENT_DIR` directory should looks like this:

```plaintext
GO_ROOT_PARENT_DIR
|
|-go/            --> your GOROOT should points to this dir
|-go1.18.10/     --> available go sdks
|-go1.19/
|-go1.20.3/
`
```
