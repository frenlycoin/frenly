package bot

import (
	"errors"
	"math"
	"time"
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

func exchange(u *User) (amountOut float64, err error) {
	if u == nil {
		return 0, errors.New("user is nil")
	}

	amountIn := float64(u.rewards(true)) / float64(Mul9)
	if amountIn <= 0 {
		return 0, errors.New("no rewards available for exchange")
	}

	reserveAKv := &KeyValue{Key: "dexFren"}
	reserveBKv := &KeyValue{Key: "dexGram"}

	if err = db.Where("key = ?", reserveAKv.Key).FirstOrCreate(reserveAKv).Error; err != nil {
		return 0, err
	}
	if err = db.Where("key = ?", reserveBKv.Key).FirstOrCreate(reserveBKv).Error; err != nil {
		return 0, err
	}

	reserveA := float64(reserveAKv.ValueInt) / float64(Mul9)
	reserveB := float64(reserveBKv.ValueInt) / float64(Mul9)

	amountOut, effectivePrice, err := CalculateTradePrice(reserveA, reserveB, amountIn)
	if err != nil {
		return 0, err
	}

	newReserveA := int64(math.Round((reserveA + amountIn) * float64(Mul9)))
	newReserveB := int64(math.Round((reserveB - amountOut) * float64(Mul9)))
	lastPrice := int64(math.Round(effectivePrice * float64(Mul9)))

	reserveAKv.ValueInt = newReserveA
	reserveBKv.ValueInt = newReserveB

	priceKv := &KeyValue{Key: "dexLastPrice"}
	if err = db.Where("key = ?", priceKv.Key).FirstOrCreate(priceKv).Error; err != nil {
		return 0, err
	}
	priceKv.ValueInt = lastPrice

	if err = db.Save(reserveAKv).Error; err != nil {
		return 0, err
	}
	if err = db.Save(reserveBKv).Error; err != nil {
		return 0, err
	}
	if err = db.Save(priceKv).Error; err != nil {
		return 0, err
	}

	u.MiningTime = time.Now()
	u.LastUpdated = time.Now()
	u.CycleCountTotal += u.CycleCount
	u.CycleCount = 1
	if err = db.Save(u).Error; err != nil {
		return 0, err
	}

	tonAmount := int64(math.Round(amountOut * float64(Mul9)))
	// frenAmount := int64(math.Round(amountIn * float64(Mul9)))
	// msg := fmt.Sprintf(lExchangeCompleted, u.Name, float64(frenAmount)/float64(Mul9), float64(tonAmount)/float64(Mul9))

	// notify(msg, GroupHall)
	notifyCashout(u.Name, float64(tonAmount)/float64(Mul9), GroupHall)

	return amountOut, nil
}
