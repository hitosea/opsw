package main

import (
	"opsw/app"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app.Execute()
}
