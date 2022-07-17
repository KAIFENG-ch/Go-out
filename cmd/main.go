package main

import (
	gout "Go-out"
	"fmt"
)

func main() {
	r := gout.Default()
	r.GET("/", func(context *gout.Context) {
		fmt.Println("hello world")
	})
	r.Run(":8080")
}
