# cache-disk


refer to [cachego/cache](https://github.com/cachego/cache)

## Install

### go module

Use go module directly import：

```go
import "github.com/cachego/disk"
```

### go get

Use go get to install：

```go
go get github.com/cachego/disk
```

## Demo

[Complete example](https://github.com/cachego/disk/example)

### 1. cache demo

code reference：

```go
package main

import (
	"fmt"

	cache "github.com/cachego/disk"
)

func main() {
	key := "key1"
	c := cache.NewInDiskStrCache(".cache")
	c.Set(key, "cache value", 0)
	v, err := c.Get(key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v) // cache value
	c.Del(key)
	v, err = c.Get(key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v) // nil
	c.Clear()
}
```

output：

```text
cache value
<nil>
````