package constant

import "time"

const (
	LogErrorTag         = "[ERROR]"
	LogInfoTag          = "[INFO]"
	ConsolePrintGreen   = "\u001B[1;32;40m%s\u001B[0m\n"
	ConsolePrintCyan    = "\u001B[1;36;40m%s\u001B[0m\n"
	LogFilePath         = "./log/runtime.%Y-%m-%d.log"
	LogFileRotationTime = time.Duration(24) * time.Hour
)
