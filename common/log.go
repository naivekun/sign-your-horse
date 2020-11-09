package common

func LogWithModule(name, msg string) string {
	return "[" + name + "]" + msg
}
