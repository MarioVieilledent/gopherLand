package graphic

import (
	"gopherLand/game"
	"image"
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/tinne26/etxt"
)

const xPlayerFixed int = 10

type Controller struct {
	game        *game.Game
	tick        uint64         // Ticks of the game
	tickFrame   uint8          // Increments each frame, go back to 0 each tick
	txtRenderer *etxt.Renderer // Used to render text on screen
}

var backgroundImage *ebiten.Image
var background3Image *ebiten.Image
var resourcesImage *ebiten.Image
var iconImage *ebiten.Image

var playerShift float64 // Shift for displaying player

//////////////////////////////
// INITIALIZATION FUNCTIONS //
//////////////////////////////

func initController() Controller {
	g := game.InitGame(xPlayerFixed)
	playerShift = 0.5 * float64(g.Ss)
	return Controller{&g, 0, 0, getTxtRenderer()}
}

func init() {
	var err error
	resourcesImage, _, err = ebitenutil.NewImageFromFile("data/images/resources/resources.png")
	if err != nil {
		log.Fatal(err)
	}
	backgroundImage, _, err = ebitenutil.NewImageFromFile("data/images/backgrounds/background.png")
	if err != nil {
		log.Fatal(err)
	}
	background3Image, _, err = ebitenutil.NewImageFromFile("data/images/backgrounds/background3.png")
	if err != nil {
		log.Fatal(err)
	}
	iconImage, _, err = ebitenutil.NewImageFromFile("data/images/icons/icon.png")
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

// Manages tick for modulo for animate blocks
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
	// c.game.Player.TouchingGround = false
	// Forces the player to stop walking while on the air
	/*
		if c.game.Player.Walking {
			c.game.Player.Walking = false
		}
	*/
	// Move vertically the player depending on its vertical velocity
	if !c.game.TouchesGround() {
		c.game.Move(0.0, (0.01 * c.game.Player.VerticalVelocity))
		c.game.Player.VerticalVelocity += 1.0
	}
}

// Manages input for controlling player
func (c *Controller) manageButtonClicks() {
	var notClickedLeft, notClickedRight bool

	// Left button clicked
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

	// Down button clicked
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		c.game.GoDown()
	}
}

///////////////////////
// DRAWING ON WINDOW //
///////////////////////

func (c *Controller) Draw(screen *ebiten.Image) {
	c.displayBackgrounds(screen)
	c.displayBlocks(screen)
	c.displayPlayer(screen)
}

// Draw backgrounds
func (c *Controller) displayBackgrounds(screen *ebiten.Image) {
	// Fixed background image
	screen.DrawImage(backgroundImage, nil)

	// Moving background image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-c.game.Player.Position.X*0.6*float64(c.game.Ss), 0)
	screen.DrawImage(background3Image, op)
}

// Draw all blocks of the map
func (c *Controller) displayBlocks(screen *ebiten.Image) {
	for x := 0; x < len(c.game.GameMap); x++ {
		for y := 0; y < len(c.game.GameMap[x]); y++ {

			block := c.game.AllBlocks[c.game.GameMap[x][y]]

			modulo := c.getModulo(block)

			if len(block.Images) > 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(c.game.Ss*x)-
					c.game.Player.Position.X*float64(c.game.Ss)+
					float64(xPlayerFixed*c.game.Ss),
					float64(c.game.Ss*y))
				screen.DrawImage(resourcesImage.SubImage(
					image.Rect(block.Images[modulo].X1, block.Images[modulo].Y1,
						block.Images[modulo].X2, block.Images[modulo].Y2)).(*ebiten.Image), op)
			}
		}
	}
}

// Draw the player
func (c *Controller) displayPlayer(screen *ebiten.Image) {
	modulo := c.getModulo(c.game.AllBlocks['p'])

	op := &ebiten.DrawImageOptions{}
	if c.game.Player.Direction == 'l' {
		// Mirroring image
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(c.game.Ss)-playerShift, -playerShift)
	} else {
		op.GeoM.Translate(-playerShift, -playerShift)
	}
	op.GeoM.Translate(float64(xPlayerFixed*c.game.Ss),
		c.game.Player.Position.Y*float64(c.game.Ss))
	if c.game.Player.Walking {
		screen.DrawImage(resourcesImage.SubImage(
			image.Rect(c.game.AllBlocks['p'].Images[modulo].X1,
				c.game.AllBlocks['p'].Images[modulo].Y1,
				c.game.AllBlocks['p'].Images[modulo].X2,
				c.game.AllBlocks['p'].Images[modulo].Y2)).(*ebiten.Image), op)
	} else {
		screen.DrawImage(resourcesImage.SubImage(
			image.Rect(c.game.AllBlocks['p'].Images[1].X1,
				c.game.AllBlocks['p'].Images[1].Y1,
				c.game.AllBlocks['p'].Images[1].X2,
				c.game.AllBlocks['p'].Images[1].Y2)).(*ebiten.Image), op)
	}

	// Config textRenderer
	c.txtRenderer.SetTarget(screen)
	c.txtRenderer.SetSizePx(42)

	// Display number of golds
	c.txtRenderer.SetColor(color.RGBA{188, 94, 16, 255})
	c.txtRenderer.Draw("Golds: "+strconv.Itoa(c.game.Player.Gold), 20, 10)

	// Display number of golds
	c.txtRenderer.SetColor(color.RGBA{147, 31, 124, 255})
	c.txtRenderer.Draw("Keys: "+strconv.Itoa(c.game.Player.Keys), 20, 40)
}

/////////////////////
// OTHER FUNCTIONS //
/////////////////////

// Used to make animation of blocks
func (c Controller) getModulo(block game.Block) int {
	max := len(block.Images)
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
	ebiten.SetWindowTitle("GopherLand")
	ebiten.SetWindowIcon([]image.Image{iconImage})

	controler := initController()

	err := ebiten.RunGame(&controler)

	if err != nil {
		log.Fatal(err)
	}
}
