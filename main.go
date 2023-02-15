package main

import (
	"fyne.io/fyne/v2/app"

	"fyne.io/setup/pkg"
)

func main() {
	a := app.NewWithID("io.fyne.setup")
	pkg.ShowSummaryWindow(a)
	a.Run()
}
