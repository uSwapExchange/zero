package main

import (
	"encoding/json"
	"html/template"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

// Token brand colors for dynamic accent theming.
var tokenColors = map[string]string{
	"BTC":   "#F7931A",
	"ETH":   "#627EEA",
	"USDT":  "#50AF95",
	"USDC":  "#2775CA",
	"SOL":   "#9945FF",
	"BNB":   "#F3BA2F",
	"XRP":   "#23B5E8",
	"DOGE":  "#C2A633",
	"AVAX":  "#E84142",
	"DOT":   "#E6007A",
	"MATIC": "#8247E5",
	"NEAR":  "#00EC97",
	"UNI":   "#FF007A",
	"LINK":  "#2A5ADA",
	"DAI":   "#F5AC37",
	"AAVE":  "#B6509E",
	"WBTC":  "#F09242",
	"WETH":  "#627EEA",
	"ARB":   "#28A0F0",
	"OP":    "#FF0420",
	"TON":   "#0098EA",
	"LTC":   "#BFBBBB",
	"SHIB":  "#FFA409",
	"TRX":   "#FF0013",
}

func hexToRGB(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("%d, %d, %d", r, g, b)
}

func tokenColorPair(ticker string) (string, string) {
	hex := "#ffffff"
	if c, ok := tokenColors[strings.ToUpper(ticker)]; ok {
		hex = c
	}
	return hex, hexToRGB(hex)
}

// PageData is the base data passed to every template.
type PageData struct {
	Title       string
	Error       string
	MetaRefresh int // seconds; 0 = no refresh
	FromColor   string
	FromColorA  string
	ToColor     string
	ToColorA    string
	CommitHash     string
	BuildTime      string
	BuildLogURL    string
	OnionURL       string
	Description    string
	CanonicalPath  string
	StructuredData template.HTML
	NoIndex        bool
}

func newPageData(title string) PageData {
	return PageData{
		Title:       title,
		FromColor:   "#ffffff",
		FromColorA:  "255, 255, 255",
		ToColor:     "#ffffff",
		ToColorA:    "255, 255, 255",
		CommitHash:  commitHash,
		BuildTime:   buildTime,
		BuildLogURL: buildLogURL,
		OnionURL:    onionURL,
	}
}

// SwapPageData is the data for the swap form page.
type SwapPageData struct {
	PageData
	From       string
	FromNet    string
	To         string
	ToNet      string
	Amount     string
	AmountOut  string // receive amount for EXACT_OUTPUT
	Recipient  string
	RefundAddr string
	Slippage   string
	CSRFToken  string
	Networks   []NetworkGroup
	SearchFrom string
	SearchTo   string
	ModalOpen  string // "from" or "to" if a modal should be open
	FromToken  *TokenInfo
	ToToken    *TokenInfo
}

// QuotePageData is the data for the quote preview page.
type QuotePageData struct {
	PageData
	From            string
	FromNet         string
	FromTicker      string
	To              string
	ToNet           string
	ToTicker        string
	AmountIn        string
	AmountInUSD     string
	AmountOut       string
	AmountOutUSD    string
	Rate            string
	Recipient       string
	RefundAddr      string
	Slippage        string
	SlippageBPS     int
	CSRFToken       string
	OriginAsset     string
	DestAsset       string
	AtomicAmount    string
	SpreadUSD       string
	SpreadPct       string
	FromToken       *TokenInfo
	ToToken         *TokenInfo
	HasJWT          bool   // true if NEAR_INTENTS_JWT is set (0% protocol fee)
	SwapType        string // FLEX_INPUT or EXACT_OUTPUT
}

// OrderPageData is the data for the order status page.
type OrderPageData struct {
	PageData
	Token         string
	Order         *OrderData
	Status        *StatusResponse
	QRCode        string
	TimeRemaining string
	IsTerminal    bool
	StatusStep    int // 0=pending, 1=processing, 2=complete
}

// CurrenciesPageData is the data for the currencies list page.
type CurrenciesPageData struct {
	PageData
	Networks   []NetworkGroup
	TotalCount int
	Search     string
}

func renderError(w http.ResponseWriter, status int, title, message, action, actionURL string) {
	w.WriteHeader(status)
	templates.ExecuteTemplate(w, "error.html", struct {
		PageData
		Message   string
		Action    string
		ActionURL string
	}{
		PageData:  newPageData(title),
		Message:   message,
		Action:    action,
		ActionURL: actionURL,
	})
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx > 0 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

// handleSwap renders the main swap form.
func handleSwap(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderError(w, 404, "Not Found", "Page not found.", "Back to Home", "/")
		return
	}

	networks, _ := getNetworkGroups()

	data := SwapPageData{
		PageData:   newPageData("uSwap Zero"),
		From:       r.URL.Query().Get("from"),
		FromNet:    r.URL.Query().Get("from_net"),
		To:         r.URL.Query().Get("to"),
		ToNet:      r.URL.Query().Get("to_net"),
		Amount:     r.URL.Query().Get("amt"),
		AmountOut:  r.URL.Query().Get("amt_out"),
		Recipient:  r.URL.Query().Get("recipient"),
		Slippage:   r.URL.Query().Get("slippage"),
		CSRFToken:  generateCSRFToken("quote"),
		Networks:   networks,
		SearchFrom: r.URL.Query().Get("search_from"),
		SearchTo:   r.URL.Query().Get("search_to"),
		ModalOpen:  r.URL.Query().Get("modal"),
	}

	// Defaults
	if data.From == "" {
		data.From = "ETH"
		data.FromNet = "eth"
	}
	if data.To == "" {
		data.To = "USDT"
		data.ToNet = "eth"
	}
	if data.Slippage == "" {
		data.Slippage = "1"
	}

	// Set accent colors from selected currencies
	data.FromColor, data.FromColorA = tokenColorPair(data.From)
	data.ToColor, data.ToColorA = tokenColorPair(data.To)

	// Look up token info for display
	data.FromToken = findToken(data.From, data.FromNet)
	data.ToToken = findToken(data.To, data.ToNet)

	data.Description = "Swap 140+ cryptocurrencies across 29 blockchains with zero fees, zero tracking, and zero hidden markup. Open source and verifiable."
	data.CanonicalPath = "/"
	data.StructuredData = homepageSchema

	// Filter networks if search is active
	if data.SearchFrom != "" || data.SearchTo != "" {
		query := data.SearchFrom
		if query == "" {
			query = data.SearchTo
		}
		filtered := filterNetworks(networks, query)
		data.Networks = filtered
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "swap.html", data)
}

// handleQuote processes the quote form and shows a price preview.
func handleQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	// Rate limit
	ip := clientIP(r)
	if !limiter.allow(ip, 30, time.Minute) {
		renderError(w, 429, "Too Many Requests", "Please wait a moment before trying again.", "Back to Home", "/")
		return
	}

	// CSRF check
	if !verifyCSRFToken(r.FormValue("csrf"), "quote", time.Hour) {
		renderError(w, 403, "Invalid Request", "Form expired. Please go back and try again.", "Back to Home", "/")
		return
	}

	fromTicker := strings.ToUpper(r.FormValue("from"))
	fromNet := r.FormValue("from_net")
	toTicker := strings.ToUpper(r.FormValue("to"))
	toNet := r.FormValue("to_net")
	amount := r.FormValue("amount")
	amountOutForm := r.FormValue("amount_out")
	recipient := strings.TrimSpace(r.FormValue("recipient"))
	refundAddr := strings.TrimSpace(r.FormValue("refund_addr"))
	slippage := r.FormValue("slippage")

	// Validation (amount is optional — determines swap type)
	var errors []string
	if recipient == "" {
		errors = append(errors, "Recipient address is required")
	}
	if refundAddr == "" {
		errors = append(errors, "Refund address is required")
	}
	if len(errors) > 0 {
		renderError(w, 400, "Validation Error", "Please check your input:\n"+strings.Join(errors, "\n"), "Go Back", "/")
		return
	}

	// Find tokens
	fromToken := findToken(fromTicker, fromNet)
	toToken := findToken(toTicker, toNet)
	if fromToken == nil || toToken == nil {
		renderError(w, 400, "Unknown Token", "Could not find the selected tokens. Try selecting them again.", "Go Back", "/")
		return
	}

	slippageBPS, err := slippageToBPS(slippage)
	if err != nil {
		slippageBPS = 100 // default 1%
	}

	// Auto-detect swap type based on which amount field is filled.
	// Both filled: prefer send amount (FLEX_INPUT).
	swapType := "FLEX_INPUT"
	if amountOutForm != "" && amount == "" {
		swapType = "EXACT_OUTPUT"
	}

	if amount == "" && amountOutForm == "" {
		renderError(w, 400, "Amount Required", "Please enter a send or receive amount.", "Go Back", "/")
		return
	}

	// FLEX_INPUT or EXACT_OUTPUT: convert the appropriate amount to atomic.
	var atomicAmount string
	if swapType == "EXACT_OUTPUT" {
		atomicAmount, err = humanToAtomic(amountOutForm, toToken.Decimals)
	} else {
		atomicAmount, err = humanToAtomic(amount, fromToken.Decimals)
	}
	if err != nil {
		renderError(w, 400, "Invalid Amount", "Could not parse the amount: "+err.Error(), "Go Back", "/")
		return
	}

	// Request dry quote from NEAR Intents
	quoteReq := &QuoteRequest{
		Dry:                true,
		SwapType:           swapType,
		SlippageTolerance:  slippageBPS,
		OriginAsset:        fromToken.DefuseAssetID,
		DepositType:        "ORIGIN_CHAIN",
		DestinationAsset:   toToken.DefuseAssetID,
		Amount:             atomicAmount,
		RefundTo:           refundAddr,
		RefundType:         "ORIGIN_CHAIN",
		Recipient:          recipient,
		RecipientType:      "DESTINATION_CHAIN",
		Deadline:           buildDeadline(time.Hour),
		QuoteWaitingTimeMs: 8000,
		AppFees:            []struct{}{},
	}

	dryResp, err := requestDryQuote(quoteReq)
	if err != nil {
		renderError(w, 502, "Quote Failed", "NEAR Intents API is temporarily unavailable. This usually resolves in a few minutes.", "Try Again", "/")
		return
	}

	// Extract amounts from dry quote response.
	// For EXACT_OUTPUT, AmountIn is estimated and AmountOut is exact.
	// For FLEX_INPUT, both are approximate.
	humanIn := dryResp.Quote.AmountInFormatted
	humanOut := dryResp.Quote.AmountOutFormatted
	if humanIn == "" {
		humanIn = amount
	}
	if humanOut == "" {
		humanOut = atomicToHuman(dryResp.Quote.AmountOut, toToken.Decimals)
	}

	if dryResp.Quote.AmountOut == "" || dryResp.Quote.AmountOut == "0" {
		renderError(w, 502, "Quote Unavailable", "No market makers are currently offering a rate for this pair/amount. Try a larger amount or a different pair.", "Go Back", "/")
		return
	}

	// USD values
	amountInUSD := ""
	amountOutUSD := ""
	spreadUSD := ""
	spreadPct := ""
	rate := ""

	inFloat, _ := parseFloat(humanIn)
	outFloat, _ := parseFloat(humanOut)

	if fromToken.Price > 0 && inFloat > 0 {
		inUSD := inFloat * fromToken.Price
		amountInUSD = formatUSD(inUSD)

		if toToken.Price > 0 && outFloat > 0 {
			outUSD := outFloat * toToken.Price
			amountOutUSD = formatUSD(outUSD)

			spread := inUSD - outUSD
			if spread < 0 {
				spread = 0
			}
			spreadUSD = formatUSD(spread)
			if inUSD > 0 {
				spreadPct = fmt.Sprintf("%.2f%%", (spread/inUSD)*100)
			}

			rate = fmt.Sprintf("1 %s = %s %s", fromTicker, formatRate(outFloat/inFloat), toTicker)
		}
	}

	data := QuotePageData{
		PageData:     newPageData("Quote Preview"),
		From:         fromTicker,
		FromNet:      fromNet,
		FromTicker:   fromTicker,
		To:           toTicker,
		ToNet:        toNet,
		ToTicker:     toTicker,
		AmountIn:     humanIn,
		AmountInUSD:  amountInUSD,
		AmountOut:    humanOut,
		AmountOutUSD: amountOutUSD,
		Rate:         rate,
		Recipient:    recipient,
		RefundAddr:   refundAddr,
		Slippage:     slippage,
		SlippageBPS:  slippageBPS,
		CSRFToken:    generateCSRFToken("swap"),
		OriginAsset:  fromToken.DefuseAssetID,
		DestAsset:    toToken.DefuseAssetID,
		AtomicAmount: atomicAmount,
		SpreadUSD:    spreadUSD,
		SpreadPct:    spreadPct,
		FromToken:    fromToken,
		ToToken:      toToken,
		HasJWT:       nearIntentsJWT != "",
		SwapType:     swapType,
	}

	data.FromColor, data.FromColorA = tokenColorPair(fromTicker)
	data.ToColor, data.ToColorA = tokenColorPair(toTicker)

	data.NoIndex = true
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "quote.html", data)
}

