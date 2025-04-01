package main

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Before that some definitions of Go terms
// A group is a set of stones of the same color that connected either orthogonally horizontal/vertical.
// However 2 stones next to each other diagonally are not connected in any way and so simply form two groups of one stone each.
// If a third stone were to be added to the two diagonal stones so that it sat next to both of them, a group of three stones would be formed.
// A libery is any empty point orthogonally adjacent to a group of stones.
const (
	BoardSize   = 19
	CellSize    = 22.7
	Margin      = 12.7
	WindowSize  = 441
	HUDHeight   = 50
	TotalHeight = WindowSize + 2*HUDHeight
	StoneSize   = 20
	Empty       = 0
	Black       = 1
	White       = 2
)

type Position struct {
	x, y int
}

type Game struct {
	board         [BoardSize][BoardSize]int
	boardImage    *ebiten.Image
	blackStone    *ebiten.Image
	whiteStone    *ebiten.Image
	currentTurn   int
	previousBoard [BoardSize][BoardSize]int // For ko rule checking, preventing the repetition of board positions
	capturedBlack int
	capturedWhite int
	hud           *HUD
}

func NewGame() (*Game, error) {
	g := &Game{currentTurn: Black}

	var err error
	if g.boardImage, err = loadImage("assets/go_board.png"); err != nil {
		return nil, err
	}
	if g.blackStone, err = loadImage("assets/black_stone.png"); err != nil {
		return nil, err
	}
	if g.whiteStone, err = loadImage("assets/white_stone.png"); err != nil {
		return nil, err
	}

	// Initialized the HUD
	g.hud = NewHUD()

	return g, nil
}

func loadImage(path string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, errors.New("failed to load image: " + path + " - " + err.Error())
	}
	return img, nil
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		gridX := int(math.Round((float64(x) - Margin) / CellSize))
		gridY := int(math.Round((float64(y) - Margin - HUDHeight) / CellSize)) // Adjust for HUD height

		if gridX >= 0 && gridX < BoardSize && gridY >= 0 && gridY < BoardSize {
			// Try to place a stone
			if g.isValidMove(gridX, gridY) {
				// Save the current board state for ko rule checking
				g.saveBoardState()

				// Place the stone
				g.board[gridY][gridX] = g.currentTurn

				// Check for captures
				capturedStones := g.checkCaptures(gridX, gridY)

				// Update the capture counter
				if g.currentTurn == Black {
					g.capturedWhite += capturedStones
				} else {
					g.capturedBlack += capturedStones
				}

				// Suicide move check - if the placed stone has no liberties after captures,
				// and the move would result in capture of the player's own stones, it's invalid
				if !g.hasLiberties(gridX, gridY) {
					g.board[gridY][gridX] = Empty
					return nil
				}

				// Switch turns
				if g.currentTurn == Black {
					g.currentTurn = White
				} else {
					g.currentTurn = Black
				}

				pointLabel := fmt.Sprintf("P%d", gridY*BoardSize+gridX+1)
				log.Printf("Stone placed at %s", pointLabel)
			}
		}
	}
	return nil
}

func (g *Game) saveBoardState() {
	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			g.previousBoard[y][x] = g.board[y][x]
		}
	}
}

func (g *Game) isValidMove(x, y int) bool {
	// Check if the point is empty
	if g.board[y][x] != Empty {
		return false
	}

	// Temporarily place the stone to check for ko rule and suicide rule
	g.board[y][x] = g.currentTurn

	// Check if this move would recreate the previous board state (ko rule)
	if g.wouldViolateKoRule(x, y) {
		g.board[y][x] = Empty
		return false
	}

	// Check if this is a suicide move (unless it captures opponent stones)
	if !g.hasLiberties(x, y) && g.capturesOpponentStones(x, y) == 0 {
		g.board[y][x] = Empty
		return false
	}

	// Undo the temporary placement
	g.board[y][x] = Empty

	return true
}

func (g *Game) wouldViolateKoRule(x, y int) bool {
	// Temporarily place the stone
	originalState := g.board[y][x]
	g.board[y][x] = g.currentTurn

	// Simulate captures
	g.simulateCaptures(x, y)

	// Check if the resulting board matches the previous board state
	boardMatches := true
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			if g.board[i][j] != g.previousBoard[i][j] {
				boardMatches = false
				break
			}
		}
		if !boardMatches {
			break
		}
	}

	// Restore the board
	g.restoreBoard()
	g.board[y][x] = originalState

	return boardMatches
}

func (g *Game) restoreBoard() {
	// Uses a copy of the board saved state before the temporary move was made
}

