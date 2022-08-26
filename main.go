package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Position struct {
	X float64
	Y float64
}

type Entity struct {
	Position
	Sprite *ebiten.Image
}

func (e *Entity) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(e.X, e.Y)
	screen.DrawImage(e.Sprite, options)
}

func (e *Entity) Width() float64 {
	result := e.Sprite.Bounds().Dx()
	return float64(result)
}

func (e *Entity) Height() float64 {
	result := e.Sprite.Bounds().Dy()
	return float64(result)
}

type Player struct {
	Entity
}

func playerSprite() *ebiten.Image {
	const playerWidth = screenWidth / 15
	const playerHeight = screenHeight / 30
	sprite := ebiten.NewImage(playerWidth, playerHeight)
	sprite.Fill(color.White)
	return sprite
}

func NewPlayer() Player {
	player := Player{}
	player.Sprite = playerSprite()
	// start at the bottom left
	player.Y += float64(screenHeight - player.Height())
	return player
}

func (p *Player) Update() error {
	const playerVelocity = screenWidth / 100
	leftPressed := ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	rightPressed := ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	if leftPressed && rightPressed {
		return nil
	} else if leftPressed {
		nextX := p.X - playerVelocity
		if nextX >= 0 {
			p.X = nextX
		} else {
			// Allows the player to lock to the side of the screen
			p.X = 0
		}
	} else if rightPressed {
		// mirror of leftPressed
		nextX := p.X + playerVelocity
		if nextX <= screenWidth-p.Width() {
			p.X = nextX
		} else {
			p.X = screenWidth - p.Width()
		}
	}
	return nil
}

// Game implements ebiten.Game interface.
type Game struct {
	Player Player
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	err := g.Player.Update()
	if err != nil {
		return err
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	g.Player.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	game := &Game{Player: NewPlayer()}
	return game
}

func main() {
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Brick Break")
	game := NewGame()
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
