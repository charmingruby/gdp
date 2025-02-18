package logger

import "github.com/fatih/color"

var (
	highlight = color.New(color.Bold, color.BgGreen, color.FgWhite).PrintFunc()
	config    = color.New(color.Bold, color.BgYellow, color.FgWhite).PrintFunc()
)
