package transaction

import (
	transaction_query "nearbyassist/internal/db/query/transaction"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func OngoingVendorTransactions(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing user ID",
		})
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user ID must be a number",
		})
	}

    transactions, err := transaction_query.VendorOngoingTransactions(id)

	return c.JSON(http.StatusOK, transactions)
}
