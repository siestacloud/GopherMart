package main

import (
	"fmt"
	"log"

	"github.com/siestacloud/gopherMart/internal/config"
)

var (
	cfg config.Cfg
)

func main() {

	err := config.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

}
