package game

type Block struct {
	Name        string          // Name of block
	Short       rune            // Short identifier for blocks (to build maps)
	Solid       bool            // Solid can be walked over
	Collectable bool            // Object can be collected
	Images      []ImagePosition // List of images for displaying
}

// Loads all blocks
func (game *Game) loadResources() {
	game.loadRessource("stone", 's', true, false, []ImagePosition{{0, 1, 0, 1}})
	game.loadRessource("dirt", 'd', true, false, []ImagePosition{{1, 2, 0, 1}})
	game.loadRessource("grass", 'g', true, false, []ImagePosition{{2, 3, 0, 1}})
	game.loadRessource("brick", 'b', true, false, []ImagePosition{{3, 4, 0, 1}})

	game.loadRessource("pillar_up_down", '0', false, false, []ImagePosition{{0, 1, 4, 5}})
	game.loadRessource("pillar_down", '1', false, false, []ImagePosition{{1, 2, 4, 5}})
	game.loadRessource("pillar_up", '2', false, false, []ImagePosition{{2, 3, 4, 5}})
	game.loadRessource("pillar_central", '3', false, false, []ImagePosition{{3, 4, 4, 5}})

	game.loadRessource("herb", 'h', false, false, []ImagePosition{{0, 1, 1, 2}, {1, 2, 1, 2},
		{2, 3, 1, 2}, {3, 4, 1, 2}})
	game.loadRessource("coin", 'c', false, true, []ImagePosition{{0, 1, 2, 3}, {1, 2, 2, 3},
		{2, 3, 2, 3}, {3, 4, 2, 3}, {4, 5, 2, 3}, {5, 6, 2, 3}})
	game.loadRessource("player", 'p', false, false, []ImagePosition{{0, 1, 3, 4}, {1, 2, 3, 4},
		{2, 3, 3, 4}, {3, 4, 3, 4}, {4, 5, 3, 4}, {5, 6, 3, 4}, {6, 7, 3, 4}, {7, 8, 3, 4}})
}

// Loads a single block
func (game *Game) loadRessource(name string, short rune, solid bool, collectable bool, images []ImagePosition) {
	imagePos := []ImagePosition{}

	for _, v := range images {
		imagePos = append(imagePos, ImagePosition{
			game.Ss * v.X1,
			game.Ss * v.X2,
			game.Ss * v.Y1,
			game.Ss * v.Y2,
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
