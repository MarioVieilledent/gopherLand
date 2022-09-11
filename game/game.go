package game

import (
	"io/ioutil"
	"log"
	"strings"
)

type Elem struct {
	Name   string // Name of element
	Short  rune   // Short identifier for elems (to build maps)
	Images []ImagePosition
}

type ImagePosition struct {
	X1 int
	X2 int
	Y1 int
	Y2 int
}

type Game struct {
	Ss          int           // Size of square elements
	width       int           // Number of blocks (width)
	height      int           // Number of blocks (heigth)
	AllElements map[rune]Elem // All elements
	GameMap     [][]rune      // Game map
}

func InitGame() Game {
	// Init game structure
	game := Game{
		64,
		21,
		21,
		map[rune]Elem{},
		[][]rune{},
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
	game.loadRessource("stone", 's', []ImagePosition{{0, 1, 0, 1}})
	game.loadRessource("dirt", 'd', []ImagePosition{{1, 2, 0, 1}})
	game.loadRessource("grass", 'g', []ImagePosition{{2, 3, 0, 1}})
	game.loadRessource("brick", 'b', []ImagePosition{{3, 4, 0, 1}})

	game.loadRessource("herb_1", 'h', []ImagePosition{{0, 1, 1, 2}, {1, 2, 1, 2},
		{2, 3, 1, 2}, {3, 4, 1, 2}})

	game.loadRessource("coin_1", 'c', []ImagePosition{{0, 1, 2, 3}, {1, 2, 2, 3},
		{2, 3, 2, 3}, {3, 4, 2, 3}, {4, 5, 2, 3}, {5, 6, 2, 3}})
}

func (game *Game) loadRessource(name string, short rune, images []ImagePosition) {
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
		ip,
	}
}
