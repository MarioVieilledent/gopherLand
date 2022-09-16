package game

import (
	"io/ioutil"
	"log"
	"strings"
)

const mapPath string = "data/maps/map.txt"

// Creates the map (array of array of runes) based on the txt map file
func (game *Game) createMap() {
	// Reads the resources file
	file, err := ioutil.ReadFile(mapPath)
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
