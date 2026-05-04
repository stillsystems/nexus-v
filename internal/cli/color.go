package cli

import "fmt"

var (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func Info(msg string) {
	fmt.Println(Blue + "ℹ" + Reset + " " + msg)
}

func Success(msg string) {
	fmt.Println(Green + "✔" + Reset + " " + msg)
}

func Warn(msg string) {
	fmt.Println(Yellow + "⚠" + Reset + " " + msg)
}

func Error(msg string) {
	fmt.Println(Red + "✖ " + msg + Reset)
}

