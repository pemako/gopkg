# envload

`fork from` [lestrrat-go/envload](https://github.com/lestrrat-go/envload.git)


- `example`

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pemako/gopkg/envload"
)

func main() {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	fmt.Printf("FOO = %s\n", os.Getenv("FOO"))
	fmt.Printf("BAR = %s\n", os.Getenv("BAR"))

	loader := envload.New()

	os.Setenv("FOO", "Hello")
	os.Setenv("BAR", "World!")

	fmt.Printf("FOO = %s\n", os.Getenv("FOO"))
	fmt.Printf("BAR = %s\n", os.Getenv("BAR"))

	if err := loader.Restore(); err != nil {
		return
	}

	fmt.Printf("FOO = %s\n", os.Getenv("FOO"))
	fmt.Printf("BAR = %s\n", os.Getenv("BAR"))
}

// Output:
// FOO = foo
// BAR = bar
// FOO = Hello
// BAR = World!
// FOO = foo
// BAR = bar
```
