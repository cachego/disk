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
