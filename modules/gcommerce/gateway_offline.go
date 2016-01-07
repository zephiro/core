package gcommerce

type GatewayOffline struct {
	di    *Module
	order *Order
	meta  map[string]interface{}
}

// Set DI instance
func (this *GatewayOffline) SetDI(di *Module) {
	this.di = di
}

func (this *GatewayOffline) SetOrder(order *Order) {
	this.order = order
}

func (this *GatewayOffline) SetMeta(meta map[string]interface{}) {
	this.meta = meta
}

func (this *GatewayOffline) Charge(amount float64) error {
	return nil
}

func (this *GatewayOffline) ModifyPrice(p float64) float64 {
	return p
}

func (this *GatewayOffline) AdjustPrice(p float64) float64 {
	return p
}