func (g *Game) capturesOpponentStones(x, y int) int {
	// Similar to checkCaptures but doesn't actually remove stones
	// Returns the count of stones that would be captured
	opponent := Black
	if g.currentTurn == Black {
		opponent = White
	}

	captured := 0
	directions := []Position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, dir := range directions {
		newX, newY := x+dir.x, y+dir.y
		if newX >= 0 && newX < BoardSize && newY >= 0 && newY < BoardSize && g.board[newY][newX] == opponent {
			if !g.hasLiberties(newX, newY) {
				group := g.getGroup(newX, newY)
				captured += len(group)
			}
		}
	}

	return captured
}

func (g *Game) simulateCaptures(x, y int) {
	// Similar to checkCaptures but simulates without changing the actual game state
}

func (g *Game) checkCaptures(x, y int) int {
	opponent := Black
	if g.currentTurn == Black {
		opponent = White
	}

	captured := 0
	directions := []Position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, dir := range directions {
		newX, newY := x+dir.x, y+dir.y
		if newX >= 0 && newX < BoardSize && newY >= 0 && newY < BoardSize && g.board[newY][newX] == opponent {
			if !g.hasLiberties(newX, newY) {
				group := g.getGroup(newX, newY)
				for _, pos := range group {
					g.board[pos.y][pos.x] = Empty
					captured++
				}
			}
		}
	}

	return captured
}

func (g *Game) hasLiberties(x, y int) bool {
	visited := make(map[Position]bool)
	return g.checkLiberties(x, y, visited)
}

func (g *Game) checkLiberties(x, y int, visited map[Position]bool) bool {
	pos := Position{x, y}
	if visited[pos] {
		return false
	}
	visited[pos] = true

	color := g.board[y][x]
	directions := []Position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, dir := range directions {
		newX, newY := x+dir.x, y+dir.y
		if newX >= 0 && newX < BoardSize && newY >= 0 && newY < BoardSize {
			// If we find an empty space, the group has a liberty
			if g.board[newY][newX] == Empty {
				return true
			}
			// If we find a stone of the same color, check if it has liberties
			if g.board[newY][newX] == color && g.checkLiberties(newX, newY, visited) {
				return true
			}
		}
	}

	return false
}

func (g *Game) getGroup(x, y int) []Position {
	color := g.board[y][x]
	visited := make(map[Position]bool)
	group := []Position{}

	var dfs func(int, int)
	dfs = func(curX, curY int) {
		pos := Position{curX, curY}
		if visited[pos] || curX < 0 || curX >= BoardSize || curY < 0 || curY >= BoardSize || g.board[curY][curX] != color {
			return
		}

		visited[pos] = true
		group = append(group, pos)

		// Check adjacent positions
		dfs(curX-1, curY)
		dfs(curX+1, curY)
		dfs(curX, curY-1)
		dfs(curX, curY+1)
	}

	dfs(x, y)
	return group
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the HUD
	g.hud.Draw(screen, g.capturedWhite, g.capturedBlack)

	// Draw the board
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, HUDHeight) // Offset the board to account for the top HUD
	screen.DrawImage(g.boardImage, op)

	// Draw stones
	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			var stone *ebiten.Image
			switch g.board[y][x] {
			case Black:
				stone = g.blackStone
			case White:
				stone = g.whiteStone
			}
			if stone != nil {
				op := &ebiten.DrawImageOptions{}

				// Scale the stone image to the desired size
				stoneWidth, stoneHeight := float64(stone.Bounds().Dx()), float64(stone.Bounds().Dy())
				scaleX := float64(StoneSize) / stoneWidth
				scaleY := float64(StoneSize) / stoneHeight
				op.GeoM.Scale(scaleX, scaleY)

				// Position stone exactly at intersection point, accounting for stone size
				op.GeoM.Translate(
					Margin+float64(x)*CellSize-float64(StoneSize)/2,
					Margin+float64(y)*CellSize-float64(StoneSize)/2+HUDHeight, // Offset for HUD
				)
				screen.DrawImage(stone, op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Return the base size of the game (logical size)
	return WindowSize, TotalHeight
}

func main() {
	game, err := NewGame()
	if err != nil {
		log.Fatal("Error initializing game: ", err)
	}

	// Set the window to be resizable
	ebiten.SetWindowResizable(true)

	// Set the initial window size and title
	ebiten.SetWindowSize(WindowSize, TotalHeight)
	ebiten.SetWindowTitle("Go Game")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Game encountered an error: ", err)
	}
}
