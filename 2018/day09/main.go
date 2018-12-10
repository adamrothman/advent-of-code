package main

import (
	"bufio"
	"container/ring"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Game struct {
	PlayerCount     int
	LastMarbleValue int
}

func readInput(filename string) (Game, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return Game{}, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return Game{}, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	var playerCount, lastMarbleValue int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n, err := fmt.Sscanf(scanner.Text(), "%d players; last marble is worth %d points", &playerCount, &lastMarbleValue)
		if n != 2 || err != nil {
			return Game{}, fmt.Errorf("parsing line: %s", err)
		}
		break
	}
	if err := scanner.Err(); err != nil {
		return Game{}, fmt.Errorf("reading input file: %s", err)
	}

	game := Game{PlayerCount: playerCount, LastMarbleValue: lastMarbleValue}
	return game, nil
}

func play(game Game) map[int]int {
	scores := make(map[int]int)
	// Players are 1-indexed
	for p := 1; p <= game.PlayerCount; p++ {
		scores[p] = 0
	}

	zero := ring.New(1)
	zero.Value = 0

	current := zero

	for nextMarbleValue := 1; nextMarbleValue <= game.LastMarbleValue; nextMarbleValue++ {
		player := ((nextMarbleValue - 1) % game.PlayerCount) + 1

		if nextMarbleValue%23 == 0 {
			removed := current.Move(-8).Unlink(1)
			current = current.Move(-6)

			scores[player] += nextMarbleValue
			scores[player] += removed.Value.(int)
		} else {
			nextMarble := ring.New(1)
			nextMarble.Value = nextMarbleValue

			current.Next().Link(nextMarble)
			current = nextMarble
		}
	}

	return scores
}

func getWinner(scores map[int]int) (winner, highScore int) {
	for player, score := range scores {
		if score > highScore {
			winner = player
			highScore = score
		}
	}
	return
}

func main() {
	filename := "input.txt"

	game, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	scores := play(game)
	winner, highScore := getWinner(scores)
	fmt.Printf("Winner of %+v is player %d with %d points\n", game, winner, highScore)

	bigGame := Game{PlayerCount: game.PlayerCount, LastMarbleValue: game.LastMarbleValue * 100}
	scores = play(bigGame)
	winner, highScore = getWinner(scores)
	fmt.Printf("Winner of %+v is player %d with %d points\n", bigGame, winner, highScore)
}
