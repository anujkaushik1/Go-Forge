package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("kaushik error = ", err)
	}

	input = strings.TrimSpace(input)
	integerInput, err := strconv.ParseInt(input, 10, 64)

	if err != nil {
		fmt.Println("kauhsik er22 = ", err)
		return
	}

	// fmt.Println(integerInput)

	newResult := integerInput + 10

	fmt.Println("new res = ", newResult)

}
