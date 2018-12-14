package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	dependencies := make(DependencyGraph)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var dependency, step string
		n, err := fmt.Sscanf(
			scanner.Text(),
			"Step %s must be finished before step %s can begin.",
			&dependency,
			&step,
		)
		if n != 2 || err != nil {
			return nil, fmt.Errorf("parsing line: %s", err)
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
		log.Fatalf("Error reading input from %s: %s\n", filename, err)
	}

	buildOrder := findBuildOrder(dependencies)
	fmt.Println("Build order:", strings.Join(buildOrder, ""))

	// Read dependencies again because it's modified by findBuildOrder
	dependencies, err = readInput(filename)
	if err != nil {
		log.Fatalf("Error reading input from %s: %s\n", filename, err)
	}

	workerCount := 5
	parallelTime := timeWork(dependencies, workerCount)
	fmt.Printf("Total time for %d parallel workers: %d\n", workerCount, parallelTime)
}
