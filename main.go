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
}

func NewPlayer() Player {
	player := Player{}
	const playerWidth = screenWidth / 10
	const playerHeight = screenHeight / 50
	sprite := ebiten.NewImage(playerWidth, playerHeight)
	sprite.Fill(color.White)
	player.Sprite = sprite
	// start at the bottom left
	player.Y += float64(screenHeight - player.Height())
	return player
}

// bottomColliding checks if the bottom of entity is colliding with the top of
// other
func bottomColliding(entity, other Entity) bool {
	// Check that entity is within the x bounds of other
	if entity.X >= other.X && entity.X <= other.X+other.Width() {
		// Check that bottom of entity is touching top of other
		if entity.Y+entity.Height() >= other.Y {
			return true
		}
	}
	return false
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
		return
	} else if nextX > screenWidth-b.Width() {
		b.X = screenWidth - b.Width()
		b.XSpeed *= -1
		return
	} else {
		b.X = nextX
	}

	// mirror of X
	nextY := b.Y + b.YSpeed
	if nextY < 0 {
		b.Y = 0
		b.YSpeed *= -1
		return
	} else if nextY > screenHeight-b.Height() {
		b.Y = screenHeight - b.Height()
		b.YSpeed *= -1
		return
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
	Destroyed bool
}

func (b *Brick) Draw(screen *ebiten.Image) {
	if !b.Destroyed {
		b.Entity.Draw(screen)
	}
}

const brickWidth = screenWidth / 20
const brickHeight = screenHeight / 30

func NewBrick() Brick {
	brick := Brick{}
	sprite := ebiten.NewImage(brickWidth, brickHeight)
	sprite.Fill(color.White)
	brick.Sprite = sprite
	return brick
}

// Game implements ebiten.Game interface.
type Game struct {
	Player Player
	Ball   Ball
	Bricks []Brick
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.Player.Update()
	g.Ball.Update()

	// check for ball hitting player paddle
	if bottomColliding(g.Ball.Entity, g.Player.Entity) {
		g.Ball.Y = g.Player.Y - g.Ball.Height()
		g.Ball.YSpeed *= -1
		return nil
	}

	for i := range g.Bricks {
		brick := &g.Bricks[i]
		if brick.Destroyed {
			continue
		}
		if bottomColliding(brick.Entity, g.Ball.Entity) {
			brick.Destroyed = true
			g.Ball.Y = brick.Y + brick.Height()
			g.Ball.YSpeed *= -1
			return nil
		}
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	g.Player.Draw(screen)
	g.Ball.Draw(screen)
	for i := range g.Bricks {
		g.Bricks[i].Draw(screen)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	game := &Game{Player: NewPlayer(), Ball: NewBall()}
	// allocate enough space to fill half the screen with bricks
	numBricks := (screenWidth / brickWidth) * (screenHeight / brickHeight / 2)
	bricks := make([]Brick, 0, numBricks)
	for i := 0; i < numBricks; i++ {
		bricks = append(bricks, NewBrick())
	}
	game.Bricks = bricks
	var currX float64
	var currY float64
	for i := range game.Bricks {
		b := &bricks[i]
		b.X = currX
		b.Y = currY

		currX += brickWidth
		if currX >= screenWidth {
			currX = 0
			currY += brickHeight
		}
	}
	return game
}

func slowMode() {
	ebiten.SetMaxTPS(30)
}

func main() {
	// Debugging
	if screenWidth%brickWidth != 0 {
		log.Printf("bricks will not tile horizontally because the screen width is not divisible by the brick width")
	}

	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Brick Break")
	game := NewGame()
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
