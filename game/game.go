package game

import (
	"io/ioutil"
	"log"
	"strings"
)

var MapSizeWidth int = 21
var MapSizeHeight int = 21

type Elem struct {
	Name   string          // Name of element
	Short  rune            // Short identifier for elems (to build maps)
	Solid  bool            // Solid can be walked over
	Images []ImagePosition // List of images for displaying
}

type ImagePosition struct {
	X1 int
	X2 int
	Y1 int
	Y2 int
}

type Player struct {
	X         float64
	Y         float64
	Speed     float64
	Direction rune // l or r
}

type Game struct {
	Ss          int           // Size of square elements
	width       int           // Number of blocks (width)
	height      int           // Number of blocks (heigth)
	AllElements map[rune]Elem // All elements
	GameMap     [][]rune      // Game map
	Player      Player
}

func InitGame() Game {
	// Init game structure
	game := Game{
		64,
		MapSizeWidth,
		MapSizeHeight,
		map[rune]Elem{},
		[][]rune{},
		Player{0.0, 0.0, 0.08, 'r'},
	}
	game.loadRessources()
	game.createMap()
	return game
}

func (game *Game) createMap() {
	for w := 0; w < game.width; w++ {
		game.GameMap = append(game.GameMap, []rune{})
		for h := 0; h < game.height; h++ {
			game.GameMap[w] = append(game.GameMap[w], ' ')
		}
	}

	file, err := ioutil.ReadFile("data/map.txt")
	if err != nil {
		log.Fatal(err)
	} else {
		lines := strings.Split(string(file), "\n")
		for x, l := range lines {
			for y, c := range l {
				if c != ' ' {
					game.GameMap[x][y] = c
				}
			}
		}
	}
}

func (game *Game) loadRessources() {
	game.loadRessource("stone", 's', true, []ImagePosition{{0, 1, 0, 1}})
	game.loadRessource("dirt", 'd', true, []ImagePosition{{1, 2, 0, 1}})
	game.loadRessource("grass", 'g', true, []ImagePosition{{2, 3, 0, 1}})
	game.loadRessource("brick", 'b', true, []ImagePosition{{3, 4, 0, 1}})
	game.loadRessource("herb_1", 'h', false, []ImagePosition{{0, 1, 1, 2}, {1, 2, 1, 2},
		{2, 3, 1, 2}, {3, 4, 1, 2}})
	game.loadRessource("coin_1", 'c', false, []ImagePosition{{0, 1, 2, 3}, {1, 2, 2, 3},
		{2, 3, 2, 3}, {3, 4, 2, 3}, {4, 5, 2, 3}, {5, 6, 2, 3}})
	game.loadRessource("player", 'p', false, []ImagePosition{{0, 1, 3, 4}, {1, 2, 3, 4},
		{2, 3, 3, 4}, {3, 4, 3, 4}, {4, 5, 3, 4}, {5, 6, 3, 4}, {5, 6, 3, 4}, {7, 8, 3, 4}})
}

func (game *Game) loadRessource(name string, short rune, solid bool, images []ImagePosition) {
	ip := []ImagePosition{}

	for _, v := range images {
		ip = append(ip, ImagePosition{
			game.Ss * v.X1,
			game.Ss * v.X2,
			game.Ss * v.Y1,
			game.Ss * v.Y2,
		})
	}

	game.AllElements[short] = Elem{
		name,
		short,
		solid,
		ip,
	}
}
