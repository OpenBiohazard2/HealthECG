package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
		DrawECGOverview(screen, healthECGViews[i], initialXOffset, initialYOffset+distanceBetweenViews*i)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// calculateGradientColor calculates the color for a line based on its position and gradient
func calculateGradientColor(baseColor [3]int, gradient [3]int, columnNum int) color.RGBA {
	red := baseColor[0] - (gradient[0] * columnNum)
	if red < 0 {
		red = 0
	}
	green := baseColor[1] - (gradient[1] * columnNum)
	if green < 0 {
		green = 0
	}
	blue := baseColor[2] - (gradient[2] * columnNum)
	if blue < 0 {
		blue = 0
	}
	return color.RGBA{uint8(red), uint8(green), uint8(blue), 255}
}

// drawVerticalLine draws a single vertical line on the screen
func drawVerticalLine(screen *ebiten.Image, x, y, width, height int, lineColor color.RGBA) {
	vector.StrokeLine(screen, float32(x), float32(y), float32(x+width), float32(y+height), 1, lineColor, false)
}

func DrawECGView(screen *ebiten.Image, ecgView HealthECGView, xOffset int, yOffset int) {
	for columnNum := 0; columnNum < 32; columnNum++ {
		startX := ecgOffsetX - columnNum
		if startX < 0 || startX >= len(ecgView.Lines) {
			continue
		}

		destX := startX + xOffset
		destY := ecgView.Lines[startX][0] + yOffset
		width := 1
		height := ecgView.Lines[startX][1] + 1

		lineColor := calculateGradientColor(ecgView.Color, ecgView.Gradient, columnNum)
		drawVerticalLine(screen, destX, destY, width, height, lineColor)
	}
}

// drawViewportIndicator draws the rectangular region overlay showing the current visible area
func drawViewportIndicator(screen *ebiten.Image, xOffset, yOffset int) {
	regionColor := color.RGBA{255, 255, 255, 255}
	regionUpperX := ecgOffsetX + xOffset + minimapOffset
	regionLowerY := yOffset
	regionLowerX := regionUpperX - 32
	if regionLowerX < 0 {
		regionLowerX = 0
	}
	regionUpperY := regionLowerY + 40

	// Draw the four sides of the rectangle
	vector.StrokeLine(screen, float32(regionLowerX), float32(regionLowerY), float32(regionLowerX), float32(regionUpperY), 1, regionColor, false)
	vector.StrokeLine(screen, float32(regionUpperX), float32(regionLowerY), float32(regionUpperX), float32(regionUpperY), 1, regionColor, false)
	vector.StrokeLine(screen, float32(regionLowerX), float32(regionLowerY), float32(regionUpperX), float32(regionLowerY), 1, regionColor, false)
	vector.StrokeLine(screen, float32(regionLowerX), float32(regionUpperY), float32(regionUpperX), float32(regionUpperY), 1, regionColor, false)
}

func DrawECGOverview(screen *ebiten.Image, ecgView HealthECGView, xOffset int, yOffset int) {
	// Draw all ECG lines in the overview
	for startX := 0; startX < len(ecgView.Lines); startX++ {
		destX := startX + xOffset + minimapOffset
		destY := ecgView.Lines[startX][0] + yOffset
		width := 1
		height := ecgView.Lines[startX][1] + 1

		lineColor := ecgView.Color
		finalColor := color.RGBA{uint8(lineColor[0]), uint8(lineColor[1]), uint8(lineColor[2]), 255}
		drawVerticalLine(screen, destX, destY, width, height, finalColor)
	}

	// Draw the viewport indicator
	drawViewportIndicator(screen, xOffset, yOffset)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Health ECG")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
