package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Bili_Pink struct {
}

var _ fyne.Theme = (*Bili_Pink)(nil)

func (m Bili_Pink) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.NRGBA{R: 0x17, G: 0x17, B: 0x18, A: 0xff}
	}

	if name == theme.ColorNameButton {
		// #FF6699
		return color.NRGBA{R: 0xFF, G: 0x66, B: 0x99, A: 0xff}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m Bili_Pink) Icon(name fyne.ThemeIconName) fyne.Resource {

	return theme.DefaultTheme().Icon(name)
}

func (m Bili_Pink) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m Bili_Pink) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
