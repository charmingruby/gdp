package logger

func OpenBracket() {
	highlight("%=n")
	println()
}

func CloseBracket() {
	highlight("%=v")
	println()
}

func Divider() {
	highlight("%=>--___--=>")
	println()
}
