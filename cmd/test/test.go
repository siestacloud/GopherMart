package main

import (
	"fmt"

	"github.com/siestacloud/gopherMart/pkg"
)

func main() {
	if err := pkg.Valid("456126121234546"); err != nil {
		fmt.Println(err)
	}
}
