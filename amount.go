package main

import (
	"fmt"
	"math/big"
	"strings"
)

// humanToAtomic converts a human-readable amount like "0.5" to atomic units
// given the token's decimal places. Uses big.Int exclusively — never floats.
// Example: humanToAtomic("0.5", 18) → "500000000000000000"
func humanToAtomic(amount string, decimals int) (string, error) {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return "", fmt.Errorf("empty amount")
	}

	// Split on decimal point
	parts := strings.SplitN(amount, ".", 2)
	if len(parts) > 2 {
		return "", fmt.Errorf("invalid amount: multiple decimal points")
	}

	whole := parts[0]
	if whole == "" {
		whole = "0"
	}

	frac := ""
	if len(parts) == 2 {
		frac = parts[1]
	}

	// Trim or pad fractional part to exactly `decimals` digits
	if len(frac) > decimals {
		// Truncate excess precision
		frac = frac[:decimals]
	} else {
		// Pad with zeros
		frac = frac + strings.Repeat("0", decimals-len(frac))
	}

	// Combine: "0" + "500000000000000000" → "0500000000000000000"
	combined := whole + frac

	// Parse as big.Int
	result := new(big.Int)
	_, ok := result.SetString(combined, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount: %q", amount)
	}

	// big.Int handles leading zeros correctly
	return result.String(), nil
}

// atomicToHuman converts atomic units back to a human-readable string.
// Example: atomicToHuman("500000000000000000", 18) → "0.5"
func atomicToHuman(atomic string, decimals int) string {
	val := new(big.Int)
	_, ok := val.SetString(atomic, 10)
	if !ok {
		return "0"
	}

	if decimals == 0 {
		return val.String()
	}

	// divisor = 10^decimals
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	whole := new(big.Int).Div(val, divisor)
	remainder := new(big.Int).Mod(val, divisor)

	if remainder.Sign() == 0 {
		return whole.String()
	}

	// Format remainder with leading zeros
	fracStr := fmt.Sprintf("%0*s", decimals, remainder.String())
	// Trim trailing zeros
	fracStr = strings.TrimRight(fracStr, "0")

	return whole.String() + "." + fracStr
}

// formatUSD formats a price string for display, e.g. "927.45" → "$927.45"
func formatUSD(amount float64) string {
	if amount >= 1000 {
		// Use Sprintf for correct rounding, then split to add commas
		s := fmt.Sprintf("%.2f", amount) // e.g. "1234.56"
		dot := strings.Index(s, ".")
		intPart := s[:dot]
		fracPart := s[dot+1:]
		var n int64
		fmt.Sscanf(intPart, "%d", &n)
		if fracPart == "00" {
			return "$" + formatCommas(n)
		}
		return "$" + formatCommas(n) + "." + fracPart
	}
	if amount >= 1 {
		return fmt.Sprintf("$%.2f", amount)
	}
	if amount >= 0.01 {
		return fmt.Sprintf("$%.4f", amount)
	}
	return fmt.Sprintf("$%.6f", amount)
}

// formatCommas adds thousand separators to an integer.
func formatCommas(n int64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

// slippageTooltipBPS converts a percentage string like "1" to basis points integer.
// "1" → 100, "0.5" → 50, "2" → 200
func slippageToBPS(pct string) (int, error) {
	// Parse as float then multiply by 100 — acceptable here since
	// slippage is a small user-facing number, not a financial amount.
	pct = strings.TrimSpace(pct)
	if pct == "" {
		return 100, nil // default 1%
	}

	val := new(big.Float)
	_, ok := val.SetString(pct)
	if !ok {
		return 0, fmt.Errorf("invalid slippage: %q", pct)
	}

	bps := new(big.Float).Mul(val, big.NewFloat(100))
	result, _ := bps.Int64()
	if result < 1 || result > 5000 {
		return 0, fmt.Errorf("slippage out of range: %d bps", result)
	}
	return int(result), nil
}
