package game

import (
	"io/ioutil"
	"log"
	"strings"
)

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

type Position struct {
	X float64
	Y float64
}

type Player struct {
	Position         Position
	EatBox           [4]float64 // Distance from center (position) (up, right, down, right)
	Speed            float64
	Direction        rune    // l or r
	TouchingGround   bool    // True if player is walking, false if falling or jumping
	VerticalVelocity float64 // Vertical velocity on air (gravity falling, or gravity jumping)
	Walking          bool    // Animate player when walking
}

type Game struct {
	Ss          int           // Size of square elements
	width       int           // Number of blocks (width)
	height      int           // Number of blocks (height)
	AllElements map[rune]Elem // All elements
	GameMap     [][]rune      // Game map
	Player      Player
}

func InitGame() Game {
	// Init game structure
	game := Game{
		64,
		0,
		0,
		map[rune]Elem{},
		[][]rune{},
		Player{
			Position{1.5, 1.5},
			[4]float64{0.45, 0.4, 0.5, 0.4},
			0.08,
			'r',
			false,
			0.0,
			false,
		},
	}
	game.loadResources()
	game.createMap()
	return game
}

///////////////////
// GAME METHODS ///
///////////////////

// Checks if player is touching the ground (not falling or jumping)
func (g *Game) CheckIfTouchesGround() bool {
	x := int(g.Player.Position.X)
	y := int(g.Player.Position.Y + g.Player.EatBox[2])
	if y < g.height-1 {
		if g.AllElements[g.GameMap[x][y]].Solid {
			g.Player.TouchingGround = true
			return true
		} else {
			g.Player.TouchingGround = false
		}
	}
	return false
}

// Check if player touches a block while jumping
func (g *Game) CheckIfTouchesCeil() bool {
	x := int(g.Player.Position.X)
	y := int(g.Player.Position.Y - g.Player.EatBox[0])
	if y < 0 {
		return true
	}
	return g.AllElements[g.GameMap[x][y]].Solid
}

func (game *Game) createMap() {
	// Reads the resources file
	file, err := ioutil.ReadFile("data/map.txt")
	if err != nil {
		log.Fatal(err) // /!\ Need a better handle of the error here /!\
	} else {
		// Takes all lines as a slice of strings
		lines := strings.Split(string(file), "\n")

		// Get the size of the map
		game.height = len(lines)
		game.width = longestStr(lines)

		// Generate empty map
		for w := 0; w < game.width; w++ {
			game.GameMap = append(game.GameMap, []rune{})
			for h := 0; h < game.height; h++ {
				game.GameMap[w] = append(game.GameMap[w], ' ')
			}
		}

		// Fill the map
		for x, l := range lines {
			for y, c := range l {
				if c != ' ' {
					game.GameMap[y][x] = c
				}
			}
		}
	}
}

// While moving toward left or right, checks for obstacle on the way
func (g Game) CheckEmptyBlock(direction string) bool {
	upMostPoint := int(g.Player.Position.Y - g.Player.EatBox[0])
	rightMostPoint := int(g.Player.Position.X + g.Player.EatBox[1])
	downMostPoint := int(g.Player.Position.Y + g.Player.EatBox[2])
	leftMostPoint := int(g.Player.Position.X - g.Player.EatBox[3])
	x := int(g.Player.Position.X)
	y := int(g.Player.Position.Y)
	switch direction {
	case "up":
		if upMostPoint >= 0 {
			return !g.AllElements[g.GameMap[x][upMostPoint]].Solid
		}
	case "right":
		if rightMostPoint < g.width-1 {
			return !g.AllElements[g.GameMap[rightMostPoint][y]].Solid
		}
	case "down":
		if downMostPoint < g.height-1 {
			return !g.AllElements[g.GameMap[x][downMostPoint]].Solid
		}
	case "left":
		if leftMostPoint >= 0 {
			return !g.AllElements[g.GameMap[leftMostPoint][y]].Solid
		}
	}
	return false
}

// Loads all resources as element
func (game *Game) loadResources() {
	game.loadRessource("stone", 's', true, []ImagePosition{{0, 1, 0, 1}})
	game.loadRessource("dirt", 'd', true, []ImagePosition{{1, 2, 0, 1}})
	game.loadRessource("grass", 'g', true, []ImagePosition{{2, 3, 0, 1}})
	game.loadRessource("brick", 'b', true, []ImagePosition{{3, 4, 0, 1}})
	game.loadRessource("herb_1", 'h', false, []ImagePosition{{0, 1, 1, 2}, {1, 2, 1, 2},
		{2, 3, 1, 2}, {3, 4, 1, 2}})
	game.loadRessource("coin_1", 'c', false, []ImagePosition{{0, 1, 2, 3}, {1, 2, 2, 3},
		{2, 3, 2, 3}, {3, 4, 2, 3}, {4, 5, 2, 3}, {5, 6, 2, 3}})
	game.loadRessource("player", 'p', false, []ImagePosition{{0, 1, 3, 4}, {1, 2, 3, 4},
		{2, 3, 3, 4}, {3, 4, 3, 4}, {4, 5, 3, 4}, {5, 6, 3, 4}, {6, 7, 3, 4}, {7, 8, 3, 4}})
}

// Loads a single element
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

/////////////////////
// OTHER FUNCTIONS //
/////////////////////

// Find the longest string of an array (used for get the width of the map depending on)
func longestStr(arr []string) (max int) {
	for _, v := range arr {
		if len(v) > max {
			max = len(v)
		}
	}
	return
}
