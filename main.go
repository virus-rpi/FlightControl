package main

import (
	"FlightControl/warp"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	App := app.NewWithID("com.virusrpi.flightcontrol")
	App.Settings().SetTheme(&FlightControlTheme{})
	go warp.InitRocketClient(App, ps)
	go listenToNewData()
	MainWindow := App.NewWindow("Flight Control")
	MainWindow.Resize(fyne.NewSize(800, 600))
	MainWindow.CenterOnScreen()

	tabControl := container.NewTabItem("Control", controlTab(App, MainWindow))
	tabAnalysis := container.NewTabItem("Analysis", analysisTab())
	tabSimulation := container.NewTabItem("Simulation", simulationTab())
	tabSetting := container.NewTabItem("Settings", widget.NewLabel("Content of Tab 4"))
	tabChecklists := container.NewTabItem("Checklists", widget.NewLabel("Content of Tab 5"))
	tabMock := container.NewTabItem("Mock", mockTab())

	tabs := container.NewAppTabs(tabControl, tabAnalysis, tabSimulation, tabSetting, tabChecklists, tabMock)

	tabs.OnSelected = func(item *container.TabItem) {
		ps.Pub(item, "selectedTab")
	}

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
					App.Preferences().SetString("RocketAddress", ipEntry.Text)
					go warp.RefreshRocketClient(App)
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
	defer ps.Shutdown()
}
