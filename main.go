package main

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	boardSize = 19
	cellSize  = 30
)

var (
	boardImage *ebiten.Image
	woodImage  *ebiten.Image
	blackStone *ebiten.Image
	whiteStone *ebiten.Image

	board         [boardSize][boardSize]int
	currentPlayer int = 1
)

type Game struct{}

func init() {
	var err error
	boardImage, err = loadImage("assets/board.jpg")
	if err != nil {
		log.Fatal(err)
	}
	woodImage, err = loadImage("assets/wood_full_original.jpg")
	if err != nil {
		log.Fatal(err)
	}
	blackStone, err = loadImage("assets/black.png")
	if err != nil {
		log.Fatal(err)
	}
	whiteStone, err = loadImage("assets/white.png")
	if err != nil {
		log.Fatal(err)
	}
}

func loadImage(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		boardX := x / cellSize
		boardY := y / cellSize

		if boardX >= 0 && boardX < boardSize && boardY >= 0 && boardY < boardSize {
			if board[boardX][boardY] == 0 {
				board[boardX][boardY] = currentPlayer
				if currentPlayer == 1 {
					currentPlayer = 2
				} else {
					currentPlayer = 1
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(woodImage, nil)
	screen.DrawImage(boardImage, nil)

	for x := 0; x < boardSize; x++ {
		for y := 0; y < boardSize; y++ {
			if board[x][y] == 1 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*cellSize), float64(y*cellSize))
				screen.DrawImage(blackStone, op)
			} else if board[x][y] == 2 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*cellSize), float64(y*cellSize))
				screen.DrawImage(whiteStone, op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Go Game")
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
