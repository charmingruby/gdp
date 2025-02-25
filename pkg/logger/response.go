package logger

func Response(content string) {
	highlight("%")
	print(" -> ")
	println(content)
}

func HighlightedErrorResponse(content string) {
	highlight("%")
	hightlightedErrResponse(" -> ")
	hightlightedErrResponse(content)
	println()
}
