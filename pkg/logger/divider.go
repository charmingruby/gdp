package logger

func OpenBracket() {
	highlight("%=n")
	println()
}

func CloseBracket() {
	highlight("%=v")
}

func Divider() {
	println()
}
