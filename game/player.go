package game

type Player struct {
	// Position and eat-box
	Position Position
	EatBox   [4][2]float64 // 4 points in rectangle around player

	// Speed and velocities for moving
	Speed                float64 // Speed when walking
	VelocityShortJump    float64 // Vertical velocity hen short jump
	VelocityDiffLongJump float64 // Add of vertical velocity for long jump

	// State
	Direction        rune    // l or r
	TouchingGround   bool    // True if player is walking, false if falling or jumping
	VerticalVelocity float64 // Vertical velocity on air (gravity falling, or gravity jumping)
	Walking          bool    // Animate player when walking

	// Stuff
	Gold      int      // Gold earned
	Keys      int      // Number of keys owned
	Inventory []Object // List of object
}

// Initialize a new player with default settings
func initPlayer(xPlayerFixed int) Player {
	return Player{
		Position{float64(xPlayerFixed) + 6.5, 2},
		[4][2]float64{
			{-0.3, -0.4},
			{0.3, -0.4},
			{0.3, 0.5},
			{-0.3, 0.5},
		},

		0.09,
		-17.0,
		-5.9,

		'r',
		false,
		0.0,
		false,

		0,
		0,
		[]Object{},
	}
}

// Move the player increasing X and Y, does not check if player can move (game's job)
func (p *Player) Move(x, y float64) {
	p.Position.X += x
	p.Position.Y += y
}

// Collects n golds
func (p *Player) CollectGold(number int) {
	p.Gold += number
}

// Return true if object can be bought, if yes, decreases Gold by cost
func (p *Player) Buy(cost int) (canBuy bool) {
	if p.Gold-cost >= 0 {
		p.Gold = -cost
		canBuy = true
	}
	return
}

func (p *Player) AddInventory(short rune) {
	switch short {
	case 'k':
		p.Inventory = append(p.Inventory, Object{"Key", "Into the unknown."})
	}
}
