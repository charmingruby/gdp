package logger

import "fmt"

func ConfigKV(key, value string) {
	config(fmt.Sprintf("%s: %s", key, value))
	println()
}

func Config(content string) {
	config(content)
	println()
}
