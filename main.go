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

func (p *Player) Update() {
	const playerSpeed = screenWidth / 100
	leftPressed := ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	rightPressed := ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	if leftPressed && rightPressed {
		return
	} else if leftPressed {
		nextX := p.X - playerSpeed
		if nextX >= 0 {
			p.X = nextX
		} else {
			// Allows the player to lock to the side of the screen
			p.X = 0
		}
	} else if rightPressed {
		// mirror of leftPressed
		nextX := p.X + playerSpeed
		if nextX <= screenWidth-p.Width() {
			p.X = nextX
		} else {
			p.X = screenWidth - p.Width()
		}
	}
	return
}

func NewPlayer() Player {
	player := Player{}
	const playerWidth = screenWidth / 15
	const playerHeight = screenHeight / 30
	sprite := ebiten.NewImage(playerWidth, playerHeight)
	sprite.Fill(color.White)
	player.Sprite = sprite
	// start at the bottom left
	player.Y += float64(screenHeight - player.Height())
	return player
}

type Ball struct {
	XSpeed float64
	YSpeed float64
	Entity
}

func (b *Ball) Update() {
	nextX := b.X + b.XSpeed
	// Check bounds
	if nextX < 0 {
		b.X = 0
		b.XSpeed *= -1
	} else if nextX > screenWidth-b.Width() {
		b.X = screenWidth - b.Width()
		b.XSpeed *= -1
	} else {
		b.X = nextX
	}

	// mirror of X
	nextY := b.Y + b.YSpeed
	if nextY < 0 {
		b.Y = 0
		b.YSpeed *= -1
	} else if nextY > screenHeight-b.Height() {
		b.Y = screenHeight - b.Height()
		b.YSpeed *= -1
	} else {
		b.Y = nextY
	}
}

func NewBall() Ball {
	const ballXSpeed = screenWidth / 150
	const ballYSpeed = screenWidth / 150
	ball := Ball{XSpeed: ballXSpeed, YSpeed: ballYSpeed}
	const ballWidth = screenWidth / 40
	const ballHeight = screenHeight / 40
	sprite := ebiten.NewImage(ballWidth, ballHeight)
	sprite.Fill(color.White)
	ball.Sprite = sprite
	ball.X += screenWidth / 2
	ball.Y += screenHeight / 2
	return ball
}

type Brick struct {
	Entity
}

func NewBrick() Brick {
	brick := Brick{}
	const brickWidth = screenWidth / 20
	const brickHeight = screenHeight / 30
	if screenWidth%brickWidth != 0 {
		log.Printf("bricks will not tile horizontally because the screen width is not divisible by the brick width")
	}
	sprite := ebiten.NewImage(brickWidth, brickHeight)
	sprite.Fill(color.White)
	brick.Sprite = sprite
	return Brick{}
}

func (b *Brick) Update() {

}

// Game implements ebiten.Game interface.
type Game struct {
	Player Player
	Ball   Ball
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.Player.Update()
	g.Ball.Update()
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	g.Player.Draw(screen)
	g.Ball.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	game := &Game{Player: NewPlayer(), Ball: NewBall()}
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
