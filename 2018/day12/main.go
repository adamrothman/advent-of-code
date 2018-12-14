package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type State []bool

type Generation struct {
	State     State
	ZeroIndex int
}

func (g Generation) Sum() (sum int) {
	for i, hasPlant := range g.State {
		if hasPlant {
			sum += i - g.ZeroIndex
		}
	}
	return
}

type Rule struct {
	Input  State
	Output bool
}

func (r Rule) Matches(window State) bool {
	if len(window) != len(r.Input) {
		return false
	}
	for i := 0; i < len(r.Input); i++ {
		if window[i] != r.Input[i] {
			return false
		}
	}
	return true
}

func readInput(filename string) (Generation, []Rule, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return Generation{}, nil, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return Generation{}, nil, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	var initial Generation
	var rules = make([]Rule, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "initial state:") {
			initial.State, err = parseInitialState(line)
			if err != nil {
				return Generation{}, nil, fmt.Errorf("parsing initial state: %s", err)
			}
		} else if len(line) == 0 {
			continue
		} else {
			rule, err := parseRule(line)
			if err != nil {
				return Generation{}, nil, fmt.Errorf("parsing rule: %s", err)
			}
			rules = append(rules, rule)
		}
	}
	if err := scanner.Err(); err != nil {
		return Generation{}, nil, fmt.Errorf("reading input file: %s", err)
	}

	return initial, rules, nil
}

const hashChar = 35

func parseState(raw string) State {
	state := make(State, len(raw))
	for i, c := range raw {
		if c == hashChar {
			state[i] = true
		}
	}
	return state
}

var initialStateRegexp = regexp.MustCompile(`^initial state: (?P<state>[#.]+)$`)

func parseInitialState(raw string) (State, error) {
	matches := initialStateRegexp.FindStringSubmatch(raw)
	if matches == nil {
		return nil, fmt.Errorf("initial state string \"%s\" does not match expected pattern", raw)
	}

	state := parseState(matches[1])
	return state, nil
}

var ruleRegexp = regexp.MustCompile(`^(?P<input>[#.]{5}) => (?P<output>[#.])$`)

func parseRule(raw string) (Rule, error) {
	matches := ruleRegexp.FindStringSubmatch(raw)
	if matches == nil {
		return Rule{}, fmt.Errorf("rule string \"%s\" does not match expected pattern", raw)
	}

	rule := Rule{
		Input:  parseState(matches[1]),
		Output: matches[2][0] == hashChar,
	}
	return rule, nil
}

const windowSize = 5
const paddingSize = windowSize - 1

func simulateGrowth(initial Generation, rules []Rule, generations int) (result Generation) {
	result = expand(initial)
	for gen := 0; gen < generations; gen++ {
		result = simulateGeneration(result, rules)
	}
	return
}

func expand(gen Generation) (expanded Generation) {
	lowestPlant, highestPlant := findLowestAndHighestPlants(gen.State)
	var headPadding, tailPadding int
	if lowestPlant != nil && *lowestPlant < paddingSize {
		headPadding = paddingSize - *lowestPlant
	}
	if highestPlant != nil && *highestPlant >= len(gen.State)-paddingSize {
		tailPadding = *highestPlant + windowSize - len(gen.State)
	}

	expanded = Generation{ZeroIndex: gen.ZeroIndex}

	paddingToAdd := headPadding + tailPadding
	if paddingToAdd > 0 {
		expanded.State = make(State, len(gen.State)+paddingToAdd)
		copy(expanded.State[headPadding:], gen.State)
		expanded.ZeroIndex += headPadding
	} else {
		expanded.State = gen.State
	}

	return
}

func findLowestAndHighestPlants(state State) (lowest, highest *int) {
	for i := 0; i < len(state); i++ {
		if state[i] && (lowest == nil || i < *lowest) {
			idx := i
			lowest = &idx
			break
		}
	}
	for i := len(state) - 1; i >= 0; i-- {
		if state[i] && (highest == nil || i > *highest) {
			idx := i
			highest = &idx
			break
		}
	}
	return
}

func simulateGeneration(g Generation, rules []Rule) Generation {
	next := Generation{
		State:     make(State, len(g.State)),
		ZeroIndex: g.ZeroIndex,
	}

	for i := 0; i+windowSize < len(g.State); i++ {
		current := i + windowSize/2
		window := g.State[i : i+windowSize]

		for _, rule := range rules {
			if rule.Matches(window) {
				next.State[current] = rule.Output
				break
			}
		}
	}

	return expand(next)
}

func main() {
	filename := "input.txt"

	initial, rules, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	twentyGens := simulateGrowth(initial, rules, 20)
	fmt.Printf("Sum of numbers of all pots with plants after 20 generations: %d\n", twentyGens.Sum())
}
