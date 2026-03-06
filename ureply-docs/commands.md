# Telegram Bot Commands

## What commands are available?

The uSwap Zero Telegram bot supports these commands:

- **/start** — Begin a new swap or reset your session
- **/status** — Check the status of an existing swap
- **/verify** — View the bot's build and deployment info
- **/forget** — Opt out of notifications and be forgotten
- **/subscribe** — Re-enable notifications after opting out

You can also use inline queries by typing **@botusername** in any chat to search for swaps or check order status without opening a direct conversation.

## What does /start do?

Send **/start** to begin a new swap. The bot shows a swap card with buttons to pick tokens, set amounts, enter addresses, and adjust slippage. You can also use /start with a deep link to pre-fill the swap — for example, from an inline query result or a shared link. If you already have a swap in progress, /start resets your session and starts fresh.

## What does /status do?

Use **/status** followed by an order token to check the status of any swap. For example: `/status abc123def456`. The bot will show the current state of that swap — pending, processing, complete, refunded, or failed. If you type just /status alone, the bot will remind you of the correct format and ask for the order token.

## What does /verify do?

Shows the bot's deployment information: the Git commit hash (linked to GitHub), build time, and a link to the Verify page where you can audit the full source code. Use this to confirm the bot is running the published open-source code. It's a transparency feature so you never have to take our word for it.

## What does /forget do?

Removes you from any notification or update lists. You'll stop receiving messages from the bot unless you interact with it directly. Your data is forgotten. You can still use the bot normally for swaps — this only affects notifications. It's a privacy feature. Use **/subscribe** to opt back in later if you change your mind.

## What does /subscribe do?

Re-enables notifications after you've used /forget. You'll receive occasional important updates from the bot. Use **/forget** again anytime to opt back out — there's no penalty for changing your preference.

## How do I use inline queries?

Type **@botusername** in any Telegram chat (group or private) to get swap suggestions. You can add token names — for example, `@botusername BTC ETH` shows BTC-to-ETH swap options. Add an amount for price estimates: `@botusername 0.5 BTC ETH`. Tap a result to share a swap link in the chat. You can also check order status inline: `@botusername status <token>`.

## What do the buttons on the swap card do?

- **[Send] / [Recv]** — Open the token picker to choose your source or destination token
- **Swap arrow** (circular arrows icon) — Flip the send and receive tokens
- **Slippage buttons** (0.5%, 1%, 2%, 3%, Custom) — Set your slippage tolerance
- **Set Send Amt / Set Recv Amt** — Enter how much you want to send or receive
- **Set Refund Address** — Enter your refund wallet address (on the source chain)
- **Set Receive Address** — Enter your destination wallet address
- **Get Quote** — Preview the swap rate before confirming
- **Quick Swap** — Start a flexible swap where you can send any amount
- **Open Quote in App** — Open the swap in the web interface for a full-screen view

## What does "Refresh Status" do?

After you've confirmed a swap and sent your deposit, tap **Refresh Status** to check the latest state. The bot fetches the current status and updates the card in place. You'll see whether your swap is still pending, processing, or complete. You can tap it as many times as you like.

## What does "Clear" do?

After a swap is finished (complete, refunded, or failed), the **Clear** button deletes all swap-related messages from the chat and resets your session. It tidies up the conversation so you're not left with old swap cards cluttering things up.

## What does "New Swap" do?

After a swap is finished, tap **New Swap** to start a fresh swap with default settings. Your previous session is cleared and a new swap card appears. It's the quickest way to go again without typing /start.

## What does "Open Order" do?

Opens your swap's order page in the web interface. You can view full details including the deposit QR code, transaction links, and transparency data. Useful if you want to share the order or view it on a larger screen.
