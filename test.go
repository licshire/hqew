package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	dat, _ := ioutil.ReadFile("hqew_list_file.txt")

	secondDomainSlice := strings.Split(string(dat), "\r\n")
	fmt.Print(secondDomainSlice)
}
