package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Dark struct{}

var _ fyne.Theme = (*Dark)(nil)

func (m Dark) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.NRGBA{R: 0x17, G: 0x17, B: 0x18, A: 0xff}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m Dark) Icon(name fyne.ThemeIconName) fyne.Resource {

	return theme.DefaultTheme().Icon(name)
}

func (m Dark) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m Dark) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
