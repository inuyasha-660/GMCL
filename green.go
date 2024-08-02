package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Green struct{}

var _ fyne.Theme = (*Green)(nil)

func (m Green) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.NRGBA{R: 0x17, G: 0x17, B: 0x18, A: 0xff}
	}

	if name == theme.ColorNameButton {
		// #227D51
		return color.NRGBA{R: 0x22, G: 0x7D, B: 0x51, A: 0xff}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m Green) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m Green) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m Green) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
