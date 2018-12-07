package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type StringSet map[string]bool
type DependencyGraph map[string]StringSet

func readInput(filename string) (DependencyGraph, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	dependencies := make(DependencyGraph)

	for scanner.Scan() {
		dependency, step, err := parseLine(scanner.Text())
		if err != nil {
			log.Printf("Error parsing line: %s", err)
			continue
		}

		if _, ok := dependencies[step]; !ok {
			dependencies[step] = make(StringSet)
		}
		if _, ok := dependencies[dependency]; !ok {
			dependencies[dependency] = make(StringSet)
		}

		dependencies[step][dependency] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input file: %s", err)
	}

	return dependencies, nil
}

var lineRegexp = regexp.MustCompile(`^Step ([A-Z]) must be finished before step ([A-Z]) can begin.$`)

func parseLine(raw string) (dependency, step string, err error) {
	matches := lineRegexp.FindStringSubmatch(raw)
	if matches == nil {
		err = fmt.Errorf("line did not match expected pattern")
		return
	}

	dependency = matches[1]
	step = matches[2]

	return
}

func findBuildOrder(dependencies DependencyGraph) []string {
	order := make([]string, 0, len(dependencies))

	for len(dependencies) > 0 {
		available := findAvailableSteps(dependencies)

		next := available[0]
		order = append(order, next)

		removeDependency(next, dependencies)
	}

	return order
}

func findAvailableSteps(dependencies DependencyGraph) []string {
	steps := make([]string, 0)
	for step, deps := range dependencies {
		if len(deps) == 0 {
			steps = append(steps, step)
		}
	}
	sort.Strings(steps)
	return steps
}

func removeDependency(step string, dependencies DependencyGraph) {
	delete(dependencies, step)
	for _, deps := range dependencies {
		delete(deps, step)
	}
}

type Worker struct {
	Step          string
	TimeRemaining uint8
}

func (w Worker) Working() bool {
	return w.Step != ""
}

func timeWork(dependencies DependencyGraph, workerCount int) (total uint) {
	workers := make([]Worker, workerCount)
	assigned := make(StringSet)

	for len(dependencies) > 0 || len(assigned) > 0 {
		available := findAvailableSteps(dependencies)

		// Assign work
		for i := 0; i < len(workers); i++ {
			w := &workers[i]

			if w.Working() {
				continue
			}

			// Pick an available step; must not already be assigned
			var chosenStep string
			for _, step := range available {
				if !assigned[step] {
					chosenStep = step
					break
				}
			}

			if chosenStep == "" {
				// All available steps are already assigned
				break
			}

			w.Step = chosenStep
			w.TimeRemaining = timeForStep(chosenStep)

			assigned[chosenStep] = true
		}

		// Tick the clock
		for i := 0; i < len(workers); i++ {
			w := &workers[i]

			if !w.Working() {
				continue
			}

			w.TimeRemaining--
			if w.TimeRemaining == 0 {
				completedStep := w.Step
				w.Step = ""
				delete(assigned, completedStep)
				removeDependency(completedStep, dependencies)
			}
		}

		total++
	}

	return
}

func timeForStep(step string) uint8 {
	return step[0] - 4
}

func main() {
	filename := "input.txt"

	dependencies, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	buildOrder := findBuildOrder(dependencies)
	fmt.Println("Build order:", strings.Join(buildOrder, ""))

	// Read dependencies again because it's modified by findBuildOrder
	dependencies, err = readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	workerCount := 5
	parallelTime := timeWork(dependencies, workerCount)
	fmt.Printf("Total time for %d parallel workers: %d\n", workerCount, parallelTime)
}
