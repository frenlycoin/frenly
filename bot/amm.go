package bot

import (
	"errors"
)

// CalculateTradePrice calculates the output amount and effective price for a trade
// using Uniswap v1 constant product formula (x * y = k)
//
// Parameters:
//
//	reserveA     - Current reserve of Token A
//	reserveB     - Current reserve of Token B
//	amountInA    - Amount of Token A being sold (input)
//
// Returns:
//
//	amountOutB   - Amount of Token B the user will receive
//	effectivePrice - Effective price (B per A)
//	err          - Error if any
func CalculateTradePrice(reserveA, reserveB, amountInA float64) (amountOutB, effectivePrice float64, err error) {
	if reserveA <= 0 || reserveB <= 0 {
		return 0, 0, errors.New("reserves must be greater than zero")
	}
	if amountInA <= 0 {
		return 0, 0, errors.New("input amount must be greater than zero")
	}

	// Constant product formula: k = x * y
	k := reserveA * reserveB

	// New reserve A after input
	newReserveA := reserveA + amountInA

	// Calculate output amount B
	amountOutB = reserveB - (k / newReserveA)

	// Effective price (how many B you get per A)
	effectivePrice = amountOutB / amountInA

	// Basic validation
	if amountOutB <= 0 {
		return 0, 0, errors.New("insufficient liquidity or trade too large")
	}

	return amountOutB, effectivePrice, nil
}
