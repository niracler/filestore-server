package main

import (
	"filestore-server/cache/myredis"
	"fmt"
)

func main() {
	a := myredis.RedisPool()
	fmt.Println(a)
}
