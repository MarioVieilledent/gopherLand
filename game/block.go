package game

type Block struct {
	Name        string          // Name of block
	Short       rune            // Short identifier for blocks (to build maps)
	Solidity    Solidity        // Solid can be walked over
	Collectable bool            // Object can be collected
	Images      []ImagePosition // List of images for displaying
}

// Solidity enum
type Solidity string

const (
	Solid    Solidity = "Solid"
	Platform Solidity = "Platform"
	NotSolid Solidity = "NotSolid"
)

// Loads all blocks
func (game *Game) loadResources() {
	// Solid blocks
	game.loadRessource("stone", 's', Solid, false, []ImagePosition{{0, 1, 0, 1}})
	game.loadRessource("dirt", 'd', Solid, false, []ImagePosition{{1, 2, 0, 1}})
	game.loadRessource("grass", 'g', Solid, false, []ImagePosition{{2, 3, 0, 1}})
	game.loadRessource("brick", 'b', Solid, false, []ImagePosition{{3, 4, 0, 1}})
	game.loadRessource("door_closed", 'C', Solid, false, []ImagePosition{{4, 5, 0, 1}})

	// Not solid blocks
	game.loadRessource("pillar_up_down", '0', NotSolid, false, []ImagePosition{{0, 1, 4, 5}})
	game.loadRessource("pillar_down", '1', NotSolid, false, []ImagePosition{{1, 2, 4, 5}})
	game.loadRessource("pillar_up", '2', NotSolid, false, []ImagePosition{{2, 3, 4, 5}})
	game.loadRessource("pillar_central", '3', NotSolid, false, []ImagePosition{{3, 4, 4, 5}})
	game.loadRessource("tree", 't', NotSolid, false, []ImagePosition{{4, 5, 1, 2}, {5, 6, 1, 2},
		{6, 7, 1, 2}, {7, 8, 1, 2}})
	game.loadRessource("herb", 'h', NotSolid, false, []ImagePosition{{0, 1, 1, 2}, {1, 2, 1, 2},
		{2, 3, 1, 2}, {3, 4, 1, 2}})
	game.loadRessource("door_opened", 'O', NotSolid, false, []ImagePosition{{5, 6, 0, 1}})

	// Platforms
	game.loadRessource("platform_none", '_', Platform, false, []ImagePosition{{6, 7, 0, 1}})
	game.loadRessource("platform_left", '/', Platform, false, []ImagePosition{{7, 8, 0, 1}})
	game.loadRessource("platform_right", '\\', Platform, false, []ImagePosition{{8, 9, 0, 1}})
	game.loadRessource("platform_all", '-', Platform, false, []ImagePosition{{9, 10, 0, 1}})

	// Collectable items
	game.loadRessource("coin", 'c', NotSolid, true, []ImagePosition{{0, 1, 2, 3}, {1, 2, 2, 3},
		{2, 3, 2, 3}, {3, 4, 2, 3}, {4, 5, 2, 3}, {5, 6, 2, 3}})
	game.loadRessource("key", 'k', NotSolid, true, []ImagePosition{{6, 7, 2, 3}, {7, 8, 2, 3}})

	// Player
	game.loadRessource("player", 'p', NotSolid, false, []ImagePosition{{0, 1, 3, 4}, {1, 2, 3, 4},
		{2, 3, 3, 4}, {3, 4, 3, 4}, {4, 5, 3, 4}, {5, 6, 3, 4}, {6, 7, 3, 4}, {7, 8, 3, 4}})

	// Air
	game.loadRessource("air", ' ', NotSolid, false, []ImagePosition{})
}

// Loads a single block
func (game *Game) loadRessource(name string, short rune, solid Solidity, collectable bool, images []ImagePosition) {
	imagePos := []ImagePosition{}

	for _, v := range images {
		imagePos = append(imagePos, ImagePosition{
			game.BlockSize * v.X1,
			game.BlockSize * v.X2,
			game.BlockSize * v.Y1,
			game.BlockSize * v.Y2,
		})
	}

	game.AllBlocks[short] = Block{
		name,
		short,
		solid,
		collectable,
		imagePos,
	}
}
