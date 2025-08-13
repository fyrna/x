# x/textarea

a simple way to create an opinioned text input

> textarea is my personal library and also an experimental package for my own uses.

## use:

```go
package main

import (
  "fmt"
  "strings"

  "github.com/fyrna/x/textarea"
)

func main() {
  editor := textarea.NewInput("Title", "Body")

  lines, err := editor.Run()
  if err != nil {
    fmt.Printf("\nAborted: %v\n", err)
    return
  }

  fmt.Println("\nResult:")
  fmt.Println(strings.Join(lines, "\n"))
}
```
