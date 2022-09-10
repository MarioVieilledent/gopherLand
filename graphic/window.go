package graphic

import (
	"gopherLand/game"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Controler struct {
	game *game.Game
}

var backgroundImage *ebiten.Image
var ressourcesImage *ebiten.Image

func initControler() Controler {
	g := game.InitGame()
	return Controler{&g}
}

func init() {
	var err error
	ressourcesImage, _, err = ebitenutil.NewImageFromFile("data/ressources.png")
	backgroundImage, _, err = ebitenutil.NewImageFromFile("data/background.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Controler) Update() error {
	return nil
}

func (c *Controler) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, nil)

	for x, line := range c.game.GameMap {
		for y, elem := range line {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(c.game.Ss*y), float64(c.game.Ss*x))
			screen.DrawImage(ressourcesImage.SubImage(
				image.Rect(elem.X1, elem.Y1, elem.X2, elem.Y2)).(*ebiten.Image), op)
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
