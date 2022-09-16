package game

import "math"

const mapPath string = "data/maps/map.txt"

// Describe the image of an block, an entity, an object
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
	Ss        int            // Square size of blocks
	width     int            // Number of blocks (width)
	height    int            // Number of blocks (height)
	AllBlocks map[rune]Block // All blocks
	GameMap   [][]rune       // Game map
	Player    Player
}

// Create all the structures and arrays to initialize the game
func InitGame(xPlayerFixed int) Game {
	// Init game structure
	game := Game{
		64,
		0,
		0,
		map[rune]Block{},
		[][]rune{},
		initPlayer(xPlayerFixed),
	}
	game.loadResources()
	game.createMap()
	return game
}

///////////////////
// GAME METHODS ///
///////////////////

// Returns coordinates of each 4 points of player's eat-box
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

// Check if player is touching ground
func (g *Game) TouchesGround() bool {
	_, frac := math.Modf(g.Player.Position.Y)
	if frac > 0.4999 && frac < 0.5000 {
		xDownRight := int(g.Player.Position.X)
		yDownRight := int(g.Player.Position.Y) + 1
		var bDownRight Block
		if g.outOfMap([]int{xDownRight}, []int{yDownRight}) {
			bDownRight = g.AllBlocks[' ']
		} else {
			bDownRight = g.AllBlocks[g.GameMap[xDownRight][yDownRight]]
		}

		xDownLeft := int(g.Player.Position.X)
		yDownLeft := int(g.Player.Position.Y) + 1
		var bDownLeft Block
		if g.outOfMap([]int{xDownLeft}, []int{yDownLeft}) {
			bDownLeft = g.AllBlocks[' ']
		} else {
			bDownLeft = g.AllBlocks[g.GameMap[xDownLeft][yDownLeft]]
		}

		// If player fall to the platform, the platform is solid
		if (bDownLeft.Solidity == Platform && bDownRight.Solidity == NotSolid) ||
			(bDownLeft.Solidity == NotSolid && bDownRight.Solidity == Platform) ||
			(bDownLeft.Solidity == Platform && bDownRight.Solidity == Platform) {
			if bDownLeft.Solidity != Solid {
				g.Player.VerticalVelocity = 0.0 // Reset the velocity of player
				g.Player.TouchingGround = true
				return true
			}
		}

		if bDownLeft.Solidity == Solid && bDownRight.Solidity == Solid {
			g.Player.VerticalVelocity = 0.0 // Reset the velocity of player
			g.Player.TouchingGround = true
			return true
		}
	}
	g.Player.TouchingGround = false
	return false
}