// handleSwapConfirm creates a real quote and redirects to the order page.
func handleSwapConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	ip := clientIP(r)
	if !limiter.allow(ip, 10, time.Minute) {
		renderError(w, 429, "Too Many Requests", "Please wait before creating another swap.", "Back to Home", "/")
		return
	}

	if !verifyCSRFToken(r.FormValue("csrf"), "swap", time.Hour) {
		renderError(w, 403, "Invalid Request", "Form expired. Please start over.", "Back to Home", "/")
		return
	}

	fromTicker := strings.ToUpper(r.FormValue("from"))
	fromNet := r.FormValue("from_net")
	toTicker := strings.ToUpper(r.FormValue("to"))
	toNet := r.FormValue("to_net")
	atomicAmount := r.FormValue("atomic_amount")
	userAmountIn := r.FormValue("amount_in")   // user's original input
	userAmountOut := r.FormValue("amount_out")  // user's original output (EXACT_OUTPUT)
	recipient := r.FormValue("recipient")
	refundAddr := r.FormValue("refund_addr")
	slippageBPS := r.FormValue("slippage_bps")
	swapType := r.FormValue("swap_type")
	if swapType == "" {
		swapType = "FLEX_INPUT"
	}

	fromToken := findToken(fromTicker, fromNet)
	toToken := findToken(toTicker, toNet)
	if fromToken == nil || toToken == nil {
		renderError(w, 400, "Unknown Token", "Token not found.", "Back to Home", "/")
		return
	}

	bps := 100
	fmt.Sscanf(slippageBPS, "%d", &bps)

	// Real quote (not dry)
	quoteReq := &QuoteRequest{
		Dry:                false,
		SwapType:           swapType,
		SlippageTolerance:  bps,
		OriginAsset:        fromToken.DefuseAssetID,
		DepositType:        "ORIGIN_CHAIN",
		DestinationAsset:   toToken.DefuseAssetID,
		Amount:             atomicAmount,
		RefundTo:           refundAddr,
		RefundType:         "ORIGIN_CHAIN",
		Recipient:          recipient,
		RecipientType:      "DESTINATION_CHAIN",
		Deadline:           buildDeadline(time.Hour),
		QuoteWaitingTimeMs: 8000,
		AppFees:            []struct{}{},
	}

	quoteResp, err := requestQuote(quoteReq)
	if err != nil {
		renderError(w, 502, "Swap Failed", "NEAR Intents API is temporarily unavailable. This usually resolves in a few minutes.", "Try Again", "/")
		return
	}

	// For FLEX_INPUT, use the user's original amount (the API may return a
	// different amountIn since FLEX_INPUT accepts a range). For EXACT_OUTPUT,
	// use the user's desired output and the API's estimated input.
	amountIn := quoteResp.Quote.AmountInFmt
	amountOut := quoteResp.Quote.AmountOutFmt
	if swapType == "FLEX_INPUT" && userAmountIn != "" {
		amountIn = userAmountIn
	}
	if swapType == "EXACT_OUTPUT" && userAmountOut != "" {
		amountOut = userAmountOut
	}

	orderData := &OrderData{
		DepositAddr: quoteResp.Quote.DepositAddress,
		Memo:        quoteResp.Quote.DepositMemo,
		FromTicker:  fromTicker,
		FromNet:     fromNet,
		ToTicker:    toTicker,
		ToNet:       toNet,
		AmountIn:    amountIn,
		AmountOut:   amountOut,
		Deadline:    quoteResp.Quote.Deadline,
		CorrID:      quoteResp.CorrelationID,
		RefundAddr:  refundAddr,
		RecvAddr:    recipient,
		SwapType:    swapType,
	}

	token, err := encryptOrderData(orderData)
	if err != nil {
		renderError(w, 500, "Internal Error", "Failed to create order token.", "Back to Home", "/")
		return
	}

	http.Redirect(w, r, "/order/"+token, http.StatusFound)
}

