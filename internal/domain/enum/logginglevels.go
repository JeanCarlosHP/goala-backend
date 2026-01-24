package enum

type LoggingLevel string

const (
	LoggingLevelDebug   LoggingLevel = "debug"
	LoggingLevelInfo    LoggingLevel = "info"
	LoggingLevelWarning LoggingLevel = "warning"
	LoggingLevelError   LoggingLevel = "error"
)

func (l LoggingLevel) String() string {
	return string(l)
}
