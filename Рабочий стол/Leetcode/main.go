package main

import (
	"fmt"
)

func main() {
	a, b := 9, 9
	fmt.Println(getSum(a, b))
}

func getSum(a int, b int) int { // 9 -- 1001   9 -- 1001
	for b != 0 {
		temp := a & b << 1 // 10010 -- 0000
		a = a ^ b          // 0000 -- 10010
		b = temp           // 10010 -- 0000
	}
	return a
}
