package main

import "godis/local"

func main() {
	c := local.New()
	c.SetMaxMemory("1MB")
}
