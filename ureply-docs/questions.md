# uSwap Zero — Support FAQ

## Can I swap without creating an account?

Yes. uSwap Zero requires no account, no sign-up, and no KYC. Just pick your tokens, enter your wallet addresses, and swap. There is nothing to register, no email to verify, and no personal information collected. You stay in full control of your funds the entire time.

## Can I swap Bitcoin to Ethereum?

Yes. uSwap Zero supports cross-chain swaps between 29+ blockchains. You can swap BTC to ETH, SOL to USDT, or any supported pair. Just select the source and destination tokens using the currency picker, paste your addresses, and go.

## What happens if my swap fails?

Your funds are automatically returned to the refund address you provided when setting up the swap. This happens if the swap times out, the price moves beyond your slippage tolerance, or there's a network issue. No action is needed on your part — just make sure the refund address you entered is correct before confirming.

## How long does a swap take?

Most swaps complete within 2–10 minutes, depending on the blockchains involved. Bitcoin swaps may take longer due to block confirmation times. The order page auto-refreshes so you can watch the progress. If your swap seems stuck, give it a few extra minutes before worrying.

## Is there a minimum swap amount?

There's no hard minimum set by uSwap Zero, but very small amounts may not get quotes from market makers. If you see "No market makers," try a slightly larger amount. The threshold varies by token and current market conditions.

## Is there a maximum swap amount?

There's no maximum set by uSwap Zero. Very large amounts depend on market maker liquidity. If you can get a quote, you can swap. For especially large trades, you may want to split into multiple swaps to get better rates.

## What is slippage?

Slippage is the maximum price change you'll accept between when you get a quote and when the swap executes. Higher slippage means your swap is more likely to fill, but you might get a slightly worse rate if the market moves. The default is 1%. You can adjust it using the custom slippage option if you need tighter or looser tolerance.

## What is the market maker spread?

The spread is the difference between the spot price and the rate offered by the market maker. It's not a fee that uSwap Zero charges — it's simply the cost of the trade on the open market. uSwap Zero adds zero markup on top of this. What you see in the quote is what the market is offering.

## What does "Quick Swap" or ANY_INPUT mean?

If you leave both amount fields empty, you enter Quick Swap mode. You'll get a deposit address where you can send any amount of the source token. Each deposit gets converted at the current market rate. This is great for when you don't know the exact amount in advance, or you want to send multiple deposits to the same address.

## Do I need JavaScript enabled?

No. The web interface works entirely without JavaScript. This is a privacy feature — there's no client-side tracking, analytics, or third-party scripts. Everything is rendered server-side. The page works the same whether JS is on or off.

## Can I use uSwap Zero on mobile?

Yes. The web interface works in any mobile browser. You can also use the Telegram bot, which provides the full swap experience through buttons and messages. Both options are fully functional on mobile devices.

## What is the refund address for?

The refund address is where your funds are returned if the swap can't be completed. It must be a valid address on the same blockchain you're sending from. Always double-check this address — if you enter the wrong one, refunded funds could be lost permanently.

## What is the recipient address?

The recipient address is where your swapped tokens will be delivered. It must be a valid address on the destination blockchain. Always double-check this address before confirming your swap. A wrong recipient address means your funds go to someone else with no way to reverse it.

## Can I cancel a swap after confirming?

Once you confirm a swap and receive a deposit address, you can cancel by simply not sending the deposit. If you've already sent the deposit, the swap will proceed. If it can't complete for any reason, funds are automatically refunded to your refund address.

## What happens if I send the wrong amount?

If you specified an exact send amount but send a different amount, the swap may still process depending on the mode. In Quick Swap (ANY_INPUT) mode, any amount works. In other modes, sending the wrong amount may result in a refund to your refund address.

## How do I know my swap completed?

The order page will show "Swap Complete" with a checkmark and a link to view the transaction on the blockchain explorer. In the Telegram bot, tap "Refresh Status" to see the latest state. You can also check your destination wallet for the incoming funds.

## Can I swap the same token on different networks?

