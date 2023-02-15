package main

import (
	"fyne.io/fyne/v2/app"

	"fyne.io/setup/pkg"
)

func main() {
	a := app.New()
	pkg.ShowSummaryWindow(a)
	a.Run()
}
