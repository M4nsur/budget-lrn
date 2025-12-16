package payments

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAmount       = errors.New("invalid amount: must be positive")
	ErrAmountTooLarge      = errors.New("amount exceeds maximum limit of 100000")
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrPaymentFailed       = errors.New("payment processing failed")
	ErrAlreadyCancelled    = errors.New("payment already cancelled")
	ErrCannotCancelPayment = errors.New("cannot cancel payment")
)

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