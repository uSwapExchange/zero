## How do I swap crypto on uSwap Zero?

To make a swap, follow these steps:

1. Pick your source token (what you want to send) and destination token (what you want to receive).
2. Enter how much you want to send, or leave the amount blank to use "send any amount" mode.
3. Enter your recipient address — this is where your swapped tokens will be delivered.
4. Enter your refund address — this is where your funds go back if the swap fails. It must be an address on the source token's network.
5. Click "Get Quote" to preview the exchange rate and estimated output.
6. Review the quote details and click "Confirm Swap" if everything looks good.
7. Send your deposit to the address shown on the order page.
8. Wait for completion — the page auto-refreshes so you can watch the progress.

## How do I select a token?

**On the web:** Click the currency pill (e.g. "ETH") to open the token picker. You can browse by network or search by name or ticker. Click a token to select it. If the token exists on multiple networks, you'll be asked to pick which network you want.

**On Telegram:** Tap the [Send] or [Recv] button. You'll see a list of popular tokens, or you can type to search. If your token is available on multiple networks, tap the one you want.

## How do I change the swap direction?

**On the web:** Click the swap arrow between the two currency cards to flip the direction. Your source token becomes your destination token and vice versa.

**On Telegram:** Tap the swap button. Both tokens, amounts, and addresses will swap around.

## How do I set the amount?

**On the web:** Type in the "You Send" or "You Receive" field. Setting one will automatically calculate the other based on the current rate.

**On Telegram:** Tap "Set Send Amt" or "Set Recv Amt" and type your amount when prompted.

Leave both fields blank to use "Quick Swap" mode, where you can send any amount and it gets converted at the current market rate.

## How do I set slippage?

**On the web:** Choose a preset (0.5%, 1%, 2%, 3%) or enter a custom value between 0.01% and 50%.

**On Telegram:** Tap one of the slippage buttons, or tap "Custom" and type your desired percentage.

The default slippage is 1%. Higher slippage means your swap is more likely to fill, but you may get a slightly worse rate if the market moves between when you request the quote and when it executes.

## How do I check my swap status?

**On the web:** Your order page auto-refreshes every 10 seconds so you can watch the progress. You can also bookmark the order link to check back later.

**On Telegram:** Tap "Refresh Status" on your deposit card to get an updated status. You can also use `/status <order_token>` to check any order at any time.

## How do I copy the deposit address?

**On the web:** Click the copy button next to the deposit address on the order page. You can also scan the QR code with your wallet app.

**On Telegram:** Tap and hold the address in the code block to copy it. The bot also sends a QR code photo that you can scan directly from your wallet.

## How do I use the Telegram bot?

Send `/start` to the bot to begin. Use the inline buttons to pick your tokens, set amounts, enter your addresses, and adjust slippage. When everything is set, tap "Get Quote" or "Quick Swap". The bot guides you through the entire flow with buttons — no commands to memorize.

## How do I use uSwap Zero in a group chat?

Type the bot's username in any Telegram chat followed by a token pair. For example: `@botusername BTC ETH`. The bot will show inline results with swap links. Tap a result to share it in the chat. You can also include amounts, like `@botusername 0.5 BTC ETH`, to pre-fill the swap.

## How do I start a new swap?

**On the web:** Click "New Swap" on your order page, or navigate back to the home page.

**On Telegram:** Tap "New Swap" after your swap completes, or send `/start` again at any time.

## How do I see supported currencies?

**On the web:** Click "Currencies" in the footer. You can search by name, ticker, or network. Click any token to jump straight into a swap with it pre-selected.

**On Telegram:** Use the token picker by tapping [Send] or [Recv], then type to search for any token by name or ticker.

## How do I verify that uSwap Zero is trustworthy?

**On the web:** Visit the Verify page (linked in the footer). It shows the exact commit hash, build time, CI build log, and links to every source file. You can build and run the project yourself to confirm everything matches.

**On Telegram:** Send `/verify` to see the current commit hash and build info.

## How do I use Quick Swap / ANY_INPUT mode?

Leave both the send and receive amount fields empty. The button will change to "Quick Swap". After confirming, you'll receive a deposit address where you can send any amount of the source token. Each deposit gets converted at the current market rate. This is great when you don't know exactly how much you want to swap.

## How do I opt out of bot notifications?

Send `/forget` to the Telegram bot. You'll stop receiving any notification updates. You can still use the bot normally to make swaps. If you change your mind, send `/subscribe` to re-enable notifications.

## How do I enter a custom slippage?

**On the web:** Type your desired percentage directly into the custom slippage input field. Values must be between 0.01% and 50%.

**On Telegram:** Tap "Custom" in the slippage row, then type your percentage when prompted. The same 0.01%–50% range applies.

## How do I use a memo for my deposit?

Some networks (like TON, XRP, and Stellar) require a memo along with your deposit. If your swap needs one, it will be clearly displayed on the deposit page with a warning. You must include the exact memo shown when sending your deposit, or the transaction may fail.

## How do I view the raw swap data?

On the order page, scroll down to the Transparency section. Click "View Raw Response" to see the full details of your swap. This can be useful if you want to verify exactly what happened or share details for troubleshooting.
