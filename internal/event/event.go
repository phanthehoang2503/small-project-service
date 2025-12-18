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
	// order domain
	RoutingKeyOrderCreated   = "order.created"
	RoutingKeyOrderCancelled = "order.cancelled"

	// payment domain
	RoutingKeyPaymentSucceeded = "payment.succeeded"
	RoutingKeyPaymentFailed    = "payment.failed"

	// inventory domain
	RoutingKeyInventoryReserved          = "inventory.reserved"
	RoutingKeyInventoryReservationFailed = "inventory.reservation.failed"
)
