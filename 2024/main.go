package main

import (
	"AoC/2024/utils"
	"fmt"
)

func main() {
	input, err := utils.ReadInput(25)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	Solution25(input)
}
