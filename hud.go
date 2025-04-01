package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type HUD struct {
	topColor    color.Color
	bottomColor color.Color
}

func NewHUD() *HUD {
	return &HUD{
		topColor:    color.White,
		bottomColor: color.Black,
	}
}

func (h *HUD) Draw(screen *ebiten.Image, whiteScore, blackScore int) {
	topHUD := ebiten.NewImage(WindowSize, HUDHeight)
	topHUD.Fill(h.topColor)
	screen.DrawImage(topHUD, nil)
	text.Draw(screen, fmt.Sprintf("White: %d", whiteScore), basicfont.Face7x13, 10, 20, color.Black)

	bottomHUD := ebiten.NewImage(WindowSize, HUDHeight)
	bottomHUD.Fill(h.bottomColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, WindowSize+HUDHeight)
	screen.DrawImage(bottomHUD, op)
	text.Draw(screen, fmt.Sprintf("Black: %d", blackScore), basicfont.Face7x13, 10, WindowSize+HUDHeight+20, color.White)
}

