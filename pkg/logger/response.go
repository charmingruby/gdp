package logger

func Response(content string) {
	highlight("%")
	print(" -> ")
	println(content)
}
