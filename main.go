package main

import (
	"fmt"
	"frenly/bot"
)

func main() {
	bot.Start()

	amountOut, price, err := bot.CalculateTradePrice(200, 20, 60)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Input: 5 FREN\n")
	fmt.Printf("Output: %.4f TON\n", amountOut)
	fmt.Printf("Effective Price: %.4f TON per FREN\n", price)
}