// handleOrder renders the order status page.
func handleOrder(w http.ResponseWriter, r *http.Request) {
	// Extract token from path: /order/{token} or /order/{token}/raw
	path := strings.TrimPrefix(r.URL.Path, "/order/")
	isRaw := strings.HasSuffix(path, "/raw")
	if isRaw {
		path = strings.TrimSuffix(path, "/raw")
	}

	if path == "" {
		renderError(w, 400, "Missing Order", "No order token provided.", "Create New Swap", "/")
		return
	}

	order, err := decryptOrderData(path)
	if err != nil {
		renderError(w, 400, "Invalid Order", "This order link is invalid or expired. It may have been created on a different server.", "Create New Swap", "/")
		return
	}

	// Fetch live status from NEAR Intents
	status, err := fetchStatus(order.DepositAddr, order.Memo)
	if err != nil {
		// If API is down, still show what we know from the token
		status = &StatusResponse{Status: "UNKNOWN"}
	}

	if isRaw {
		w.Header().Set("Content-Type", "application/json")
		if status.RawJSON != nil {
			w.Write(status.RawJSON)
		} else {
			json.NewEncoder(w).Encode(status)
		}
		return
	}

	// Determine status step and terminal state
	isTerminal := false
	statusStep := 0
	switch status.Status {
	case "PENDING_DEPOSIT":
		statusStep = 0
	case "PROCESSING":
		statusStep = 1
	case "SUCCESS":
		statusStep = 2
		isTerminal = true
	case "REFUNDED", "FAILED", "INCOMPLETE_DEPOSIT":
		statusStep = 2
		isTerminal = true
	default:
		statusStep = 0
	}

	// Calculate time remaining
	timeRemaining := ""
	if order.Deadline != "" {
		dl, err := time.Parse(time.RFC3339, order.Deadline)
		if err == nil {
			remaining := time.Until(dl)
			if remaining > 0 {
				mins := int(remaining.Minutes())
				if mins >= 60 {
					timeRemaining = fmt.Sprintf("%dh %dm", mins/60, mins%60)
				} else {
					timeRemaining = fmt.Sprintf("%dm", mins)
				}
			} else {
				timeRemaining = "Expired"
			}
		}
	}

	// Generate QR code
	qrData := order.DepositAddr
	qrSVG := generateQRSVG(qrData, 200)

	refresh := 0
	if !isTerminal {
		refresh = 10
	}

	data := OrderPageData{
		PageData:      newPageData("Order Status"),
		Token:         path,
		Order:         order,
		Status:        status,
		QRCode:        qrSVG,
		TimeRemaining: timeRemaining,
		IsTerminal:    isTerminal,
		StatusStep:    statusStep,
	}
	data.MetaRefresh = refresh
	data.NoIndex = true
	data.FromColor, data.FromColorA = tokenColorPair(order.FromTicker)
	data.ToColor, data.ToColorA = tokenColorPair(order.ToTicker)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "order.html", data)
}

