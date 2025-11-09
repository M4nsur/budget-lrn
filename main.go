package main

import "fmt"

func main() {
	transactions := []float64{}
	for  {
		transaction := inputUser()
		if transaction == 0 {
			break
		} else {
			transactions = append(transactions, transaction)
			fmt.Println(transactions)
		}
	}
}

func inputUser () (float64) {
	var valuesUser float64
	fmt.Println("Введите сумму, либо 0, чтобы закрыть программу")
	fmt.Scan(&valuesUser)
	return valuesUser
}
