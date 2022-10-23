package logger

func LogLevelFromString(level string) int {
	switch level {
	case LevelDebug:
		return -1
	case LevelInfo:
		return 0
	case LevelWarn:
		return 1
	case LevelPanic:
		return 4
	case LevelFatal:
		return 5
	default:
		return 0
	}
}
