package game

import (
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Elem struct {
	Name  string // Name of element
	Short string // Short identifier for elems (to build maps)
	// Positions in the image containing all the graphics
	X1 int
	X2 int
	Y1 int
	Y2 int
}

type Game struct {
	Tick        uint64          // Ticks (one tick each 100 ms)
	Ss          int             // Size of square elements
	width       int             // Number of blocks (width)
	height      int             // Number of blocks (heigth)
	AllElements map[string]Elem // All elements
	GameMap     [][]Elem        // Game map
}

func InitGame() Game {
	// Init game structure
	game := Game{
		0,
		64,
		21,
		21,
		map[string]Elem{},
		[][]Elem{},
	}
	game.loadRessources()
	game.createMap()
	// starts the gameLoop
	// go game.startGameLoop()
	return game
}

func (game *Game) startGameLoop() {
	for {
		time.Sleep(100 * time.Millisecond)
		game.Tick++
	}
}

func (game *Game) createMap() {
	for w := 0; w < game.width; w++ {
		game.GameMap = append(game.GameMap, []Elem{})
		for h := 0; h < game.height; h++ {
			game.GameMap[w] = append(game.GameMap[w], Elem{})
		}
	}

	file, err := ioutil.ReadFile("data/map.txt")
	if err != nil {
		log.Fatal(err)
	} else {
		lines := strings.Split(string(file), "\n")
		for x, l := range lines {
			for y, c := range l {
				if c == ' ' {

				} else {
					game.GameMap[x][y] = game.AllElements[string(c)]
				}
			}
		}
	}
}

func (game *Game) loadRessources() {
	game.loadRessource("stone", "s", 0, 1, 0, 1)
	game.loadRessource("dirt", "d", 1, 2, 0, 1)
	game.loadRessource("grass", "g", 2, 3, 0, 1)
	game.loadRessource("brick", "b", 3, 4, 0, 1)

	game.loadRessource("herb_1", "h", 0, 1, 1, 2)
	// game.loadRessource("herb_2", "h", 1, 2, 1, 2)
	// game.loadRessource("herb_3", "h", 2, 3, 1, 2)
	// game.loadRessource("herb_4", "h", 3, 4, 1, 2)

	game.loadRessource("coin_1", "c", 0, 1, 2, 3)
	// game.loadRessource("coin_2", "c", 1, 2, 2, 3)
	// game.loadRessource("coin_3", "c", 2, 3, 2, 3)
	// game.loadRessource("coin_4", "c", 3, 4, 2, 3)
	// game.loadRessource("coin_5", "c", 4, 5, 2, 3)
	// game.loadRessource("coin_6", "c", 5, 6, 2, 3)
}

func (game *Game) loadRessource(name string, short string, x1, x2, y1, y2 int) {
	elem := Elem{
		name,
		short,
		game.Ss * x1,
		game.Ss * x2,
		game.Ss * y1,
		game.Ss * y2,
	}
	game.AllElements[short] = elem
}
