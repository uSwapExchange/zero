# Site Structure — uSwap Zero

## Proposed URL Hierarchy

```
/
├── / .......................... Homepage / swap form
├── /currencies ................ Supported tokens directory (existing)
├── /how-it-works .............. Process explanation (existing)
├── /verify .................... Build verification (existing)
│
├── /case-study ................ Competitor markup analysis (existing)
│
├── /compare/ .................. Comparison hub (new)
│   ├── /compare/changenow ..... vs ChangeNow
│   ├── /compare/simpleswap .... vs SimpleSwap
│   ├── /compare/stealthex ..... vs StealthEX
│   └── /compare/exolix ........ vs Exolix
│
├── /swap/ ..................... Token pair pages (new, programmatic)
│   ├── /swap/btc-to-eth
│   ├── /swap/eth-to-usdt
│   ├── /swap/btc-to-sol
│   └── ... (top 100 pairs)
│
├── /learn/ .................... Educational content (new)
│   ├── /learn/hidden-fees-crypto-swaps
│   ├── /learn/what-is-near-intents
│   ├── /learn/cross-chain-swaps-explained
│   ├── /learn/privacy-in-crypto-trading
│   └── /learn/how-to-verify-open-source-software
│
├── /faq ....................... Frequently asked questions (new)
│
├── /fees ...................... Fee transparency page (new)
│
├── /robots.txt ................ Crawl directives (new)
├── /sitemap.xml ............... XML sitemap (new)
└── /llms.txt .................. AI crawler guidance (new)
```

## Internal Linking Strategy

### Hub-and-Spoke Model

**Hub: Homepage (`/`)**
- Links to: /currencies, /how-it-works, /fees, /compare/
- Receives links from: all pages (via header/footer nav)

**Hub: /currencies**
- Links to: individual /swap/ pair pages
- Links to: /how-it-works
- Receives links from: /swap/ pages (back to full list)

**Hub: /compare/**
- Links to: /fees, /case-study
- Links to: individual comparison pages
- Receives links from: /learn/ articles, /case-study

**Hub: /learn/**
- Links to: relevant product pages (/fees, /how-it-works, /verify)
- Links to: /compare/ pages where relevant
- Receives links from: /faq, /case-study

### Cross-Linking Rules

1. Every `/swap/` page links to `/fees` and `/how-it-works`
2. Every `/compare/` page links to `/case-study` and `/fees`
3. Every `/learn/` page links to at least 2 other internal pages
4. `/case-study` links to all `/compare/` pages
5. `/faq` links to relevant `/learn/` and `/how-it-works` pages
6. Footer navigation on all pages: Home, Currencies, How it Works, Fees, Compare, Source, Telegram

## Navigation Updates

### Current Footer
```
Currencies | How it Works | Case Study | Verify | Source | Telegram
```

### Proposed Footer
```
Currencies | How it Works | Fees | Compare | Case Study | Verify | Source | Telegram
```

## Programmatic Pages: /swap/ Token Pairs

### URL Pattern
`/swap/{from}-to-{to}` (lowercase ticker symbols)

### Template Content
- Live quote for the pair (fetched on page load)
- Supported networks for each token
- Fee breakdown (showing $0 markup)
- Link to start swap (pre-filled form)
- Mini FAQ: "How long does {FROM} to {TO} take?", "What are the fees?"

### Priority Pairs (Phase 1 — Top 20)
BTC↔ETH, BTC↔USDT, BTC↔SOL, ETH↔USDT, ETH↔SOL, BTC↔XMR, BTC↔LTC, ETH↔BNB, BTC↔DOGE, SOL↔USDT, BTC↔NEAR, ETH↔MATIC, BTC↔ADA, BTC↔DOT, ETH↔ARB, BTC↔AVAX, ETH↔OP, BTC↔TRX, BTC↔XRP, SOL↔ETH

### Expansion
- Phase 2: Top 50 pairs
- Phase 3: All pairs with meaningful search volume (100+)

## Sitemap Structure

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <!-- Core pages (high priority) -->
  <url><loc>https://zero.uswap.net/</loc><priority>1.0</priority><changefreq>daily</changefreq></url>
  <url><loc>https://zero.uswap.net/currencies</loc><priority>0.8</priority><changefreq>weekly</changefreq></url>
  <url><loc>https://zero.uswap.net/how-it-works</loc><priority>0.8</priority><changefreq>monthly</changefreq></url>
  <url><loc>https://zero.uswap.net/fees</loc><priority>0.8</priority><changefreq>monthly</changefreq></url>
  <url><loc>https://zero.uswap.net/case-study</loc><priority>0.7</priority><changefreq>monthly</changefreq></url>
  <url><loc>https://zero.uswap.net/verify</loc><priority>0.5</priority><changefreq>weekly</changefreq></url>
  <url><loc>https://zero.uswap.net/faq</loc><priority>0.6</priority><changefreq>monthly</changefreq></url>

  <!-- Comparison pages -->
  <url><loc>https://zero.uswap.net/compare/changenow</loc><priority>0.7</priority><changefreq>monthly</changefreq></url>
  <url><loc>https://zero.uswap.net/compare/simpleswap</loc><priority>0.7</priority><changefreq>monthly</changefreq></url>
  <!-- ... -->

  <!-- Token pair pages -->
  <url><loc>https://zero.uswap.net/swap/btc-to-eth</loc><priority>0.6</priority><changefreq>daily</changefreq></url>
  <!-- ... -->
</urlset>
```

## robots.txt

```
User-agent: *
Allow: /
Disallow: /api/
Disallow: /tg/
Disallow: /wrapper-logs
Disallow: /order/
Disallow: /static/

Sitemap: https://zero.uswap.net/sitemap.xml

# AI crawlers welcome
User-agent: GPTBot
Allow: /

User-agent: ClaudeBot
Allow: /

User-agent: PerplexityBot
Allow: /
```
