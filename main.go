package main

import (
	"fmt"
	"net/url"
)


func main() {
	res, _ := url.ParseQuery("title_eq=%22Game%20Developer%22")

	fmt.Println(res)
}