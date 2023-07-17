package main

import (
	"fmt"

	"golang-beginner/internal/config"
)

func main(){
	cfg := config.MustLoad()

	fmt.Println(cfg)
}