package main

import (
	"errors"
	"fmt"
	"log"
)

type MockPaymentMethod struct {
	nextID      int
	failPayment bool
	failCancel  bool
}

func NewMockPaymentMethod() *MockPaymentMethod {
	return &MockPaymentMethod{
		nextID: 1,
	}
}

func (m *MockPaymentMethod) pay(usd int) (int, error) {
	if m.failPayment {
		return 0, errors.New("payment gateway error")
	}
	id := m.nextID
	m.nextID++
	return id, nil
}

func (m *MockPaymentMethod) cancel(id int) error {
	if m.failCancel {
		return errors.New("cancellation not allowed by gateway")
	}
	return nil
}

type PaymentInfo struct {
	Description string
	Usd         int
	Cancelled   bool
}

type PaymentMethod interface {
	pay(usd int) (int, error)
	cancel(id int) error
}

type PaymentModule struct {
	paymentsInfo  map[int]PaymentInfo
	paymentMethod PaymentMethod
}

var (
	ErrInvalidAmount       = errors.New("invalid amount: must be positive")
	ErrAmountTooLarge      = errors.New("amount exceeds maximum limit of 100000")
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrPaymentFailed       = errors.New("payment processing failed")
	ErrAlreadyCancelled    = errors.New("payment already cancelled")
	ErrCannotCancelPayment = errors.New("cannot cancel payment")
)

func NewPaymentModule(paymentMethod PaymentMethod) (*PaymentModule, error) {
	if paymentMethod == nil {
		return nil, errors.New("payment method cannot be nil")
	}
	return &PaymentModule{
		paymentsInfo:  make(map[int]PaymentInfo),
		paymentMethod: paymentMethod,
	}, nil
}

func (p *PaymentModule) Pay(descr string, usd int) (int, error) {
	if usd <= 0 {
		return 0, ErrInvalidAmount
	}
	if usd > 100000 {
		return 0, ErrAmountTooLarge
	}
	if descr == "" {
		return 0, errors.New("description cannot be empty")
	}
	
	id, err := p.paymentMethod.pay(usd)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrPaymentFailed, err)
	}
	if id <= 0 {
		return 0, errors.New("invalid payment ID returned")
	}
	
	info := PaymentInfo{
		Description: descr,
		Usd:         usd,
		Cancelled:   false,
	}
	p.paymentsInfo[id] = info
	return id, nil
}

func (p *PaymentModule) Cancel(id int) error {
	info, ok := p.paymentsInfo[id]
	if !ok {
		return fmt.Errorf("%w: ID %d", ErrPaymentNotFound, id)
	}
	if info.Cancelled {
		return fmt.Errorf("%w: ID %d", ErrAlreadyCancelled, id)
	}
	
	err := p.paymentMethod.cancel(id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCannotCancelPayment, err)
	}
	
	info.Cancelled = true
	p.paymentsInfo[id] = info
	return nil
}

func (p *PaymentModule) Info(id int) (PaymentInfo, error) {
	info, ok := p.paymentsInfo[id]
	if !ok {
		return PaymentInfo{}, fmt.Errorf("%w: ID %d", ErrPaymentNotFound, id)
	}
	return info, nil
}

func (p *PaymentModule) InfoAll() map[int]PaymentInfo {
	tempPaymentsInfo := make(map[int]PaymentInfo, len(p.paymentsInfo))
	for id, info := range p.paymentsInfo {
		tempPaymentsInfo[id] = info
	}
	return tempPaymentsInfo
}

func testSuccessfulPayment(pm *PaymentModule) {
	fmt.Println("\nTest: Successful Payment")
	id, err := pm.Pay("Coffee", 5)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Payment successful. ID: %d\n", id)
	
	info, _ := pm.Info(id)
	fmt.Printf("Description: %s, Amount: $%d, Cancelled: %v\n", 
		info.Description, info.Usd, info.Cancelled)
}

func testMultiplePayments(pm *PaymentModule) {
	fmt.Println("\nTest: Multiple Payments")
	
	payments := []struct {
		descr string
		usd   int
	}{
		{"Book", 25},
		{"Subscription", 99},
		{"Donation", 50},
	}
	
	for _, p := range payments {
		id, err := pm.Pay(p.descr, p.usd)
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Payment ID %d: %s - $%d\n", id, p.descr, p.usd)
	}
}

