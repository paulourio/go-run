# go-urn

URN (Uniform Resource Name) identifier parser compliant with RFC 8141:

    urn:nid:nss[?+resolution][?=query][#fragment]

This module implements as best as possible the specification of URNs from RFC 8141.
See unit tests for further examples.

## Installation

```bash
go get github.com/paulourio/go-urn
```

## Example

```go
package main

import (
    "fmt"
    "net/url"

    "github.com/paulourio/go-urn"
)

func main() {
    var input = "urn:ietf:rfc:8141?+Resolve=http?=foo=bar"

    id, err := urn.Parse(input)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%#v\n", id)
    // Output:
    // &urn.URN{
    //     Scheme:        "urn",
    //     NID:           "ietf",
    //     NSS:           "rfc:8141",
    //     Resolve:       "Resolve=http",
    //     Query:         "foo=bar",
    //     Fragment:      "",
    //     ForceFragment: false,
    //     normalized:    false
    // }

    kv, _ := url.ParseQuery(id.Query)
    fmt.Printf("%#v\n", kv)
    // Output:
    // url.Values{
    //     "foo": []string{"bar"}
    // }
}
```
