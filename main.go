package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	App := app.NewWithID("com.virusrpi.flightcontrol")
	App.Settings().SetTheme(&FlightControlTheme{})
	go func() { initWebsocket(App) }()
	MainWindow := App.NewWindow("Flight Control")
	MainWindow.Resize(fyne.NewSize(800, 600))
	MainWindow.CenterOnScreen()

	tabControl := container.NewTabItem("Control", controlTab(App, MainWindow))
	tabAnalysis := container.NewTabItem("Analysis", widget.NewLabel("Content of Tab 2"))
	tabSimulation := container.NewTabItem("Simulation", NewThreeDWidget())
	tabSetting := container.NewTabItem("Settings", widget.NewLabel("Content of Tab 4"))
	tabChecklists := container.NewTabItem("Checklists", widget.NewLabel("Content of Tab 5"))

	tabs := container.NewAppTabs(tabControl, tabAnalysis, tabSimulation, tabSetting, tabChecklists)

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Load log", func() { println("Load log") }),
			fyne.NewMenuItem("Export log", func() { println("Export log") }),
			fyne.NewMenuItem("Set WaRa IP", func() {
				ipEntry := widget.NewEntry()
				ipEntry.SetPlaceHolder("Enter IP")
				dialog.ShowForm("Set WaRa IP", "OK", "Cancel", []*widget.FormItem{
					widget.NewFormItem("IP", ipEntry),
				}, func(ok bool) {
					if !ok {
						return
					}
					App.Preferences().SetString("WaRaIP", ipEntry.Text)
					updateWebsocket(App)
				}, MainWindow)
			}),
		),
		fyne.NewMenu("Options",
			fyne.NewMenuItem("Toggle fullscreen", func() { MainWindow.SetFullScreen(!MainWindow.FullScreen()) }),
		),
	)

	if !fyne.CurrentDevice().IsMobile() {
		MainWindow.SetMainMenu(mainMenu)
	}
	MainWindow.SetContent(tabs)
	MainWindow.ShowAndRun()
}
