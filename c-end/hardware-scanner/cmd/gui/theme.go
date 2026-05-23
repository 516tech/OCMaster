package main

import (
	"image/color"

	g "github.com/AllenDang/giu"
)

// ImHexDark matches ImHex's color palette
// WindowBg: #11161B, FrameBg: #26313C, Tab: #1A2028, Accent: #2864C8
func ImHexDark() *g.StyleSetter {
	return g.Style().
		SetColor(g.StyleColorText, color.RGBA{200, 200, 200, 255}).
		SetColor(g.StyleColorTextDisabled, color.RGBA{100, 100, 100, 255}).
		SetColor(g.StyleColorWindowBg, color.RGBA{17, 22, 27, 255}).
		SetColor(g.StyleColorChildBg, color.RGBA{22, 28, 35, 255}).
		SetColor(g.StyleColorPopupBg, color.RGBA{30, 37, 45, 255}).
		SetColor(g.StyleColorBorder, color.RGBA{50, 58, 68, 255}).
		SetColor(g.StyleColorFrameBg, color.RGBA{38, 49, 60, 255}).
		SetColor(g.StyleColorFrameBgHovered, color.RGBA{50, 64, 78, 255}).
		SetColor(g.StyleColorFrameBgActive, color.RGBA{30, 39, 48, 255}).
		SetColor(g.StyleColorTitleBg, color.RGBA{14, 18, 23, 255}).
		SetColor(g.StyleColorTitleBgActive, color.RGBA{17, 22, 27, 255}).
		SetColor(g.StyleColorMenuBarBg, color.RGBA{20, 26, 32, 255}).
		SetColor(g.StyleColorScrollbarBg, color.RGBA{17, 22, 27, 255}).
		SetColor(g.StyleColorScrollbarGrab, color.RGBA{60, 70, 85, 255}).
		SetColor(g.StyleColorScrollbarGrabHovered, color.RGBA{80, 92, 108, 255}).
		SetColor(g.StyleColorScrollbarGrabActive, color.RGBA{100, 112, 128, 255}).
		SetColor(g.StyleColorButton, color.RGBA{38, 49, 60, 255}).
		SetColor(g.StyleColorButtonHovered, color.RGBA{40, 100, 200, 255}).
		SetColor(g.StyleColorButtonActive, color.RGBA{30, 80, 160, 255}).
		SetColor(g.StyleColorHeader, color.RGBA{38, 49, 60, 255}).
		SetColor(g.StyleColorHeaderHovered, color.RGBA{50, 64, 78, 255}).
		SetColor(g.StyleColorHeaderActive, color.RGBA{55, 71, 87, 255}).
		SetColor(g.StyleColorSeparator, color.RGBA{50, 58, 68, 255}).
		SetColor(g.StyleColorTab, color.RGBA{26, 32, 40, 255}).
		SetColor(g.StyleColorTabHovered, color.RGBA{40, 100, 200, 200}).
		SetColor(g.StyleColorTabActive, color.RGBA{38, 49, 60, 255}).
		SetColor(g.StyleColorTableHeaderBg, color.RGBA{30, 37, 45, 255}).
		SetColor(g.StyleColorTableBorderStrong, color.RGBA{50, 58, 68, 255}).
		SetColor(g.StyleColorTableBorderLight, color.RGBA{38, 44, 52, 255}).
		SetColor(g.StyleColorTableRowBg, color.RGBA{22, 28, 35, 255}).
		SetColor(g.StyleColorTableRowBgAlt, color.RGBA{28, 34, 42, 255}).
		SetColor(g.StyleColorTextSelectedBg, color.RGBA{40, 100, 200, 150}).
		SetStyleFloat(g.StyleVarFrameRounding, 3).
		SetStyleFloat(g.StyleVarWindowRounding, 5).
		SetStyleFloat(g.StyleVarScrollbarRounding, 3).
		SetStyleFloat(g.StyleVarGrabRounding, 3).
		SetStyleFloat(g.StyleVarTabRounding, 4)
}
