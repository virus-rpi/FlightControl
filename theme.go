package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type FlightControlTheme struct{}

func (t *FlightControlTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *FlightControlTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *FlightControlTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == "text" {
		if fyne.CurrentDevice().IsMobile() {
			return 24
		}
		return 12
	}
	return theme.DefaultTheme().Size(name)
}

func (t *FlightControlTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}
