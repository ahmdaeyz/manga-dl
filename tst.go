package main

import (
	"fmt"
	"strings"
)

func main() {

	str := "Volume 07"
	str = strings.Replace(str, "Volume ", "", -1)
	fmt.Println(str)
}
