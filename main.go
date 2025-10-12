package main

import (
	"fmt"

	"github.com/shavits/boot-gator/internal/config"
)

func main() {
	curConfig, err := config.Read()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(curConfig)
	config.SetUser("Shahar")
	curConfig, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(curConfig)

}
