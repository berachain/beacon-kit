package deposit

const (
	// broadcastDeposit is the flag for broadcasting a deposit transaction.
	broadcastDeposit = "broadcast"

	// privateKey is the flag for the private key to sign the deposit message.
	privateKey = "private-key"
)

const (
	// broadcastDepositShorthand is the shorthand flag for the broadcastDeposit flag.
	broadcastDepositShorthand = "b"
)

const (
	// defaultBroadcastDeposit is the default value for the broadcastDeposit flag.
	defaultBroadcastDeposit = false

	// defaultPrivateKey is the default value for the privateKey flag.
	defaultPrivateKey = ""
)

const (
	// broadcastDepositFlagUsage is the usage description for the broadcastDeposit flag.
	broadcastDepositMsg = "broadcast the deposit transaction"

	// privateKeyFlagUsage is the usage description for the privateKey flag.
	privateKeyMsg = "private key to sign the deposit message"
)
