package main

import "fmt"

func main() {
	CalculateLuhn(4561261212345467)
	fmt.Println(Valid(4561261212345466))

} // CalculateLuhn return the check number
func CalculateLuhn(number int) int {
	checkNumber := checksum(number)
	fmt.Println(checkNumber)

	if checkNumber == 0 {
		return 0
	}
	return 10 - checkNumber
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

// Valid check number is valid or not based on Luhn algorithm
func Valid(number int) bool {
	return (number%10+checksum2(number/10))%10 == 0
}

func checksum2(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
