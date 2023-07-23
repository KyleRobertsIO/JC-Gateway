package main

import "kyleroberts.io/src/app"

func main() {
	app := app.Application{
		Name: "Job Container Manager",
	}
	app.Start()
}
