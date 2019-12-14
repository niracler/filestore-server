package main

import (
	"filestore/cache/myredis"
	"fmt"
)

func main() {
	a := myredis.RedisPool()
	fmt.Println(a)
}
