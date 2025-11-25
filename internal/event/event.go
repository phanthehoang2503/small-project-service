package event

// Exchanges
const (
	ExchangeLogs    = "logs_exchange"
	ExchangeProduct = "product_exchange"
	ExchangeOrder   = "order_exchange"
)

// Routing keys
const (
	// logs
	RoutingKeyLogInfo  = "log.info"
	RoutingKeyLogError = "log.error"
	RoutingKeyLogWarn  = "log.warn"

	// product domain
	RoutingKeyProductCreated = "product.created"
	RoutingKeyProductUpdated = "product.updated"
	RoutingKeyProductDeleted = "product.deleted"

	// order domain
	RoutingKeyOrderRequested = "order.requested"
	RoutingKeyOrderCreated   = "order.created"
	RoutingKeyOrderPaid      = "order.paid"
	RoutingKeyOrderCancelled = "order.cancelled"

	// payment domain
	RoutingKeyPaymentFailed = "payment.failed"
)