// handleCurrencies renders the full currency list.
func handleCurrencies(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	networks, err := getNetworkGroups()
	if err != nil {
		renderError(w, 502, "Unavailable", "Could not load currency list. NEAR Intents API may be temporarily unavailable.", "Try Again", "/currencies")
		return
	}

	totalCount := 0
	if search != "" {
		networks = filterNetworks(networks, search)
	}
	for _, ng := range networks {
		totalCount += len(ng.Tokens)
	}

	pd := newPageData("Supported Currencies")
	pd.Description = "Browse all supported tokens on uSwap Zero. Swap any pair across 29 blockchains with zero fees and zero tracking."
	pd.CanonicalPath = "/currencies"
	data := CurrenciesPageData{
		PageData:   pd,
		Networks:   networks,
		TotalCount: totalCount,
		Search:     search,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "currencies.html", data)
}

// handleHowItWorks renders the educational page.
func handleHowItWorks(w http.ResponseWriter, r *http.Request) {
	pd := newPageData("How It Works")
	pd.Description = "How uSwap Zero works: choose tokens, get a quote, send your deposit, receive your swap. Zero fees, zero JavaScript, fully open source."
	pd.CanonicalPath = "/how-it-works"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "how_it_works.html", pd)
}

// ResellerStats holds formatted display strings for a single reseller.
type ResellerStats struct {
	TotalSwaps   string
	TotalVolume  string
	TotalRevenue string
	FirstTx      string
	DaysActive   int
	DailyRevenue string
	UniqueSenders string
	BiggestUSD   string
}

// CombinedStats holds formatted combined stats.
type CombinedStats struct {
	TotalVolume  string
	TotalRevenue string
	TotalSwaps   string
	UniqueUsers  string
}

// CaseStudyPageData is the data for the case study page.
type CaseStudyPageData struct {
	PageData
	Eagle    ResellerStats
	Lizard   ResellerStats
	SwapMy   ResellerStats
	Combined CombinedStats
}

// caseStudyData is initialized once at startup from the embedded JSON.
var caseStudyData CaseStudyPageData

// rawAnalysis is the structure matching the JSON file.
type rawAnalysis struct {
	EagleSwap  rawReseller `json:"EagleSwap"`
	LizardSwap rawReseller `json:"LizardSwap"`
	SwapMy     rawReseller `json:"SwapMy"`
}

type rawReseller struct {
	TotalSwaps     int     `json:"total_swaps"`
	TotalVolumeUSD float64 `json:"total_volume_usd"`
	TotalRevenueUSD float64 `json:"total_revenue_usd"`
	UniqueSenders  int     `json:"unique_senders"`
	FirstTx        string  `json:"first_tx"`
	DaysActive     int     `json:"days_active"`
	DailyRevenueUSD float64 `json:"daily_revenue_usd"`
	BiggestSwapUSD float64 `json:"biggest_swap_usd"`
}

