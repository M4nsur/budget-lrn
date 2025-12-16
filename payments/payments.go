package payments

type PaymentMethod interface {
	pay(id int) int
	cancel(id int)
}

type PaymentModule struct {
	paymentsInfo map[int]PaymentInfo
	paymentMethod PaymentMethod
}

func NewPaymentModule(PaymentMethod PaymentMethod) *PaymentModule {
	return &PaymentModule{
		paymentsInfo: make(map[int]PaymentInfo),
		paymentMethod: PaymentMethod,
	}
}

func (p *PaymentModule) Pay(descr string, usd int) int{
	id := p.paymentMethod.pay(usd)
	info := PaymentInfo{
		Description: descr,
		Usd: usd,
		Cancelled: false,
	}
	p.paymentsInfo[id] = info
	return  id
}


func (p *PaymentModule) Cancel(id int) {
	info, ok := p.paymentsInfo[id]
	if !ok {
		return
	}
	
	p.paymentMethod.cancel(id)
	info.Cancelled = true
	p.paymentsInfo[id] = info
} 
	

func (p PaymentModule) Info(id int) PaymentInfo {
	info, ok := p.paymentsInfo[id]

	if !ok {
		return PaymentInfo{}
	}
	return info
}

func (p PaymentModule) InfoAll() map[int]PaymentInfo {
	tempPaymentsInfo := make(map[int]PaymentInfo, len(p.paymentsInfo)) 
	for i, v := range p.paymentsInfo {
		tempPaymentsInfo[i] = v
	}

	return tempPaymentsInfo
}

