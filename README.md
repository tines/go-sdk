# Tines go-sdk

> **Note**: This library is under active development as we work towards full coverage of the Tines API. 
> Although we do our best to minimize any breaking changes during development, v0.x may include breaking
> changes in minor versions, although patch versions will only include backwards-compatible changes and 
> bug fixes. Once we reach a stable v1.x version, breaking changes will only be introduced in new major
> versions according to SemVer principles.

This is the official Go library for interacting with the [Tines API](https://www.tines.com/api/). This 
library allows you to do things like:

- Create, read, update, list, and delete Stories
- Manage teams, users, credentials, and resources
- Administer tenant settings

## Installation

You need a working Go environment. Although the Tines GO SDK may work with multiple versions of Go, we officially
support only currently-supported Go versions according to [Go project's release policy](https://go.dev/doc/devel/release#policy),
and compatibility with other Go versions is not guaranteed.

```
go get github.com/tines/go-sdk
```

## Getting Started

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/tines/go-sdk/tines"
)

func main() {
    // Construct a new API client
    cli, err := tines.NewClient(
        tines.SetTenantUrl(os.Getenv("TINES_TENANT_URL")),
        tines.SetApiKey(os.Getenv("TINES_API_KEY")),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Set a Context for making API calls
    ctx := context.Background()

    // Fetch data about Team ID #1 
    t, err := cli.GetTeam(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }

    // Print the Team name
    fmt.Println(t.Name)
}

```

Naming conventions for this package follow a `{Verb}{Object(s)}` pattern mirroring the actions and objects outlined
in the [API Documentation](https://www.tines.com/api/). For example, retrieving an individual Folder is `GetFolder()`,
updating an individual Folder is `UpdateFolder()`, listing all Folders is `ListFolders()`, etc. 

The Tines SDK supports [Uber's zap logging library](https://github.com/uber-go/zap/) for debugging purposes. To enable logging, 
pass in a configured logger when creating a new client. The Tines SDK only logs at a debug level - any errors will be passed
back to your application and should be handled according to your normal error-handling logic.

```go
logger := zap.Must(zap.NewDevelopment())
defer logger.Sync()

cli, err := tines.NewClient(
    tines.SetTenantUrl(os.Getenv("TINES_TENANT_URL")),
    tines.SetApiKey(os.Getenv("TINES_API_KEY")),
    tines.SetLogger(logger)
)
```

## Contributing

Pull Requests are welcome, but please open an issue (or comment in an existing issue) to discuss any non-trivial 
changes before submitting code.