func formatResellerStats(r rawReseller) ResellerStats {
	return ResellerStats{
		TotalSwaps:    formatCommas(int64(r.TotalSwaps)),
		TotalVolume:   formatUSD(r.TotalVolumeUSD),
		TotalRevenue:  formatUSD(r.TotalRevenueUSD),
		FirstTx:       r.FirstTx,
		DaysActive:    r.DaysActive,
		DailyRevenue:  formatUSD(r.DailyRevenueUSD),
		UniqueSenders: formatCommas(int64(r.UniqueSenders)),
		BiggestUSD:    formatUSD(r.BiggestSwapUSD),
	}
}

func initCaseStudy() {
	var raw rawAnalysis
	if err := json.Unmarshal(analysisJSON, &raw); err != nil {
		log.Printf("WARNING: Failed to parse case study data: %v", err)
		return
	}

	caseStudyData.Eagle = formatResellerStats(raw.EagleSwap)
	caseStudyData.Lizard = formatResellerStats(raw.LizardSwap)
	caseStudyData.SwapMy = formatResellerStats(raw.SwapMy)
	caseStudyData.Combined = CombinedStats{
		TotalVolume:  formatUSD(raw.EagleSwap.TotalVolumeUSD + raw.LizardSwap.TotalVolumeUSD + raw.SwapMy.TotalVolumeUSD),
		TotalRevenue: formatUSD(raw.EagleSwap.TotalRevenueUSD + raw.LizardSwap.TotalRevenueUSD + raw.SwapMy.TotalRevenueUSD),
		TotalSwaps:   formatCommas(int64(raw.EagleSwap.TotalSwaps + raw.LizardSwap.TotalSwaps + raw.SwapMy.TotalSwaps)),
		UniqueUsers:  formatCommas(int64(raw.EagleSwap.UniqueSenders + raw.LizardSwap.UniqueSenders + raw.SwapMy.UniqueSenders)),
	}
}

// handleCaseStudy renders the competitor analysis page.
func handleCaseStudy(w http.ResponseWriter, r *http.Request) {
	pd := newPageData("The Crypto Swap Reseller Problem")
	pd.Description = "On-chain evidence reveals hidden fees in crypto swap resellers. Data from NEAR Intents shows how services charge secret markups."
	pd.CanonicalPath = "/case-study"
	pd.StructuredData = caseStudySchema
	data := CaseStudyPageData{
		PageData: pd,
		Eagle:    caseStudyData.Eagle,
		Lizard:   caseStudyData.Lizard,
		SwapMy:   caseStudyData.SwapMy,
		Combined: caseStudyData.Combined,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "case_study.html", data)
}

// VerifyPageData is the data for the /verify page.
type VerifyPageData struct {
	PageData
	GoVersion   string
	Uptime      string
	Requests    string
	BinarySize  string
	EnvVars     []EnvVarStatus
}

// EnvVarStatus shows whether an env var is configured.
type EnvVarStatus struct {
	Key   string
	Set   bool
}

// handleVerify renders the deployment verification page.
func handleVerify(w http.ResponseWriter, r *http.Request) {
	// Go version from build info
	goVersion := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		goVersion = info.GoVersion
	}

	// Uptime
	uptime := time.Since(serverStartTime).Round(time.Second).String()

	// Request count
	reqs := formatCommas(atomic.LoadInt64(&requestCounter))

	// Binary size
	binSize := "unknown"
	if exe, err := os.Executable(); err == nil {
		if fi, err := os.Stat(exe); err == nil {
			mb := float64(fi.Size()) / 1024 / 1024
			binSize = fmt.Sprintf("%.1f MB", mb)
		}
	}

	// Env var status (key names only — never values)
	envKeys := []string{
		"ORDER_SECRET", "NEAR_INTENTS_JWT", "NEAR_INTENTS_EXPLORER_JWT", "NEAR_INTENTS_API_URL", "PORT",
		"TG_BOT_TOKEN", "TG_APP_URL", "TG_WEBHOOK_SECRET",
		"TG_MONITOR_GROUP_ID", "TG_MAIN_CHAT_ID",
		"TG_SWAPMY_THREAD_ID", "TG_EAGLESWAP_THREAD_ID", "TG_LIZARDSWAP_THREAD_ID",
	}
	var envVars []EnvVarStatus
	for _, k := range envKeys {
		envVars = append(envVars, EnvVarStatus{Key: k, Set: os.Getenv(k) != ""})
	}

	pd := newPageData("Verify")
	pd.Description = "Verify uSwap Zero's deployment. Compare the live binary against the public source code, build logs, and commit hash."
	pd.CanonicalPath = "/verify"
	data := VerifyPageData{
		PageData:  pd,
		GoVersion: goVersion,
		Uptime:    uptime,
		Requests:  reqs,
		BinarySize: binSize,
		EnvVars:   envVars,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "verify.html", data)
}

// handleGenIcon serves dynamically generated token icon SVGs.
func handleGenIcon(w http.ResponseWriter, r *http.Request) {
	ticker := strings.TrimPrefix(r.URL.Path, "/icons/gen/")
	ticker = strings.ToUpper(ticker)
	if ticker == "" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	fmt.Fprint(w, generateTokenIconSVG(ticker))
}

