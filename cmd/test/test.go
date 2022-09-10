package main

import (
	"fmt"

	"github.com/siestacloud/gopherMart/pkg"
)

func main() {
	if err := pkg.Valid("346436439"); err != nil {
		fmt.Println(err)
	}
}
