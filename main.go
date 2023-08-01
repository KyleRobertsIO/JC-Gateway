package main

import "github.com/kylerobertsio/aci-job-manager/src/app"

func main() {
	app := app.Application{
		Name: "Job Container Manager",
	}
	app.Start()
}