// handleAPIOrderLookup decrypts an order token and returns order details as JSON.
// GET /api/order/{token}
func handleAPIOrderLookup(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/api/order/")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "No order token provided."})
		return
	}

	order, err := decryptOrderData(token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or expired order token."})
		return
	}

	// Fetch live status from NEAR Intents
	status, _ := fetchStatus(order.DepositAddr, order.Memo)
	liveStatus := "UNKNOWN"
	var swapDetails *SwapDetails
	if status != nil {
		liveStatus = status.Status
		swapDetails = status.SwapDetails
	}

	// Calculate time remaining (only meaningful for non-terminal states)
	timeRemaining := ""
	if order.Deadline != "" {
		dl, err := time.Parse(time.RFC3339, order.Deadline)
		if err == nil {
			remaining := time.Until(dl)
			if remaining > 0 {
				mins := int(remaining.Minutes())
				if mins >= 60 {
					timeRemaining = fmt.Sprintf("%dh %dm", mins/60, mins%60)
				} else {
					timeRemaining = fmt.Sprintf("%dm", mins)
				}
			} else {
				timeRemaining = "Expired"
			}
		}
	}

	// Resolve display names for networks
	fromToken := findToken(order.FromTicker, order.FromNet)
	toToken := findToken(order.ToTicker, order.ToNet)
	fromNetwork := order.FromNet
	toNetwork := order.ToNet
	if fromToken != nil {
		fromNetwork = fromToken.ChainName
	}
	if toToken != nil {
		toNetwork = toToken.ChainName
	}

	// Use actual amounts from swap details if the swap completed
	amountIn := order.AmountIn
	amountOut := order.AmountOut
	if swapDetails != nil {
		if swapDetails.AmountInFmt != "" {
			amountIn = swapDetails.AmountInFmt
		}
		if swapDetails.AmountOutFmt != "" {
			amountOut = swapDetails.AmountOutFmt
		}
	}

	result := map[string]interface{}{
		"depositAddress": order.DepositAddr,
		"memo":           order.Memo,
		"fromToken":      order.FromTicker,
		"fromNetwork":    fromNetwork,
		"toToken":        order.ToTicker,
		"toNetwork":      toNetwork,
		"amountIn":       amountIn,
		"amountOut":      amountOut,
		"deadline":       order.Deadline,
		"timeRemaining":  timeRemaining,
		"correlationId":  order.CorrID,
		"refundAddress":  order.RefundAddr,
		"receiveAddress": order.RecvAddr,
		"swapType":       order.SwapType,
		"status":         liveStatus,
	}

	// Include transaction links if available
	if swapDetails != nil {
		if len(swapDetails.DestTxs) > 0 {
			result["destinationTx"] = swapDetails.DestTxs[0].Hash
			result["destinationExplorerUrl"] = swapDetails.DestTxs[0].ExplorerURL
		}
		if swapDetails.RefundReason != "" {
			result["refundReason"] = swapDetails.RefundReason
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleAPIEstimateTime returns estimated swap times by blockchain.
// GET /api/estimate-time?from={chain}&to={chain}
func handleAPIEstimateTime(w http.ResponseWriter, r *http.Request) {
	// Typical confirmation times by blockchain (in minutes).
	// These are conservative estimates for deposit confirmation + swap execution.
	chainTimes := map[string]int{
		"btc": 30, "bitcoin": 30,
		"eth": 5, "ethereum": 5,
		"sol": 1, "solana": 1,
		"base": 3,
		"arb": 2, "arbitrum": 2,
		"op": 2, "optimism": 2,
		"pol": 3, "polygon": 3,
		"avax": 2, "avalanche": 2,
		"bsc": 2, "bnb chain": 2,
		"ton": 2,
		"tron": 3,
		"near": 1,
		"sui": 1,
		"apt": 2, "aptos": 2,
		"doge": 10, "dogecoin": 10,
		"ltc": 10, "litecoin": 10,
		"xrp": 1,
		"bch": 20, "bitcoin cash": 20,
		"xlm": 1, "stellar": 1,
	}

	fromChain := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("from")))
	toChain := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("to")))

	type chainEstimate struct {
		Chain   string `json:"chain"`
		Minutes int    `json:"minutes"`
		Label   string `json:"label"`
	}

	formatLabel := func(mins int) string {
		if mins < 2 {
			return "~1 minute"
		}
		return fmt.Sprintf("~%d minutes", mins)
	}

	// If specific chains requested, return just those
	if fromChain != "" || toChain != "" {
		fromMins := 5 // default
		toMins := 2   // default (swap execution)
		if v, ok := chainTimes[fromChain]; ok {
			fromMins = v
		}
		if v, ok := chainTimes[toChain]; ok {
			toMins = v
		}
		totalMins := fromMins + toMins

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"fromChain":     fromChain,
			"toChain":       toChain,
			"depositTime":   chainEstimate{Chain: fromChain, Minutes: fromMins, Label: formatLabel(fromMins)},
			"deliveryTime":  chainEstimate{Chain: toChain, Minutes: toMins, Label: formatLabel(toMins)},
			"totalEstimate": chainEstimate{Minutes: totalMins, Label: formatLabel(totalMins)},
		})
		return
	}

	// No params: return full table of known chains
	seen := make(map[string]bool)
	var all []chainEstimate
	for chain, mins := range chainTimes {
		if len(chain) <= 4 && !seen[chain] { // short codes only to avoid duplicates
			seen[chain] = true
			all = append(all, chainEstimate{Chain: chain, Minutes: mins, Label: formatLabel(mins)})
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Minutes < all[j].Minutes })

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"chains": all,
	})
}

// filterNetworks filters network groups by a search query.
func filterNetworks(networks []NetworkGroup, query string) []NetworkGroup {
	q := strings.ToLower(query)
	var filtered []NetworkGroup
	for _, ng := range networks {
		var tokens []TokenInfo
		for _, t := range ng.Tokens {
			if strings.Contains(strings.ToLower(t.Ticker), q) ||
				strings.Contains(strings.ToLower(t.Name), q) ||
				strings.Contains(strings.ToLower(ng.Name), q) {
				tokens = append(tokens, t)
			}
		}
		if len(tokens) > 0 {
			filtered = append(filtered, NetworkGroup{Name: ng.Name, Tokens: tokens})
		}
	}
	return filtered
}