Yes. For example, you can swap USDT on Ethereum to USDT on Solana. Select the same token but different networks in the currency picker. This is a common use case for bridging assets between chains without paying high bridge fees.

## What is the deadline on my deposit?

Each swap has a deposit deadline, usually about 1 hour. You must send your deposit before this deadline. After it expires, the deposit address is no longer valid and you need to create a new swap. The deadline is clearly shown on the order page.

## Is my data private?

uSwap Zero collects no personal data. No accounts, no cookies, no JavaScript tracking. The web UI works without JS. The Telegram bot stores your session temporarily and clears it after 2 hours of inactivity. You can send /forget to the bot at any time to be removed from any notification list.

## Can I track an old swap?

On the web, if you have the order link, you can check the status anytime — just open it in your browser. On the Telegram bot, use /status followed by your order token to check any swap. Bookmark or save your order links if you want to check back later.

## What networks require a memo?

Some blockchains like TON, XRP (Ripple), Stellar, and NEAR require a memo with deposits. If your swap needs a memo, it will be clearly displayed on the deposit page with a warning. You must include it exactly as shown, or your deposit may fail or be lost.

## Is uSwap Zero open source?

Yes. The full source code is available on GitHub. You can audit every line of code, build it yourself, and verify the running instance matches the source. Transparency is a core value of the project.

## What does the Verify page show?

The Verify page shows the exact Git commit, build time, CI build log, and links to every source file. You can confirm that the running server matches the published source code. This is how you know the service hasn't been tampered with.

## Can I use uSwap Zero via Tor?

If an onion address is configured, you can access uSwap Zero over Tor for maximum privacy. Check the footer of the website for the .onion link. Combined with the no-JS design, Tor usage provides strong anonymity.

## How do I report a problem?

Join the Telegram group at t.me/uSwapZero for support, updates, and help. Describe your issue and include your order link if you have one — it helps us look into what happened. Do not share your private keys or seed phrases with anyone.

## What currencies are supported?

uSwap Zero supports 140+ tokens across 29+ blockchains including Bitcoin, Ethereum, Solana, Base, Arbitrum, TON, TRON, BNB Chain, Polygon, Optimism, Avalanche, NEAR, and more. Visit the Currencies page on the website for the full live list, updated automatically.

## Why did I get a different amount than the quote showed?

The quote shows an estimated receive amount. The actual amount may differ slightly due to market price movement between quoting and execution. Your slippage setting limits how much the price can move — if it moves more than your slippage tolerance, the swap is refunded instead of executing at a bad rate.

## Can I do multiple swaps at the same time?

On the web, each swap gets its own order page with a unique link, so you can run multiple swaps in different tabs. On the Telegram bot, you work with one swap at a time per chat. There's no limit to how many swaps you can have in progress.

## My swap has been pending for a long time. What should I do?

First, check the order page — it auto-refreshes with the latest status. Most delays are caused by slow block confirmations on the source chain. Bitcoin swaps can take 30+ minutes during network congestion. If your swap is still pending after an hour, join t.me/uSwapZero and share your order link so the community can help.

## Why does the quote say "No market makers"?

This means no market maker is currently willing to fill your trade at the given parameters. Common causes: the amount is too small, the token pair has low liquidity, or the market is very volatile. Try adjusting the amount, increasing your slippage tolerance, or waiting a few minutes and trying again.

## Are there any fees?

uSwap Zero charges zero fees. None. The only cost is the market maker spread (the difference between spot price and the offered rate) and the network transaction fees on the source and destination blockchains. These are not charged by uSwap Zero — they're inherent costs of trading and using blockchains.

## What if I entered the wrong recipient address?

If you haven't sent the deposit yet, simply start a new swap with the correct address. If you already sent the deposit, the swap will deliver to whatever address you entered. Blockchain transactions are irreversible, so always double-check addresses before sending.

## How does uSwap Zero make money if there are no fees?

uSwap Zero is a free, open-source project. It does not charge fees or take a cut. It exists to provide a private, zero-cost swap option for the crypto community.