func testInvalidAmount(pm *PaymentModule) {
	fmt.Println("\nTest: Invalid Amount")
	
	testCases := []struct {
		descr  string
		amount int
		name   string
	}{
		{"Test", 0, "Zero amount"},
		{"Test", -100, "Negative amount"},
		{"Test", 150000, "Amount too large"},
	}
	
	for _, tc := range testCases {
		_, err := pm.Pay(tc.descr, tc.amount)
		if err != nil {
			fmt.Printf("%s rejected: %v\n", tc.name, err)
		}
	}
}

func testEmptyDescription(pm *PaymentModule) {
	fmt.Println("\nTest: Empty Description")
	_, err := pm.Pay("", 100)
	if err != nil {
		fmt.Printf("Empty description rejected: %v\n", err)
	}
}

func testCancelPayment(pm *PaymentModule) {
	fmt.Println("\nTest: Cancel Payment")
	
	id, err := pm.Pay("Refundable Item", 75)
	if err != nil {
		log.Printf("Error creating payment: %v\n", err)
		return
	}
	fmt.Printf("Payment created. ID: %d\n", id)
	
	err = pm.Cancel(id)
	if err != nil {
		log.Printf("Error cancelling: %v\n", err)
		return
	}
	fmt.Printf("Payment %d cancelled\n", id)
	
	info, _ := pm.Info(id)
	fmt.Printf("Cancelled status: %v\n", info.Cancelled)
}

func testDoubleCancellation(pm *PaymentModule) {
	fmt.Println("\nTest: Double Cancellation")
	
	id, _ := pm.Pay("Test Item", 50)
	pm.Cancel(id)
	
	err := pm.Cancel(id)
	if err != nil {
		fmt.Printf("Double cancellation rejected: %v\n", err)
	}
}

func testCancelNonexistentPayment(pm *PaymentModule) {
	fmt.Println("\nTest: Cancel Nonexistent Payment")
	
	err := pm.Cancel(99999)
	if err != nil {
		fmt.Printf("Cancelling nonexistent payment rejected: %v\n", err)
	}
}

func testInfoAll(pm *PaymentModule) {
	fmt.Println("\nTest: Info All Payments")
	
	pm.Pay("Item A", 10)
	pm.Pay("Item B", 20)
	pm.Pay("Item C", 30)
	
	allPayments := pm.InfoAll()
	fmt.Printf("Total payments: %d\n", len(allPayments))
	
	for id, info := range allPayments {
		fmt.Printf("ID %d: %s - $%d (Cancelled: %v)\n", 
			id, info.Description, info.Usd, info.Cancelled)
	}
}

func testPaymentMethodFailure(mock *MockPaymentMethod) {
	fmt.Println("\nTest: Payment Method Failure")
	
	mock.failPayment = true
	pm, _ := NewPaymentModule(mock)
	
	_, err := pm.Pay("Test", 100)
	if err != nil {
		fmt.Printf("Payment failure handled: %v\n", err)
	}
	
	mock.failPayment = false
}

func testCancellationFailure(pm *PaymentModule, mock *MockPaymentMethod) {
	fmt.Println("\nTest: Cancellation Failure")
	
	id, _ := pm.Pay("Test Item", 50)
	
	mock.failCancel = true
	err := pm.Cancel(id)
	if err != nil {
		fmt.Printf("Cancellation failure handled: %v\n", err)
	}
	
	mock.failCancel = false
}

func main() {
	fmt.Println("Payment Module Tests")
	
	mock := NewMockPaymentMethod()
	pm, err := NewPaymentModule(mock)
	if err != nil {
		log.Fatal(err)
	}
	
	testSuccessfulPayment(pm)
	testMultiplePayments(pm)
	testInvalidAmount(pm)
	testEmptyDescription(pm)
	testCancelPayment(pm)
	testDoubleCancellation(pm)
	testCancelNonexistentPayment(pm)
	testInfoAll(pm)
	testPaymentMethodFailure(mock)
	testCancellationFailure(pm, mock)
	
	fmt.Println("\nTests completed")
}