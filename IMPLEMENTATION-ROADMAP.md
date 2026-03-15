# Implementation Roadmap — uSwap Zero SEO

## Phase 1 — Foundation (Weeks 1–4)

### Week 1: Technical SEO Setup
- [ ] Add `<meta name="description">` to `layout.html` template (dynamic per page)
- [ ] Add Open Graph tags (`og:title`, `og:description`, `og:image`, `og:url`, `og:type`)
- [ ] Add Twitter Card tags (`twitter:card`, `twitter:title`, `twitter:description`)
- [ ] Add `<link rel="canonical">` self-referencing tags
- [ ] Create `/robots.txt` route in `handlers.go`
- [ ] Create `/sitemap.xml` route in `handlers.go` (auto-generated from registered routes)
- [ ] Add Organization JSON-LD schema to homepage
- [ ] Add WebApplication JSON-LD schema to homepage
- [ ] Add BreadcrumbList schema to subpages

**Implementation notes:**
- All changes are in Go templates + handlers — no external dependencies needed
- Sitemap can be generated dynamically from the route table
- robots.txt is a static response from a handler

### Week 2: Fee Transparency Page
- [ ] Create `templates/fees.html` template
- [ ] Add `/fees` route handler
- [ ] Content: live fee comparison (uSwap Zero $0 vs competitor estimates)
- [ ] Add FAQ schema to fees page
- [ ] Internal links from homepage footer and how-it-works

### Week 3: FAQ Page
- [ ] Create `templates/faq.html` template
- [ ] Add `/faq` route handler
- [ ] Write 10–15 FAQs covering: fees, security, privacy, supported tokens, speed, refunds
- [ ] Add FAQPage JSON-LD schema
- [ ] Internal links from all pages via footer

### Week 4: Optimize Existing Content
- [ ] Add Article schema to `/case-study`
- [ ] Add FAQ schema to `/how-it-works`
- [ ] Improve `/currencies` page with token descriptions and network info
- [ ] Create `llms.txt` (static file or route)
- [ ] Verify Core Web Vitals with Lighthouse

**Phase 1 deliverables:** Technical SEO foundation complete, 2 new pages, all existing pages optimized.

---

## Phase 2 — Expansion (Weeks 5–12)

### Weeks 5–7: Comparison Pages
- [ ] Create `templates/compare.html` template (reusable for all competitors)
- [ ] Add `/compare/{competitor}` route handler with data-driven content
- [ ] Create comparison data structs in Go (fees, features, tracking, etc.)
- [ ] Publish: vs ChangeNow, vs SimpleSwap, vs StealthEX
- [ ] Add ComparisonPage or Article schema
- [ ] Internal link from case study and fees pages

### Weeks 8–9: Token Pair Pages (Programmatic)
- [ ] Create `templates/swap_pair.html` template
- [ ] Add `/swap/{from}-to-{to}` route handler
- [ ] Fetch live quotes on page load (server-rendered)
- [ ] Generate pages for top 20 pairs from currency list
- [ ] Add estimated swap time per pair
- [ ] Add FAQ schema per pair page
- [ ] Add to sitemap.xml dynamically

### Weeks 10–12: Educational Content
- [ ] Create `templates/learn.html` template (article layout)
- [ ] Add `/learn/{slug}` route handler
- [ ] Write and publish 3 articles (see Content Calendar)
- [ ] Add Article schema to each
- [ ] Internal linking to product pages

**Phase 2 deliverables:** 3 comparison pages, 20 token pair pages, 3 educational articles.

---

## Phase 3 — Scale (Weeks 13–24)

### Content Expansion
- [ ] Expand token pair pages to top 50, then top 100
- [ ] Write 2 more educational articles
- [ ] Add 1 more comparison page
- [ ] Update case study with fresh transaction data

### Link Building
- [ ] Submit to crypto project directories
- [ ] Pitch case study to crypto media (CoinDesk, The Block, Decrypt)
- [ ] List on NEAR Protocol ecosystem pages
- [ ] Post on open-source crypto project lists (awesome-crypto, etc.)
- [ ] Engage in relevant Reddit/forum communities with helpful answers

### Performance Monitoring
- [ ] Set up Google Search Console
- [ ] Monitor keyword rankings weekly
- [ ] Track Core Web Vitals monthly
- [ ] Review and optimize underperforming pages

**Phase 3 deliverables:** 100+ token pair pages, 5 educational articles, 4 comparison pages, initial backlinks.

---

## Phase 4 — Authority (Months 7–12)

### Thought Leadership
- [ ] Publish quarterly "Crypto Swap Fee Report" with on-chain data
- [ ] Create developer documentation for NEAR Intents integration
- [ ] Publish year-end transparency report

### Optimization
- [ ] A/B test meta descriptions for top pages
- [ ] Update all comparison pages with fresh competitor data
- [ ] Expand programmatic pages to all pairs with search volume
- [ ] Monitor and respond to AI search citations

**Phase 4 deliverables:** Authority content, ongoing optimization, sustained organic growth.

---

## Resource Requirements

| Resource | Phase 1 | Phase 2 | Phase 3 | Phase 4 |
|----------|---------|---------|---------|---------|
| Developer time | 20–30 hrs | 30–40 hrs | 20–30 hrs | 10–15 hrs |
| Content writing | 5 hrs | 15–20 hrs | 15–20 hrs | 10–15 hrs |
| Outreach | 0 hrs | 0 hrs | 10–15 hrs | 5–10 hrs |
| Tools | None (all in Go) | None | Search Console | Same |

**Total estimated effort:** ~150–200 hours over 12 months.

---

## Dependencies

- Token pair pages require the existing `/quote` API endpoint to work for server-side rendering
- Comparison pages require periodic manual review for accuracy
- Link building requires someone to do outreach
- Search Console requires DNS/domain verification

---

## Quick Wins (Implement Immediately)

These take <1 hour each and have immediate impact:

1. **Add meta descriptions** to all 6 existing pages
2. **Create robots.txt** allowing all crawlers, blocking /api/ and /tg/
3. **Add OG tags** for social sharing (use existing favicon as og:image)
4. **Add canonical tags** to prevent duplicate content issues
5. **Submit to Google Search Console** and request indexing
