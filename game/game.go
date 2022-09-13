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

type Game struct {
	Ss          int           // Size of square elements
	width       int           // Number of blocks (width)
	height      int           // Number of blocks (height)
	AllElements map[rune]Elem // All elements
	GameMap     [][]rune      // Game map
	Player      Player
}

type Player struct {
	Position         Position
	EatBox           [4][2]float64 // 4 points in rectangle around player
	Speed            float64
	Direction        rune    // l or r
	TouchingGround   bool    // True if player is walking, false if falling or jumping
	VerticalVelocity float64 // Vertical velocity on air (gravity falling, or gravity jumping)
	Walking          bool    // Animate player when walking
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
			[4][2]float64{
				{-0.3, -0.4},
				{0.3, -0.4},
				{0.3, 0.5},
				{-0.3, 0.5},
			},
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
func (g Game) CheckIfTouchesGround() bool {
	xDownRight := int(g.Player.Position.X + g.Player.EatBox[2][0])
	yDownRight := int(g.Player.Position.Y + g.Player.EatBox[2][1])

	xDownLeft := int(g.Player.Position.X + g.Player.EatBox[3][0])
	yDownLeft := int(g.Player.Position.Y + g.Player.EatBox[3][1])

	if yDownLeft < g.height-1 && yDownRight < g.height-1 {
		if g.AllElements[g.GameMap[xDownLeft][yDownLeft]].Solid ||
			g.AllElements[g.GameMap[xDownRight][yDownRight]].Solid {
			g.Player.TouchingGround = true
			return true
		} else {
			g.Player.TouchingGround = false
		}
	}
	return false
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

// Returns coordinates of each 4 points of player's eatbox
func (g Game) GetEatBoxPoints() (xUpLeft, yUpLeft, xUpRight, yUpRight,
	xDownRight, yDownRight, xDownLeft, yDownLeft int) {
	xUpLeft = int(g.Player.Position.X + g.Player.EatBox[0][0])
	yUpLeft = int(g.Player.Position.Y + g.Player.EatBox[0][1])

	xUpRight = int(g.Player.Position.X + g.Player.EatBox[1][0])
	yUpRight = int(g.Player.Position.Y + g.Player.EatBox[1][1])

	xDownRight = int(g.Player.Position.X + g.Player.EatBox[2][0])
	yDownRight = int(g.Player.Position.Y + g.Player.EatBox[2][1])

	xDownLeft = int(g.Player.Position.X + g.Player.EatBox[3][0])
	yDownLeft = int(g.Player.Position.Y + g.Player.EatBox[3][1])

	return
}

// Moves the player (checking if space is available)
func (g *Game) Move(x, y float64) (moving bool) {

	xUpLeft := int(g.Player.Position.X + g.Player.EatBox[0][0] + x)
	yUpLeft := int(g.Player.Position.Y + g.Player.EatBox[0][1] + y)

	xUpRight := int(g.Player.Position.X + g.Player.EatBox[1][0] + x)
	yUpRight := int(g.Player.Position.Y + g.Player.EatBox[1][1] + y)

	xDownRight := int(g.Player.Position.X + g.Player.EatBox[2][0] + x)
	yDownRight := int(g.Player.Position.Y + g.Player.EatBox[2][1] + y)

	xDownLeft := int(g.Player.Position.X + g.Player.EatBox[3][0] + x)
	yDownLeft := int(g.Player.Position.Y + g.Player.EatBox[3][1] + y)

	if x > 0 {
		if xUpLeft < g.width-1 && xDownLeft < g.width-1 {
			if !g.AllElements[g.GameMap[xUpRight][yUpRight]].Solid &&
				!g.AllElements[g.GameMap[xDownRight][yDownRight]].Solid {
				g.Player.Position.X += x
				if g.Player.TouchingGround {
					moving = true
				}
			}
		}
	} else if x < 0 {
		if xUpRight > 0 && xDownRight > 0 {
			if !g.AllElements[g.GameMap[xUpLeft][yUpLeft]].Solid &&
				!g.AllElements[g.GameMap[xDownLeft][yDownLeft]].Solid {
				g.Player.Position.X += x
				if g.Player.TouchingGround {
					moving = true
				}
			}
		}
	}

	if y > 0 {
		if yUpLeft < g.height-1 && yUpRight < g.height-1 {
			if !g.AllElements[g.GameMap[xDownRight][yDownRight]].Solid &&
				!g.AllElements[g.GameMap[xDownLeft][yDownLeft]].Solid {
				g.Player.Position.Y += y
				moving = true
			} else {
				// If player hits the ground
				g.Player.TouchingGround = true
				g.Player.VerticalVelocity = 0.0 // Reset the velocity of player
			}
		}
	} else if y < 0 {
		if yDownRight > 0 && yDownLeft > 0 {
			if !g.AllElements[g.GameMap[xUpLeft][yUpLeft]].Solid &&
				!g.AllElements[g.GameMap[xUpRight][yUpRight]].Solid {
				g.Player.Position.Y += y
				moving = true
			} else {
				// If player hit a ceil
				g.Player.VerticalVelocity = 0 // Reset of its velocity
			}
		}
	}

	return
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
