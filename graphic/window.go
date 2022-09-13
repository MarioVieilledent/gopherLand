package graphic

import (
	"gopherLand/game"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Controller struct {
	game      *game.Game
	tick      uint64 // Ticks of the game
	tickFrame uint8  // Increments each frame, go back to 0 each tick
}

var backgroundImage *ebiten.Image
var resourcesImage *ebiten.Image

var playerShift float64 // Shift for displaying player

//////////////////////////////
// INITIALIZATION FUNCTIONS //
//////////////////////////////

func initController() Controller {
	g := game.InitGame()
	playerShift = 0.5 * float64(g.Ss)
	return Controller{&g, 0, 0}
}

func init() {
	var err error
	resourcesImage, _, err = ebitenutil.NewImageFromFile("data/ressources.png")
	if err != nil {
		log.Fatal(err)
	}
	backgroundImage, _, err = ebitenutil.NewImageFromFile("data/background.png")
	if err != nil {
		log.Fatal(err)
	}
}

//////////////////////
// UPDATE FUNCTIONS //
//////////////////////

// Update function, called each frame
func (c *Controller) Update() error {
	// Tick management (each 12 frames = 200 ms)
	c.manageTick()

	// Manages jumping
	c.manageJumpOrFall()

	// Manages button clicks
	c.manageButtonClicks()

	return nil
}

// Manages tick for modulo for animate elements
func (c *Controller) manageTick() {
	if c.tickFrame == 6 {
		c.tickFrame = 0
		c.tick++
	} else {
		c.tickFrame++
	}
}

// Manages jump and fall of player
func (c *Controller) manageJumpOrFall() {
	touchingGround := c.game.CheckIfTouchesGround()
	if !touchingGround {
		// Forces the player to stop walking while on the air
		if c.game.Player.Walking {
			c.game.Player.Walking = false
		}
		// Move vertically the player depending on its vertical velocity
		c.game.Move(0.0, (0.01 * c.game.Player.VerticalVelocity))
		c.game.Player.VerticalVelocity += 1.0
		// If player hit ceil while jumping
	} else if c.game.Player.VerticalVelocity != 0 {
		c.game.Player.VerticalVelocity = 0
	}
}

// Manages input for controlling player
func (c *Controller) manageButtonClicks() {
	var notClickedLeft, notClickedRight bool

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		c.game.Player.Direction = 'l'
		c.game.Player.Walking = c.game.Move(-c.game.Player.Speed, 0)
	} else {
		notClickedLeft = true
	}

	// Right button clicked
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		c.game.Player.Direction = 'r'
		c.game.Player.Walking = c.game.Move(c.game.Player.Speed, 0)
	} else {
		notClickedRight = true
	}

	if notClickedLeft && notClickedRight && c.game.Player.Walking {
		c.game.Player.Walking = false
	}

	// Up button clicked
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if c.game.Player.TouchingGround {
			c.game.Player.TouchingGround = false
			c.game.Player.VerticalVelocity = -20.0
			c.game.Move(0.0, -0.01)
		}
	}
}

///////////////////////
// DRAWING ON WINDOW //
///////////////////////

func (c *Controller) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, nil)
	c.displayBlocks(screen)
	c.displayPlayer(screen)
}

// Draw all blocks of the map
func (c *Controller) displayBlocks(screen *ebiten.Image) {
	for x := 0; x < len(c.game.GameMap); x++ {
		for y := 0; y < len(c.game.GameMap[x]); y++ {

			elem := c.game.AllElements[c.game.GameMap[x][y]]

			modulo := c.getModulo(elem)

			if len(elem.Images) > 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(c.game.Ss*x), float64(c.game.Ss*y))
				screen.DrawImage(resourcesImage.SubImage(
					image.Rect(elem.Images[modulo].X1, elem.Images[modulo].Y1,
						elem.Images[modulo].X2, elem.Images[modulo].Y2)).(*ebiten.Image), op)
			}
		}
	}
}

// Draw the player
func (c *Controller) displayPlayer(screen *ebiten.Image) {
	modulo := c.getModulo(c.game.AllElements['p'])

	op := &ebiten.DrawImageOptions{}
	if c.game.Player.Direction == 'l' {
		// Mirroring image
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(c.game.Ss)-playerShift, -playerShift)
	} else {
		op.GeoM.Translate(-playerShift, -playerShift)
	}
	op.GeoM.Translate(c.game.Player.Position.X*float64(c.game.Ss),
		c.game.Player.Position.Y*float64(c.game.Ss))
	if c.game.Player.Walking {
		screen.DrawImage(resourcesImage.SubImage(
			image.Rect(c.game.AllElements['p'].Images[modulo].X1,
				c.game.AllElements['p'].Images[modulo].Y1,
				c.game.AllElements['p'].Images[modulo].X2,
				c.game.AllElements['p'].Images[modulo].Y2)).(*ebiten.Image), op)
	} else {
		screen.DrawImage(resourcesImage.SubImage(
			image.Rect(c.game.AllElements['p'].Images[1].X1,
				c.game.AllElements['p'].Images[1].Y1,
				c.game.AllElements['p'].Images[1].X2,
				c.game.AllElements['p'].Images[1].Y2)).(*ebiten.Image), op)
	}
}

/////////////////////
// OTHER FUNCTIONS //
/////////////////////

// Used to make animation of elements
func (c Controller) getModulo(elem game.Elem) int {
	max := len(elem.Images)
	modulo := 0
	if max > 1 {
		modulo = int(c.tick) % max
	}
	return modulo
}

func (c *Controller) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func OpenWindow() {
	ebiten.SetWindowSize(1280, 720)

	controler := initController()

	err := ebiten.RunGame(&controler)

	if err != nil {
		log.Fatal(err)
	}
}
