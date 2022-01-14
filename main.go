package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	initialXOffset       = 50
	initialYOffset       = 50
	distanceBetweenViews = 75
	minimapOffset        = 150
)

var (
	backgroundColor = color.NRGBA{0x10, 0x10, 0x10, 0xff}
	ecgOffsetX      = 0
	healthECGViews  = [5]HealthECGView{
		NewHealthECGFine(),
		NewHealthECGYellowCaution(),
		NewHealthECGOrangeCaution(),
		NewHealthECGDanger(),
		NewHealthECGPoison(),
	}
)

func init() {
}

type Game struct{}

func (g *Game) Update() error {
	ecgOffsetX = (ecgOffsetX + 1) % 128
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	for i := 0; i < len(healthECGViews); i++ {
		DrawECGView(screen, healthECGViews[i], initialXOffset, initialYOffset+distanceBetweenViews*i)
		DrawMinimap(screen, healthECGViews[i], initialXOffset, initialYOffset+distanceBetweenViews*i)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func DrawECGView(screen *ebiten.Image, ecgView HealthECGView, xOffset int, yOffset int) {
	for columnNum := 0; columnNum < 32; columnNum++ {
		startX := ecgOffsetX - columnNum
		if startX < 0 || startX >= len(ecgView.Lines) {
			continue
		}
		// Draw a vertical line for the current line
		destX := startX + xOffset
		destY := ecgView.Lines[startX][0] + yOffset
		width := 1
		height := ecgView.Lines[startX][1] + 1

		// Lines to the left will have a darker color
		// If the color component is negative, set to 0 so that it will render correctly
		lineColor := ecgView.Color
		gradientColor := ecgView.Gradient
		red := lineColor[0] - (gradientColor[0] * columnNum)
		if red < 0 {
			red = 0
		}
		green := lineColor[1] - (gradientColor[1] * columnNum)
		if green < 0 {
			green = 0
		}
		blue := lineColor[2] - (gradientColor[2] * columnNum)
		if blue < 0 {
			blue = 0
		}
		finalColor := color.RGBA{uint8(red), uint8(green), uint8(blue), 255}
		ebitenutil.DrawLine(screen, float64(destX), float64(destY), float64(destX+width), float64(destY+height), finalColor)
	}
}

func DrawMinimap(screen *ebiten.Image, ecgView HealthECGView, xOffset int, yOffset int) {
	for startX := 0; startX < len(ecgView.Lines); startX++ {
		// Draw a vertical line for the current line
		destX := startX + xOffset + minimapOffset
		destY := ecgView.Lines[startX][0] + yOffset
		width := 1
		height := ecgView.Lines[startX][1] + 1

		lineColor := ecgView.Color
		finalColor := color.RGBA{uint8(lineColor[0]), uint8(lineColor[1]), uint8(lineColor[2]), 255}
		ebitenutil.DrawLine(screen, float64(destX), float64(destY), float64(destX+width), float64(destY+height), finalColor)
	}

	// Draw rectangular region over rendered area
	regionColor := color.RGBA{255, 255, 255, 255}
	regionUpperX := ecgOffsetX + xOffset + minimapOffset
	regionLowerY := yOffset
	regionLowerX := regionUpperX - 32
	if regionLowerX < 0 {
		regionLowerX = 0
	}
	regionUpperY := regionLowerY + 40

	ebitenutil.DrawLine(screen, float64(regionLowerX), float64(regionLowerY), float64(regionLowerX), float64(regionUpperY), regionColor)
	ebitenutil.DrawLine(screen, float64(regionUpperX), float64(regionLowerY), float64(regionUpperX), float64(regionUpperY), regionColor)
	ebitenutil.DrawLine(screen, float64(regionLowerX), float64(regionLowerY), float64(regionUpperX), float64(regionLowerY), regionColor)
	ebitenutil.DrawLine(screen, float64(regionLowerX), float64(regionUpperY), float64(regionUpperX), float64(regionUpperY), regionColor)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Health ECG")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
