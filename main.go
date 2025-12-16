package main

import (
	"fmt"
	"payments"
)

func main() {
	fmt.Println("=== Тестирование Payment Module ===\n")
	
	testPayPal()
	fmt.Println()
	
	testCreditCard()
	fmt.Println()
	
	testBitcoin()
	fmt.Println()
	
	testErrors()
	fmt.Println()
	
	testCancelErrors()
}

func testPayPal() {
	fmt.Println("--- Тест PayPal ---")
	
	paypal := payments.PayPal{
		Email:  "user@example.com",
		ApiKey: "secret_key_123",
	}
	
	module, err := payments.NewPaymentModule(paypal)
	if err != nil {
		fmt.Printf("Ошибка создания модуля: %v\n", err)
		return
	}
	
	id, err := module.Pay("Покупка книги", 50)
	if err != nil {
		fmt.Printf("Ошибка оплаты: %v\n", err)
		return
	}
	fmt.Printf("Платеж успешен! ID: %d\n", id)
	
	info, err := module.Info(id)
	if err != nil {
		fmt.Printf("Ошибка получения информации: %v\n", err)
		return
	}
	fmt.Printf("Информация: %s - $%d (отменен: %v)\n", 
		info.Description, info.Usd, info.Cancelled)
	
	err = module.Cancel(id)
	if err != nil {
		fmt.Printf("Ошибка отмены: %v\n", err)
		return
	}
	
	info, _ = module.Info(id)
	fmt.Printf("После отмены - отменен: %v\n", info.Cancelled)
}

func testCreditCard() {
	fmt.Println("--- Тест CreditCard ---")
	
	card := payments.CreditCard{
		Number:     "4532123456789012",
		Cvv:        "123",
		ExpiryDate: "12/25",
	}
	
	module, _ := payments.NewPaymentModule(card)
	
	id1, _ := module.Pay("Подписка", 100)
	id2, _ := module.Pay("Книга", 25)
	id3, _ := module.Pay("Курс", 200)
	
	fmt.Printf("Создано платежей: 3 (ID: %d, %d, %d)\n", id1, id2, id3)
	
	allInfo := module.InfoAll()
	fmt.Printf("Всего операций: %d\n", len(allInfo))
	
	for id, info := range allInfo {
		status := "активен"
		if info.Cancelled {
			status = "отменен"
		}
		fmt.Printf("  ID %d: %s - $%d [%s]\n", 
			id, info.Description, info.Usd, status)
	}
}

func testBitcoin() {
	fmt.Println("--- Тест Bitcoin ---")
	
	bitcoin := payments.Bitcoin{
		WalletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
	}
	
	module, _ := payments.NewPaymentModule(bitcoin)
	
	id, err := module.Pay("Донат", 30)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	
	fmt.Printf("Bitcoin платеж успешен! ID: %d\n", id)
}

func testErrors() {
	fmt.Println("--- Тест обработки ошибок ---")
	
	paypal := payments.PayPal{
		Email:  "test@example.com",
		ApiKey: "key123",
	}
	module, _ := payments.NewPaymentModule(paypal)
	
	_, err := module.Pay("Test", 0)
	if err != nil {
		fmt.Printf("Отрицательная сумма: %v\n", err)
	}
	
	_, err = module.Pay("Test", 150000)
	if err != nil {
		fmt.Printf("Слишком большая сумма: %v\n", err)
	}
	
	_, err = module.Pay("", 100)
	if err != nil {
		fmt.Printf("Пустое описание: %v\n", err)
	}
	
	_, err = module.Info(99999)
	if err != nil {
		fmt.Printf("Несуществующий ID: %v\n", err)
	}
	
	_, err = payments.NewPaymentModule(nil)
	if err != nil {
		fmt.Printf("Nil payment method: %v\n", err)
	}
	
	invalidPayPal := payments.PayPal{}
	module2, _ := payments.NewPaymentModule(invalidPayPal)
	_, err = module2.Pay("Test", 50)
	if err != nil {
		fmt.Printf("Невалидный PayPal: %v\n", err)
	}
}

func testCancelErrors() {
	fmt.Println("--- Тест ошибок отмены ---")
	
	paypal := payments.PayPal{
		Email:  "test@example.com",
		ApiKey: "key123",
	}
	module, _ := payments.NewPaymentModule(paypal)
	
	id, _ := module.Pay("Test payment", 50)
	
	err := module.Cancel(id)
	if err != nil {
		fmt.Printf("Ошибка первой отмены: %v\n", err)
	} else {
		fmt.Printf("Первая отмена успешна\n")
	}
	
	err = module.Cancel(id)
	if err != nil {
		fmt.Printf("Повторная отмена: %v\n", err)
	}
	
	err = module.Cancel(99999)
	if err != nil {
		fmt.Printf("Отмена несуществующего: %v\n", err)
	}
}