package main

import (
	"fmt"
)

type objs struct {
	id   int
	name string
}

func main() {
	arr := []objs{objs{id: 1, name: "name1"}, objs{id: 2, name: "name2"}}
	fmt.Println(arr)

	for i, _ := range arr {
		change(&arr[i])

	}
	fmt.Println(arr)

}

func change(item *objs) {
	item.name = "qwe"
}
