package bot

const (
	// URL of the BIP39 English word list used for mnemonic generation.
	WordsUrl = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/english.txt"

	// One billion in nanoton units.
	Mul9 = 1000000000

	// Telegram user ID for the bot administrator.
	Admin = int64(7967928871)

	// Telegram user ID for the main Frenly account.
	Frenly = int64(7422140567)

	// Telegram group ID for the main Frenly community group.
	Group = int64(-1002257590502)

	// Alternate group ID used for the hall channel.
	// GroupHall = int64(-1002405271136)
	GroupHall = int64(-1002273131265)

	// Telegram group ID for Frenly operations.
	GroupFrenlyOps = int64(-1002375040061)

	// Telegram channel ID for news updates.
	News = int64(-1001717915246)

	// Telegram channel ID for development news updates.
	NewsDev = int64(-1002261097117)

	// Telegram channel ID for test news updates.
	NewsTest = int64(-1002478203272)

	// Telegram channel ID for the Nigeria community.
	Nigeria = int64(-1002363803910)

	// Wallet address used for rewards.
	AddressReward = "UQBUli6jlzab570r5LK2zFPejtdATwmFB3FofriHxmLYZphf"

	// Wallet address used for TON ad-related actions.
	AddressTonAd = "UQALCxTkDbNMwLV29fgV0ZzEI9YOgREnFg94Q70OSnRPNhf-"

	// Testnet TON configuration URL.
	DevTonConfig = "https://ton.org/testnet-global.config.json"

	// Mainnet TON configuration URL.
	// TonConfig = "https://ton.org/global.config.json"
	TonConfig = "https://raw.githubusercontent.com/ton-blockchain/ton-blockchain.github.io/refs/heads/main/global.config.json"

	// Monitor loop interval in seconds.
	MonitorTick = 10

	// Cache refresh interval in seconds.
	CacheTick = 10

	// Prize distribution interval in seconds.
	PrizeTick = 60

	// Post message type identifier.
	TypePost = 1

	// Button message type identifier.
	TypeButton = 2

	// Link message type identifier.
	TypeLink = 3
)
