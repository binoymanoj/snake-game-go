package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	directionUp     = Point{x: 0, y: -1}
	directionLeft   = Point{x: -1, y: 0}
	directionDown   = Point{x: 0, y: 1}
	directionRight  = Point{x: 1, y: 0}
	mplusFaceSource *text.GoTextFaceSource
)

const (
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
	gameSpeed    = time.Second / 6
)

type Point struct {
	x, y int // co-ordinates
}

type Game struct {
	snake      []Point
	direction  Point
	lastUpdate time.Time
	food       Point
	gameOver   bool
}

func (g *Game) Update() error {
	if g.gameOver {
		// TODO: implment score card with name and restart game with Enter/Space key
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.direction = directionUp
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.direction = directionLeft
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.direction = directionDown
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.direction = directionRight
	}

	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}

	g.lastUpdate = time.Now()

	g.updateSnake(&g.snake, g.direction)
	return nil
}

func (g *Game) updateSnake(snake *[]Point, direction Point) {
	head := (*snake)[0]

	newHead := Point{x: head.x + direction.x, y: head.y + direction.y}

	if g.isCollision(newHead, *snake) {
		g.gameOver = true
		return
	}

	if newHead == g.food {
		*snake = append([]Point{newHead}, *snake...)
		g.spawnFood()
	} else {
		*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)
	}

}

func (g Game) isCollision(p Point, snake []Point) bool {
	if p.x < 0 || p.y < 0 || p.x >= screenWidth/gridSize || p.y >= screenHeight/gridSize {
		return true
	}

	for _, snakepoint := range snake {
		if snakepoint == p {
			return true
		}
	}

	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.snake {
		vector.DrawFilledRect(screen, float32(p.x*gridSize), float32(p.y*gridSize), gridSize, gridSize, color.White, true)
	}

	vector.DrawFilledRect(screen, float32(g.food.x*gridSize), float32(g.food.y*gridSize), gridSize, gridSize, color.RGBA{255, 0, 0, 255}, true)

	if g.gameOver {
		face := &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   48,
		}

		t := "Game Over!"
		w, h := text.Measure(t, face, face.Size)

		opt := &text.DrawOptions{}
		opt.GeoM.Translate(screenWidth/2-w/2, screenHeight/2-h/2)
		opt.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, t, face, opt)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) spawnFood() {
	g.food = Point{rand.Intn(screenWidth / gridSize), rand.Intn(screenHeight / gridSize)}
}

func main() {
	s, err := text.NewGoTextFaceSource(
		bytes.NewReader(
			fonts.MPlus1pRegular_ttf,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	game := &Game{
		snake:     []Point{{x: screenWidth / gridSize / 2, y: screenHeight / gridSize / 2}},
		direction: Point{x: 1, y: 0},
	}

	game.spawnFood()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Game Go")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
