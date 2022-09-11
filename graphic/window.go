package graphic

import (
	"fmt"
	"gopherLand/game"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Controler struct {
	game      *game.Game
	tick      uint64 // Ticks of the game
	tickFrame uint8  // Increments each frame, go back to 0 each tick
}

var backgroundImage *ebiten.Image
var ressourcesImage *ebiten.Image

func initControler() Controler {
	g := game.InitGame()
	return Controler{&g, 0, 0}
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

func (c *Controler) Update() error {
	if c.tickFrame == 12 {
		c.tickFrame = 0
		c.tick++
	} else {
		c.tickFrame++
	}
	// x, y := ebiten.CursorPosition()
	// fmt.Println(x, y)
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		fmt.Println("Droite")
	}
	return nil
}

func (c *Controler) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, nil)

	for y := 0; y < len(c.game.GameMap); y++ {
		for x := 0; x < len(c.game.GameMap[y]); x++ {

			elem := c.game.AllElements[c.game.GameMap[x][y]]

			max := len(elem.Images)
			modulo := 0
			if max > 1 {
				modulo = int(c.tick) % max
			}

			if len(elem.Images) > 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(c.game.Ss*y), float64(c.game.Ss*x))
				screen.DrawImage(ressourcesImage.SubImage(
					image.Rect(elem.Images[modulo].X1, elem.Images[modulo].Y1,
						elem.Images[modulo].X2, elem.Images[modulo].Y2)).(*ebiten.Image), op)
			}
		}
	}
}

func (c *Controler) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
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
