package main

import (
	"opsw/command"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	command.Execute()
}