// parseFloat is a simple float parser for display purposes only.
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// formatRate formats an exchange rate for display.
func formatRate(rate float64) string {
	if rate >= 1000 {
		return formatUSD(rate)[1:] // strip $
	}
	if rate >= 1 {
		return fmt.Sprintf("%.2f", rate)
	}
	if rate >= 0.0001 {
		return fmt.Sprintf("%.6f", rate)
	}
	// Very small rate
	return fmt.Sprintf("%.8f", math.Abs(rate))
}

// ──────────────────────────────────────────────────────────────
// SEO: Structured data schemas
// ──────────────────────────────────────────────────────────────

var homepageSchema = template.HTML(`<script type="application/ld+json">
{"@context":"https://schema.org","@graph":[{"@type":"Organization","name":"uSwap Exchange","url":"https://zero.uswap.net","logo":"https://zero.uswap.net/static/apple-touch-icon.png","sameAs":["https://github.com/uSwapExchange/zero","https://t.me/uSwapZero"]},{"@type":"WebApplication","name":"uSwap Zero","url":"https://zero.uswap.net","applicationCategory":"FinanceApplication","operatingSystem":"Any","offers":{"@type":"Offer","price":"0","priceCurrency":"USD"},"description":"Zero-fee, zero-tracking, open-source cryptocurrency swap. 140+ tokens across 29+ blockchains with no markup and no JavaScript tracking."}]}
</script>`)

var caseStudySchema = template.HTML(`<script type="application/ld+json">
{"@context":"https://schema.org","@type":"Article","headline":"The Crypto Swap Reseller Problem","description":"On-chain evidence reveals hidden fees in crypto swap resellers using NEAR Intents.","author":{"@type":"Organization","name":"uSwap Exchange"},"publisher":{"@type":"Organization","name":"uSwap Exchange","logo":{"@type":"ImageObject","url":"https://zero.uswap.net/static/apple-touch-icon.png"}},"mainEntityOfPage":"https://zero.uswap.net/case-study"}
</script>`)

var faqSchema = template.HTML(`<script type="application/ld+json">
{"@context":"https://schema.org","@type":"FAQPage","mainEntity":[{"@type":"Question","name":"What is uSwap Zero?","acceptedAnswer":{"@type":"Answer","text":"uSwap Zero is a zero-fee, zero-tracking, open-source crypto swap frontend. It lets you swap 140+ cryptocurrencies across 29+ blockchains with no markup, no hidden fees, and no JavaScript tracking."}},{"@type":"Question","name":"Is KYC required?","acceptedAnswer":{"@type":"Answer","text":"No. uSwap Zero requires no KYC, no identity verification, no account, and no sign-up. There is nothing to register, no email to provide, and no documents to upload. Just pick your tokens, enter your wallet addresses, and swap."}},{"@type":"Question","name":"Are there really no fees?","acceptedAnswer":{"@type":"Answer","text":"uSwap Zero charges zero markup. The only cost is the market maker spread, the small difference between buy and sell prices that exists on every exchange. We add nothing on top."}},{"@type":"Question","name":"Which cryptocurrencies are supported?","acceptedAnswer":{"@type":"Answer","text":"Over 140 tokens across 29+ blockchains including Bitcoin, Ethereum, Solana, NEAR, Polygon, Arbitrum, Optimism, Avalanche, BNB Chain, and more."}},{"@type":"Question","name":"How long does a swap take?","acceptedAnswer":{"@type":"Answer","text":"Most swaps complete in 1 to 30 minutes. Fast chains like Solana and NEAR take about 1 minute; Bitcoin takes around 30 minutes due to block confirmation times."}},{"@type":"Question","name":"Is uSwap Zero safe?","acceptedAnswer":{"@type":"Answer","text":"uSwap Zero is non-custodial and open source. Your funds go directly through the NEAR Intents protocol. The entire source code is public, the build is verifiable, and there is zero JavaScript tracking."}},{"@type":"Question","name":"What happens if something goes wrong?","acceptedAnswer":{"@type":"Answer","text":"If a swap cannot be completed, your funds are automatically refunded to your refund address by the NEAR Intents protocol."}},{"@type":"Question","name":"Why is there no JavaScript tracking?","acceptedAnswer":{"@type":"Answer","text":"Privacy. Most swap services embed analytics that track every click and wallet address. uSwap Zero uses zero analytics, zero cookies, and zero external requests."}},{"@type":"Question","name":"Is uSwap Zero open source?","acceptedAnswer":{"@type":"Answer","text":"Yes. The complete source code is available on GitHub under the MIT License. It is a single Go binary with zero external dependencies."}},{"@type":"Question","name":"Do you collect any personal data?","acceptedAnswer":{"@type":"Answer","text":"No. Zero personal data is collected. No analytics, no cookies, no IP logging, no wallet tracking. Swap details are encrypted in your order URL and never stored on our servers."}},{"@type":"Question","name":"What is the market maker spread?","acceptedAnswer":{"@type":"Answer","text":"The spread is the difference between what a market maker pays for a token and what they sell it for. Unlike other swap services, uSwap Zero adds zero additional markup on top of this spread."}}]}
</script>`)

var feesSchema = template.HTML(`<script type="application/ld+json">
{"@context":"https://schema.org","@type":"FAQPage","mainEntity":[{"@type":"Question","name":"Does uSwap Zero charge any fees?","acceptedAnswer":{"@type":"Answer","text":"No. uSwap Zero charges zero fees, zero markup, and zero commissions. The only cost is the market maker spread, the small difference between buy and sell prices that exists on every exchange."}},{"@type":"Question","name":"What is the market maker spread?","acceptedAnswer":{"@type":"Answer","text":"The spread is the natural difference between a market maker's buy and sell price. It typically ranges from 0.1% to 0.5% depending on the token pair and liquidity. This cost exists on every exchange. uSwap Zero adds nothing on top of it."}},{"@type":"Question","name":"How is uSwap Zero different from other swap services?","acceptedAnswer":{"@type":"Answer","text":"Most crypto swap services add a hidden markup of 0.25% to 4% on top of the market maker spread. uSwap Zero adds nothing. Our case study proves this with on-chain transaction data."}}]}
</script>`)

