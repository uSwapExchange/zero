# Something Went Wrong — uSwap Zero Support

Troubleshooting common errors and unexpected situations on uSwap Zero.

---

## My swap is stuck or taking too long

This is the most common concern we hear, and we completely understand the anxiety. Most swaps complete within a few minutes, but blockchain congestion can sometimes cause delays.

Your order page auto-refreshes every 10 seconds, so you don't need to keep reloading. If your swap shows "Processing," that's actually good news — it means your deposit was detected and the swap is actively being executed. Just hang tight.

If it's been more than 30 minutes with no progress, something may have gone wrong on the network side. In that case, your funds will be automatically refunded to the refund address you provided when you created the swap. No action is needed on your part.

---

## Why does it say "Expired" on my order?

Each swap has a deposit window — usually about 1 hour — during which you need to send your funds. If you didn't send your deposit before that window closed, the deposit address is no longer valid and the order is marked as expired.

Don't worry, no funds were lost. Simply go back to the home page, start a new swap, and make sure to send your deposit promptly after the order is created.

---

## I'm getting "Quote Failed"

This means the swap service is temporarily unavailable. It's almost always a brief hiccup that resolves on its own within a few minutes.

Head back to the home page and try your swap again. If the problem keeps happening, the service may be experiencing unusually high demand. Give it a few minutes and try once more.

---

## I'm getting "Quote Unavailable" or "No market makers"

This means no one is currently offering a rate for your specific token pair or amount. Here are a few things to try:

1. **Try a larger amount** — very small swaps sometimes don't have enough liquidity.
2. **Try a different token pair** — some pairs are more actively traded than others.
3. **Wait a few minutes and try again** — availability can change quickly.

---

## I'm getting "Too Many Requests"

You've sent too many requests in a short period. This limit is in place to keep the service running smoothly for everyone.

Just wait about a minute and try again. You should be good to go after a brief pause.

---

## I'm getting "Form expired. Please go back and try again."

Your session on the quote page expired after being inactive for about an hour. This is a security measure to make sure you're always getting a fresh, accurate quote.

Simply go back to the home page and start a new quote. It only takes a moment.

---

## I'm getting "Invalid Order" or "This order link is invalid or expired"

The order link you're using is no longer valid. This can happen if the link was generated during a previous session or server update.

Go to the home page and start a new swap. You'll get a fresh, valid order link.

---

## I'm getting "Unknown Token"

The token you selected is no longer available in our system. Token availability can change as markets evolve.

Go back to the home page and select your tokens again from the currency picker. The list will show all currently supported tokens.

---

## I'm getting "Invalid Amount"

The amount you entered couldn't be understood. Make sure you're entering a plain number — for example, `0.5` or `100`.

Don't include currency symbols, letters, or special characters. Just the number.

---

## My swap was refunded

We're sorry your swap couldn't be completed. The good news is that your funds are automatically returned to the refund address you provided. Common reasons for a refund include:

- Price movement exceeded your slippage tolerance.
- Temporary network issues.
- Liquidity problems for your token pair.

Your order page will show the reason if one is available. Check your refund wallet — the funds should appear there shortly, depending on blockchain confirmation times.

---

## My swap shows "Incomplete Deposit"

This means the deposit couldn't be fully processed. If you sent funds, they will be returned to your refund address automatically — no action needed from you.

If you'd like to try again, simply start a new swap from the home page.

---

## I can't find my token in the list

Use the search bar in the token picker to search by name or ticker symbol. uSwap Zero supports 140+ tokens across 29+ blockchains, so there's a good chance your token is there.

If you still can't find it, that token isn't currently supported. You can check the Currencies page for the full list of available tokens.

---

## My deposit failed because I forgot the memo

Some networks — including TON, XRP, and Stellar — require a memo along with your deposit. If you sent your deposit without the required memo, it may not be detected by the system.

Unfortunately, uSwap Zero cannot recover deposits sent without the correct memo. You'll need to contact the support channels for the specific blockchain network to seek help recovering those funds. We know this is frustrating, and we're sorry we can't do more here.

---

## I'm getting "Address seems too short"

The wallet address you entered appears to be too short (less than 10 characters). This usually means the address was cut off when you copied it.

Double-check that you're pasting the complete wallet address for the correct blockchain, and try again.

---

## I'm getting "Invalid slippage"

Slippage must be a number between 0.01 and 50. Just enter the number — no need to include a `%` sign.

For most swaps, the default slippage works fine. Only adjust it if you have a specific reason to.

---

## The bot says "Unknown command"

The bot didn't recognize what you typed. Here are the commands you can use:

- **/start** — Begin a new swap
- **/status** — Check an existing order
- **/verify** — See build and version info
- **/forget** — Opt out and remove your data
- **/subscribe** — Opt back in after opting out

---

## Status check failed

The status check couldn't reach the swap service right now. This is almost always temporary.

Tap **"Refresh Status"** again in a moment. You can also check your order on the web by tapping the **"Open Order"** button in the bot message.

---

## My swap is showing "Processing" for a long time

"Processing" means your deposit was successfully received and the swap is being executed on the blockchain. Depending on network congestion, this step can take several minutes.

Please be patient — the system is working on it. If the swap ultimately can't be completed, your funds will be refunded automatically to your refund address. You don't need to do anything.

---

## The page shows "NEAR Intents API is temporarily unavailable"

The underlying swap service is experiencing a brief outage. These are usually resolved within a few minutes.

Go back to the home page and try your swap again shortly. We appreciate your patience.
