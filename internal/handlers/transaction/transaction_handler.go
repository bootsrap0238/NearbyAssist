package transaction

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/transaction/v1"

	"github.com/labstack/echo/v4"
)

func TransactionsHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "health check for transactions route"

	r.GET("/ongoing/client/:userId", transaction.OngoingClientTransactions).Name = "get client ongoing transactions"
	r.GET("/ongoing/vendor/:userId", transaction.OngoingVendorTransactions).Name = "get vendor ongoing transactions"
	r.GET("/history/client/:userId", transaction.GetClientTransactionHistory).Name = "get client transaction history"
	r.GET("/history/vendor/:userId", transaction.GetVendorTransactionHistory).Name = "get vendor transaction history"
}