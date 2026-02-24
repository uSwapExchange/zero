package main

import (
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

const tokenCacheTTL = 5 * time.Minute

// NetworkGroup holds tokens grouped by their chain/network.
type NetworkGroup struct {
	Name   string
	Tokens []TokenInfo
}

type tokenCache struct {
	mu        sync.RWMutex
	tokens    []TokenInfo
	byAssetID map[string]*TokenInfo
	networks  []NetworkGroup
	updatedAt time.Time
}

var cache = &tokenCache{}

// refreshTokenCache fetches and caches the token list from NEAR Intents.
func refreshTokenCache() error {
	tokens, err := fetchTokens()
	if err != nil {
		return err
	}

	byAssetID := make(map[string]*TokenInfo, len(tokens))
	networkMap := make(map[string][]TokenInfo)

	// Map API blockchain codes to display names
	chainDisplayName := map[string]string{
		"eth": "Ethereum", "btc": "Bitcoin", "sol": "Solana", "base": "Base",
		"arb": "Arbitrum", "ton": "TON", "tron": "TRON", "bsc": "BNB Chain",
		"pol": "Polygon", "op": "Optimism", "avax": "Avalanche", "near": "NEAR",
		"sui": "Sui", "apt": "Aptos", "aptos": "Aptos", "doge": "Dogecoin", "ltc": "Litecoin",
		"xrp": "XRP", "bch": "Bitcoin Cash", "xlm": "Stellar", "stellar": "Stellar", "zec": "Zcash",
		"cardano": "Cardano", "starknet": "StarkNet", "gnosis": "Gnosis",
		"bera": "Berachain", "monad": "Monad", "plasma": "Plasma",
		"xlayer": "X Layer", "aleo": "Aleo", "adi": "ADI",
	}

	for i := range tokens {
		t := &tokens[i]
		// Normalize ticker
		if t.Ticker == "" && t.Symbol != "" {
			t.Ticker = t.Symbol
		}
		t.Ticker = strings.ToUpper(t.Ticker)

		byAssetID[t.DefuseAssetID] = t

		// Map blockchain code to display name
		netName := t.ChainName
		if displayName, ok := chainDisplayName[strings.ToLower(netName)]; ok {
			netName = displayName
		}
		if netName == "" {
			netName = "Other"
		}
		networkMap[netName] = append(networkMap[netName], *t)
	}

	// Sort networks: popular first, then alphabetical
	networkOrder := map[string]int{
		"Ethereum": 1, "Bitcoin": 2, "Solana": 3, "Base": 4,
		"Arbitrum": 5, "TON": 6, "TRON": 7, "BNB Chain": 8,
		"Polygon": 9, "Optimism": 10, "Avalanche": 11, "NEAR": 12,
	}

	var networks []NetworkGroup
	for name, toks := range networkMap {
		// Sort tokens within each network by price (descending)
		sort.Slice(toks, func(i, j int) bool {
			return toks[i].Price > toks[j].Price
		})
		networks = append(networks, NetworkGroup{Name: name, Tokens: toks})
	}
	sort.Slice(networks, func(i, j int) bool {
		oi, oki := networkOrder[networks[i].Name]
		oj, okj := networkOrder[networks[j].Name]
		if oki && okj {
			return oi < oj
		}
		if oki {
			return true
		}
		if okj {
			return false
		}
		return networks[i].Name < networks[j].Name
	})

	cache.mu.Lock()
	cache.tokens = tokens
	cache.byAssetID = byAssetID
	cache.networks = networks
	cache.updatedAt = time.Now()
	cache.mu.Unlock()

	log.Printf("Token cache refreshed: %d tokens across %d networks", len(tokens), len(networks))
	return nil
}

// getTokens returns the cached token list, refreshing if stale.
func getTokens() ([]TokenInfo, error) {
	cache.mu.RLock()
	if time.Since(cache.updatedAt) < tokenCacheTTL && len(cache.tokens) > 0 {
		tokens := cache.tokens
		cache.mu.RUnlock()
		return tokens, nil
	}
	cache.mu.RUnlock()

	if err := refreshTokenCache(); err != nil {
		// Return stale data if available
		cache.mu.RLock()
		defer cache.mu.RUnlock()
		if len(cache.tokens) > 0 {
			log.Printf("Token refresh failed, using stale cache: %v", err)
			return cache.tokens, nil
		}
		return nil, err
	}

	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.tokens, nil
}

// getNetworkGroups returns tokens grouped by network, refreshing if stale.
func getNetworkGroups() ([]NetworkGroup, error) {
	_, err := getTokens() // triggers refresh if needed
	if err != nil {
		return nil, err
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.networks, nil
}

// findToken looks up a token by ticker and network.
// Prioritizes exact blockchain match over asset ID substring match.
func findToken(ticker, network string) *TokenInfo {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	ticker = strings.ToUpper(ticker)
	network = strings.ToLower(network)

	if network == "" {
		// No network filter — return first match by ticker
		for i := range cache.tokens {
			if strings.EqualFold(cache.tokens[i].Ticker, ticker) {
				return &cache.tokens[i]
			}
		}
		return nil
	}

	// First pass: exact blockchain match (preferred)
	for i := range cache.tokens {
		t := &cache.tokens[i]
		if strings.EqualFold(t.Ticker, ticker) && strings.EqualFold(t.ChainName, network) {
			return t
		}
	}

	// Second pass: asset ID substring match (fallback)
	for i := range cache.tokens {
		t := &cache.tokens[i]
		if strings.EqualFold(t.Ticker, ticker) &&
			strings.Contains(strings.ToLower(t.DefuseAssetID), network) {
			return t
		}
	}

	return nil
}

// findTokenByAssetID looks up a token by its NEAR Intents asset ID.
func findTokenByAssetID(assetID string) *TokenInfo {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	if t, ok := cache.byAssetID[assetID]; ok {
		return t
	}
	return nil
}

// searchTokens returns tokens matching a query string (ticker, name, or network).
func searchTokens(query string) []TokenInfo {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	if query == "" {
		return cache.tokens
	}

	q := strings.ToLower(query)
	var results []TokenInfo
	for _, t := range cache.tokens {
		if strings.Contains(strings.ToLower(t.Ticker), q) ||
			strings.Contains(strings.ToLower(t.Name), q) ||
			strings.Contains(strings.ToLower(t.ChainName), q) {
			results = append(results, t)
		}
	}
	return results
}

// startCacheRefresher starts a background goroutine to keep the cache fresh.
func startCacheRefresher() {
	// Initial load
	if err := refreshTokenCache(); err != nil {
		log.Printf("Initial token cache load failed (will retry): %v", err)
	}

	go func() {
		ticker := time.NewTicker(tokenCacheTTL)
		defer ticker.Stop()
		for range ticker.C {
			if err := refreshTokenCache(); err != nil {
				log.Printf("Token cache refresh failed: %v", err)
			}
		}
	}()
}
