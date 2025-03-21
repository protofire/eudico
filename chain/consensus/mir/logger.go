package mir

import (
	ipfslogging "github.com/ipfs/go-log/v2"

	mirlogging "github.com/filecoin-project/mir/pkg/logging"
)

const managerLoggerName = "mir-manager"

var _ mirlogging.Logger = &Logger{}

// Logger implements Mir's Log interface.
type Logger struct {
	logger *ipfslogging.ZapEventLogger
	id     string
}

func NewLogger(id string) *Logger {
	return &Logger{
		logger: ipfslogging.Logger(managerLoggerName),
		id:     id,
	}
}

// Log logs a message with additional context.
func (l *Logger) Log(level mirlogging.LogLevel, text string, args ...interface{}) {
	// adding mirID to logs.
	args = append(args, []interface{}{"nodeID", l.id}...)

	switch level {
	case mirlogging.LevelError:
		l.logger.Errorw(text, args...)
	case mirlogging.LevelInfo:
		l.logger.Infow(text, args...)
	case mirlogging.LevelWarn:
		l.logger.Warnw(text, args...)
	case mirlogging.LevelDebug:
		l.logger.Debugw(text, args...)
	}
}

func (l *Logger) MinLevel() mirlogging.LogLevel {
	level := ipfslogging.GetConfig().SubsystemLevels[managerLoggerName]
	switch level {
	case ipfslogging.LevelDebug:
		return mirlogging.LevelDebug
	case ipfslogging.LevelInfo:
		return mirlogging.LevelInfo
	case ipfslogging.LevelWarn:
		return mirlogging.LevelWarn
	case ipfslogging.LevelError:
		return mirlogging.LevelError
	case ipfslogging.LevelDPanic:
		return mirlogging.LevelError
	case ipfslogging.LevelPanic:
		return mirlogging.LevelError
	case ipfslogging.LevelFatal:
		return mirlogging.LevelError
	default:
		return mirlogging.LevelError
	}
}

func (l *Logger) IsConcurrent() bool {
	return true
}