// Moves the player (checking if space is available)
func (g *Game) Move(x, y float64) (moving bool) {

	xUpLeft := int(g.Player.Position.X + g.Player.EatBox[0][0] + x)
	yUpLeft := int(g.Player.Position.Y + g.Player.EatBox[0][1] + y)
	var bUpLeft Block
	if g.outOfMap([]int{xUpLeft}, []int{yUpLeft}) {
		bUpLeft = g.AllBlocks[' ']
	} else {
		bUpLeft = g.AllBlocks[g.GameMap[xUpLeft][yUpLeft]]
	}

	xUpRight := int(g.Player.Position.X + g.Player.EatBox[1][0] + x)
	yUpRight := int(g.Player.Position.Y + g.Player.EatBox[1][1] + y)
	var bUpRight Block
	if g.outOfMap([]int{xUpRight}, []int{yUpRight}) {
		bUpRight = g.AllBlocks[' ']
	} else {
		bUpRight = g.AllBlocks[g.GameMap[xUpRight][yUpRight]]
	}

	xDownRight := int(g.Player.Position.X + g.Player.EatBox[2][0] + x)
	yDownRight := int(g.Player.Position.Y + g.Player.EatBox[2][1] + y)
	var bDownRight Block
	if g.outOfMap([]int{xDownRight}, []int{yDownRight}) {
		bDownRight = g.AllBlocks[' ']
	} else {
		bDownRight = g.AllBlocks[g.GameMap[xDownRight][yDownRight]]
	}

	xDownLeft := int(g.Player.Position.X + g.Player.EatBox[3][0] + x)
	yDownLeft := int(g.Player.Position.Y + g.Player.EatBox[3][1] + y)
	var bDownLeft Block
	if g.outOfMap([]int{xDownLeft}, []int{yDownLeft}) {
		bDownLeft = g.AllBlocks[' ']
	} else {
		bDownLeft = g.AllBlocks[g.GameMap[xDownLeft][yDownLeft]]
	}

	g.Collect() // Collect items if player is on collectable item

	g.Action() // Do action (like open a door with a key)

	if x > 0 {
		if xUpLeft < g.width && xDownLeft < g.width {
			if (bUpRight.Solidity == NotSolid || bUpRight.Solidity == Platform) &&
				(bDownRight.Solidity == NotSolid || bDownRight.Solidity == Platform) {
				g.Player.Move(x, 0.0)
				if g.Player.TouchingGround {
					moving = true
				}
			}
		}
	} else if x < 0 {
		if xUpRight >= 0 && xDownRight >= 0 {
			if (bUpLeft.Solidity == NotSolid || bUpLeft.Solidity == Platform) &&
				(bDownLeft.Solidity == NotSolid || bDownLeft.Solidity == Platform) {
				g.Player.Move(x, 0.0)
				if g.Player.TouchingGround {
					moving = true
				}
			}
		}
	}

	if y > 0 {
		if yUpLeft < g.height && yUpRight < g.height {
			if (bDownLeft.Solidity == NotSolid || bDownLeft.Solidity == Platform) &&
				(bDownRight.Solidity == NotSolid || bDownRight.Solidity == Platform) {
				g.Player.Move(0.0, y)
				moving = true
			} else {
				// If player hits the ground
				g.Player.TouchingGround = true
				g.Player.VerticalVelocity = 0.0 // Reset the velocity of player
			}
		}
	} else if y < 0 {
		if yDownRight >= 0 && yDownLeft >= 0 {
			if (bUpLeft.Solidity == NotSolid || bUpLeft.Solidity == Platform) &&
				(bUpRight.Solidity == NotSolid || bUpRight.Solidity == Platform) {
				g.Player.Move(0.0, y)
				moving = true
			} else {
				// If player hit a ceil
				g.Player.VerticalVelocity = 0 // Reset of its velocity
			}
		}
	}

	return
}

// Checks if player is over a collectable item, if yes, collects it
func (g *Game) Collect() {
	x := int(g.Player.Position.X)
	y := int(g.Player.Position.Y)
	if !g.outOfMap([]int{x}, []int{y}) {
		b := g.AllBlocks[g.GameMap[x][y]]
		if b.Collectable {
			g.GameMap[x][y] = ' '
			switch b.Short {
			case 'c':
				g.Player.CollectGold(1)
			case 'k':
				g.Player.Keys++
			}
		}
	}
}

// Checks if player is over an element that has an action and does it
func (g *Game) Action() {
	x := int(g.Player.Position.X)
	y := int(g.Player.Position.Y)

	// Block directly at the left of player
	if !g.outOfMap([]int{x - 1}, []int{y}) {
		b := g.AllBlocks[g.GameMap[x-1][y]]

		// If blocks is a closed door
		if b.Short == 'C' && g.Player.Keys > 0 {
			g.GameMap[x-1][y] = 'O'
			g.Player.Keys--
		}
	}

	// Block directly at the right of player
	if !g.outOfMap([]int{x + 1}, []int{y}) {
		b := g.AllBlocks[g.GameMap[x+1][y]]

		// If blocks is a closed door
		if b.Short == 'C' && g.Player.Keys > 0 {
			g.GameMap[x+1][y] = 'O'
			g.Player.Keys--
		}
	}
}

// When down key is pressed, check if block underneath is a platform
func (g *Game) GoDown() {
	x := int(g.Player.Position.X)
	y := int(g.Player.Position.Y) + 1
	if !g.outOfMap([]int{x}, []int{y}) {
		if g.AllBlocks[g.GameMap[x][y]].Solidity == Platform {
			g.Player.Move(0, 0.1)
			g.Player.TouchingGround = false
		}
	}
}

// Checks if coordinates are inside the map to not get an error out of bounds
func (g *Game) outOfMap(x []int, y []int) bool {
	for _, v := range x {
		if v < 0 {
			return true
		}
		if v >= g.width {
			return true
		}
	}
	for _, v := range y {
		if v < 0 {
			return true
		}
		if v >= g.height {
			return true
		}
	}
	return false
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
