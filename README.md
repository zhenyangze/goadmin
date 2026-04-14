# GoAdmin

A Dcat Admin-inspired admin framework for Go.

## Features

- **Resource-oriented DSL**: Grid, Form, Show, Tree builders
- **25 Grid Displayers**: Badge, Label, Image, ProgressBar, QRCode, Switch, Checkbox, Radio, Select, Modal, Table, Tree, etc.
- **22 Widgets**: Checkbox, Radio, Table, LazyTable, Terminal, Chart, Tab, etc.
- **Auth & RBAC**: Built-in authentication, role-based access control, audit logging
- **Form Fields**: 50+ field types including upload, repeater, editor
- **Repository Hooks**: Automatic audit logging for create/update/delete
- **net/http First**: Can be mounted in Gin, Echo, Chi, or plain stdlib

## Installation

```bash
go get github.com/zhenyangze/goadmin
```

## Quick Start

```go
package main

import (
    "net/http"
    "github.com/zhenyangze/goadmin"
)

func main() {
    app := goadmin.New(goadmin.Config{
        Title: "My Admin",
        SessionSecret: "your-secret-key",
    })

    // Register resources...

    http.ListenAndServe(":8080", app.Handler())
}
```

## Documentation

See [go-admin-build](https://github.com/zhenyangze/go-admin-build) for full documentation and demo.

## License

MIT
