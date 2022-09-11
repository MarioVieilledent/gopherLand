package graphic

import (
	"gopherLand/game"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Controller struct {
	game      *game.Game
	tick      uint64 // Ticks of the game
	tickFrame uint8  // Increments each frame, go back to 0 each tick
}

var backgroundImage *ebiten.Image
var ressourcesImage *ebiten.Image

func initControler() Controller {
	g := game.InitGame()
	return Controller{&g, 0, 0}
}

func init() {
	var err error
	ressourcesImage, _, err = ebitenutil.NewImageFromFile("data/ressources.png")
	if err != nil {
		log.Fatal(err)
	}
	backgroundImage, _, err = ebitenutil.NewImageFromFile("data/background.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (c Controller) getModulo(elem game.Elem) int {
	max := len(elem.Images)
	modulo := 0
	if max > 1 {
		modulo = int(c.tick) % max
	}
	return modulo
}

func (c Controller) checkEmptyBlock(direction string) bool {
	switch direction {
	case "left":
		x := int(c.game.Player.X)
		y := int(math.Ceil(c.game.Player.Y))
		if y > 0 {
			return !c.game.AllElements[c.game.GameMap[x][y-1]].Solid
		}
	case "right":
		x := int(c.game.Player.X)
		y := int(c.game.Player.Y)
		if y < game.MapSizeWidth {
			return !c.game.AllElements[c.game.GameMap[x][y+1]].Solid
		}
	case "up":
		x := int(math.Ceil(c.game.Player.X))
		y := int(c.game.Player.Y)
		if x > 0 {
			return !c.game.AllElements[c.game.GameMap[x-1][y]].Solid
		}
	case "down":
		x := int(c.game.Player.X)
		y := int(c.game.Player.Y)
		if x < game.MapSizeWidth {
			return !c.game.AllElements[c.game.GameMap[x+1][y]].Solid
		}
	}
	return false
}

func (c *Controller) Update() error {
	// Tick management (each 12 frames = 200 ms)
	if c.tickFrame == 12 {
		c.tickFrame = 0
		c.tick++
	} else {
		c.tickFrame++
	}

	// Control management
	// x, y := ebiten.CursorPosition()
	// fmt.Println(x, y)
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		c.game.Player.Direction = 'l'
		if c.checkEmptyBlock("left") {
			c.game.Player.Y -= c.game.Player.Speed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		c.game.Player.Direction = 'r'
		if c.checkEmptyBlock("right") {
			c.game.Player.Y += c.game.Player.Speed
		}
	}

	// Player gravity
	if c.checkEmptyBlock("down") {
		c.game.Player.X += 0.1
	}

	return nil
}

func (c *Controller) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, nil)

	// For each block (Elem)
	for y := 0; y < len(c.game.GameMap); y++ {
		for x := 0; x < len(c.game.GameMap[y]); x++ {

			elem := c.game.AllElements[c.game.GameMap[x][y]]

			modulo := c.getModulo(elem)

			if len(elem.Images) > 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(c.game.Ss*y), float64(c.game.Ss*x))
				screen.DrawImage(ressourcesImage.SubImage(
					image.Rect(elem.Images[modulo].X1, elem.Images[modulo].Y1,
						elem.Images[modulo].X2, elem.Images[modulo].Y2)).(*ebiten.Image), op)
			}
		}
	}

	// Player display
	modulo := c.getModulo(c.game.AllElements['p'])

	op := &ebiten.DrawImageOptions{}
	if c.game.Player.Direction == 'l' {
		// Mirroring image
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(c.game.Ss), 0)
	}
	op.GeoM.Translate(c.game.Player.Y*float64(c.game.Ss), c.game.Player.X*float64(c.game.Ss))
	screen.DrawImage(ressourcesImage.SubImage(
		image.Rect(c.game.AllElements['p'].Images[modulo].X1, c.game.AllElements['p'].Images[modulo].Y1,
			c.game.AllElements['p'].Images[modulo].X2, c.game.AllElements['p'].Images[modulo].Y2)).(*ebiten.Image), op)
}

func (c *Controller) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func OpenWindow() {
	ebiten.SetWindowSize(1280, 720)

	controler := initControler()

	err := ebiten.RunGame(&controler)

	if err != nil {
		log.Fatal(err)
	}
}
