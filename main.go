package main

import "fmt"

func main() {
	transactions := []float64{}
	for {
		transaction := inputUser()
		if transaction == 0 {
			break
		} else {
			transactions = append(transactions, transaction)
			sum := getSum(transactions)
			fmt.Printf("общий баланс:%.2f \n", sum)
		}
	}
}

func inputUser () (float64) {
	var valuesUser float64
	fmt.Println("Введите сумму, либо N, чтобы закрыть программу")
	fmt.Scan(&valuesUser)
	return valuesUser
}

func getSum (transactions []float64) (float64) {
	var sum float64
	for _, value := range transactions {
		sum += value
	}
	return sum
}