// ──────────────────────────────────────────────────────────────
// SEO: New content pages
// ──────────────────────────────────────────────────────────────

// handlePrivacy renders the privacy policy page.
func handlePrivacy(w http.ResponseWriter, r *http.Request) {
	pd := newPageData("Privacy Policy")
	pd.Description = "uSwap Zero collects zero personal data. No KYC, no accounts, no analytics, no cookies, no IP logging. Fully non-custodial and open source."
	pd.CanonicalPath = "/privacy"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "privacy.html", pd)
}

// handleTerms renders the terms of service page.
func handleTerms(w http.ResponseWriter, r *http.Request) {
	pd := newPageData("Terms of Service")
	pd.Description = "uSwap Zero terms of service. No KYC, no sign-up, non-custodial, zero fees, open source. Plain language, no legalese."
	pd.CanonicalPath = "/terms"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "terms.html", pd)
}

// handleFees renders the fee transparency page.
func handleFees(w http.ResponseWriter, r *http.Request) {
	pd := newPageData("Fees")
	pd.Description = "uSwap Zero charges zero fees. No markup, no hidden charges — just the market maker spread. See exactly what you pay."
	pd.CanonicalPath = "/fees"
	pd.StructuredData = feesSchema
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "fees.html", pd)
}

// handleFAQ renders the frequently asked questions page.
func handleFAQ(w http.ResponseWriter, r *http.Request) {
	pd := newPageData("FAQ")
	pd.Description = "Common questions about uSwap Zero: how swaps work, security, supported tokens, fees, and privacy."
	pd.CanonicalPath = "/faq"
	pd.StructuredData = faqSchema
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "faq.html", pd)
}

// ──────────────────────────────────────────────────────────────
// SEO: Machine-readable meta pages
// ──────────────────────────────────────────────────────────────

// handleRobotsTxt serves crawl directives.
func handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("User-agent: *\n" +
		"Allow: /\n" +
		"Disallow: /api/\n" +
		"Disallow: /tg/\n" +
		"Disallow: /wrapper-logs\n" +
		"Disallow: /order/\n" +
		"Disallow: /quote\n" +
		"Disallow: /swap\n" +
		"\n" +
		"Sitemap: https://zero.uswap.net/sitemap.xml\n"))
}

// handleSitemapXML serves the XML sitemap.
func handleSitemapXML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>https://zero.uswap.net/</loc><changefreq>daily</changefreq><priority>1.0</priority></url>
  <url><loc>https://zero.uswap.net/currencies</loc><changefreq>weekly</changefreq><priority>0.8</priority></url>
  <url><loc>https://zero.uswap.net/how-it-works</loc><changefreq>monthly</changefreq><priority>0.8</priority></url>
  <url><loc>https://zero.uswap.net/fees</loc><changefreq>monthly</changefreq><priority>0.8</priority></url>
  <url><loc>https://zero.uswap.net/case-study</loc><changefreq>monthly</changefreq><priority>0.7</priority></url>
  <url><loc>https://zero.uswap.net/faq</loc><changefreq>monthly</changefreq><priority>0.6</priority></url>
  <url><loc>https://zero.uswap.net/verify</loc><changefreq>weekly</changefreq><priority>0.5</priority></url>
  <url><loc>https://zero.uswap.net/privacy</loc><changefreq>monthly</changefreq><priority>0.4</priority></url>
  <url><loc>https://zero.uswap.net/terms</loc><changefreq>monthly</changefreq><priority>0.4</priority></url>
</urlset>`))
}

// handleLLMSTxt serves AI crawler guidance.
func handleLLMSTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("# uSwap Zero\n" +
		"\n" +
		"> Zero-fee, zero-tracking, open-source crypto swap frontend.\n" +
		"\n" +
		"uSwap Zero lets you swap 140+ cryptocurrencies across 29+ blockchains\n" +
		"with no fees, no tracking, and no hidden markup. Powered by NEAR Intents protocol.\n" +
		"\n" +
		"## Key Facts\n" +
		"\n" +
		"- Zero fees: no markup, no commissions, no hidden charges\n" +
		"- Zero tracking: no analytics, no cookies, no JavaScript tracking\n" +
		"- Open source: MIT License, single Go binary, zero external dependencies\n" +
		"- Non-custodial: funds go directly through NEAR Intents protocol\n" +
		"- Verifiable: public source code, reproducible builds\n" +
		"\n" +
		"## Pages\n" +
		"\n" +
		"- [Home](https://zero.uswap.net/): Swap form for 140+ tokens\n" +
		"- [Currencies](https://zero.uswap.net/currencies): Full list of supported tokens and networks\n" +
		"- [How It Works](https://zero.uswap.net/how-it-works): Step-by-step swap process\n" +
		"- [Fees](https://zero.uswap.net/fees): Fee transparency — we charge nothing\n" +
		"- [Case Study](https://zero.uswap.net/case-study): On-chain evidence of hidden fees in competitor swap services\n" +
		"- [FAQ](https://zero.uswap.net/faq): Common questions about swaps, fees, privacy, and security\n" +
		"- [Verify](https://zero.uswap.net/verify): Deployment verification and build metadata\n" +
		"- [Privacy Policy](https://zero.uswap.net/privacy): No KYC, no data collection, no tracking\n" +
		"- [Terms of Service](https://zero.uswap.net/terms): Plain language terms, non-custodial, no sign-up\n" +
		"- [Source Code](https://github.com/uSwapExchange/zero): Complete source on GitHub\n"))
}
