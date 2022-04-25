package main

func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
func SumTail(numbers []int) int {
	sum := 0
	if len(numbers) > 1 {

		for i := 1; i < len(numbers); i++ {
			sum += numbers[i]
		}
		return sum
	}
	return 0
}
func SumAll(numbers ...[]int) []int {
	var sum []int
	for _, number := range numbers {
		sum = append(sum, Sum(number))
	}
	return sum
}
func SumAll2(numbersToSum ...[]int) []int {
	lengthOfNumbers := len(numbersToSum)
	sums := make([]int, lengthOfNumbers)

	for i, numbers := range numbersToSum {
		sums[i] = Sum(numbers)
	}

	return sums
}
func SumAllTails(numbers ...[]int) []int {
	var sum []int
	for _, number := range numbers {
		sum = append(sum, SumTail(number))
	}
	return sum
}
