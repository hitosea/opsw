package main

import (
	"opsw/app/command"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	command.Execute()
}
