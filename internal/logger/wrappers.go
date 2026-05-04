package logger

func Debug(msg any, v ...any) {
	New().Debug(msg, v...)
}

func Debugf(msg string, v ...any) {
	New().Debugf(msg, v...)
}

func Info(msg any, v ...any) {
	New().Info(msg, v...)
}

func Infof(msg string, v ...any) {
	New().Infof(msg, v...)
}

func Warn(msg any, v ...any) {
	New().Warn(msg, v...)
}

func Warnf(msg string, v ...any) {
	New().Warnf(msg, v...)
}

func Error(msg any, v ...any) {
	New().Error(msg, v...)
}

func Errorf(msg string, v ...any) {
	New().Errorf(msg, v...)
}

func Fatal(msg any, v ...any) {
	New().Fatal(msg, v...)
}

func Fatalf(msg string, v ...any) {
	New().Fatalf(msg, v...)
}
