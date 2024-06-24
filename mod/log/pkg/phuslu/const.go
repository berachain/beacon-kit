package phuslu

// colours.
const (
	reset      = "\x1b[0m"
	black      = "\x1b[30m"
	red        = "\x1b[31m"
	green      = "\x1b[32m"
	yellow     = "\x1b[33m"
	blue       = "\x1b[34m"
	magenta    = "\x1b[35m"
	cyan       = "\x1b[36m"
	white      = "\x1b[37m"
	gray       = "\x1b[90m"
	lightWhite = "\x1b[97m"
)

// log levels.
const (
	traceColor   = magenta
	debugColor   = yellow
	infoColor    = green
	warnColor    = yellow
	errorColor   = red
	fatalColor   = red
	panicColor   = red
	defaultColor = gray
	traceLabel   = "TRCE"
	debugLabel   = "DBUG"
	infoLabel    = "INFO"
	warnLabel    = "WARN"
	errorLabel   = "ERRR"
	fatalLabel   = "FATAL"
	panicLabel   = "PANIC"
	defaultLabel = " ???"
)
