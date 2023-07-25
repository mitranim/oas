## Overview

"oas" is short for "**O**pen**A**PI **S**pec". Go package for generating OpenAPI docs at runtime.

**Non**-features:

  * No CLI.
  * No code parsing.
  * No code generation.
  * No magic comments.
  * No slowness.
  * No dependencies.
  * Nothing added to your HTTP stack.

Features:

  * Struct definitions for OpenAPI 3.1.
  * Uses reflection to make OAS schemas from your types.
    * No more maintaining separate definitions by hand.
    * The source of truth is **your Go types**. Not some external YAML.
    * Examines _actual_ encoding behavior of your types, at runtime, to determine formats and nullability.
    * Supports references and cyclic types.
  * Uses Go structs to describe what can't be reflected (routes, descriptions, etc).
    * Structured, statically-typed format.
    * Not an ad-hoc data format in breakage-prone comments.
    * Not some external YAML.
  * The docs are Go structures. You can do anything with them:
    * Inspect and modify in Go.
    * Encode as JSON or YAML.
    * Write to disk or stdout at build time.
    * Serve to clients at runtime.
    * Visualize using an external tool.
  * Tiny and dependency-free.

See [limitations](#limitations) below.

API docs: https://pkg.go.dev/github.com/mitranim/oas

## Why

* No external CLI. Just a library.
  * Automatically downloaded by Go.
  * Automatically versioned via `go.mod`.
  * No manual CLI installation.
  * No manual CLI versioning.
* No Go code generation.
  * No `make generate` when pulling or switching branches.
  * No waiting 1 minute for that stupidly slow generator.
  * No manual remapping from crappy generated types to the types you **actually** want to use.
  * No crappy 3rd party "middleware" in your HTTP stack.
* No forced generation of OAS files. It's entirely optional.
  * No manual reruns of a generate command.
  * No bloating your commits and Git diffs with generated JSON/YAML.
  * Can generate and serve OAS purely at runtime.
* No figuring out how to deal with built artifacts.
  * If you want built artifacts, it's an option. _Your Go app_ is a CLI tool. Add a command to encode its OAS as JSON, and run that at build time.

## Usage

This example focuses on the OAS docs, registering docs for routes, with schemas from Go types. Routing and server setup is elided.

```golang
import (
  "encoding/json"
  "net/http"

  o "github.com/mitranim/oas"
)

var doc = o.Doc{
  Openapi: o.Ver,
  Info:    &o.Info{Title: `API documentation for my server`},
}

type PageInput struct {
  Limit  uint64 `json:"limit"`
  Offset uint64 `json:"offset"`
}

type PersonPage struct {
  PageHead
  Vals []Person `json:"vals"`
}

type PageHead struct {
  Keys []string `json:"keys"`
  More bool     `json:"more"`
}

type Person struct {
  Id   string `json:"id"`
  Name string `json:"name"`
}

var _ = doc.Route(`/api/persons`, http.MethodGet, o.Op{
  ReqBody: doc.JsonBodyOpt(PageInput{}),
  Resps:   doc.RespsOkJson(PersonPage{}),
  Desc:    `Serves a single page from a person feed, paginated.`,
})

func servePersonFeed(rew http.ResponseWriter, req *http.Request) {
  // Mock implementation for example's sake.
  try(json.NewEncoder(rew).Encode(PersonPage{}))
}

var _ = doc.Route(`/openapi.json`, http.MethodGet, o.Op{
  Resps: doc.RespsOkJson(nil),
  Desc: `
Serves the OpenAPI documentation for this server in JSON format.
The docs' docs are elided from the docs to avoid bloat.
`,
})

func serveDocs(rew http.ResponseWriter, _ *http.Request) {
  try(json.NewEncoder(rew).Encode(&doc))
}

func try(err error) {
  if err != nil {
    panic(err)
  }
}
```

## Limitations

The following features are currently missing, but may be added on demand:

* Describing generic/parametrized types.
  * Current plan: wait for Go generics, which are expected in 1.18 on Feb 2022.
* Router integration.
  * Integration with `github.com/mitranim/rout` is planned.

This package doesn't provide a UI. You're expected to feed the resulting JSON into one of **many** externally-available tools for Swagger UI / OpenAPI UI. Many tools can consume specs from a URL, such as your server's endpoint for serving the spec.

## License

https://unlicense.org

## Misc

I'm receptive to suggestions. If this library _almost_ satisfies you but needs changes, open an issue or chat me up. Contacts: https://mitranim.com/#contacts
