package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type LogLine struct {
	Time    time.Time
	Message string

	Guard       uint64
	FallsAsleep bool
	WakesUp     bool
}

const timeLayout = "2006-01-02 15:04"

var logLineRE = regexp.MustCompile(`^\[(\d{4}\-\d{2}\-\d{2} \d{2}:\d{2})\] ([\w #]+)$`)
var guardRE = regexp.MustCompile(`^Guard #(\d+) begins shift$`)

func parseLogLine(raw string) (LogLine, error) {
	matches := logLineRE.FindStringSubmatch(raw)
	if matches == nil {
		return LogLine{}, fmt.Errorf("log line string \"%s\" does not match pattern", raw)
	}

	ts, err := time.Parse(timeLayout, matches[1])
	if err != nil {
		return LogLine{}, fmt.Errorf("parsing time from log line: %s", err)
	}

	line := LogLine{Time: ts, Message: matches[2]}

	if line.Message == "falls asleep" {
		line.FallsAsleep = true
	} else if line.Message == "wakes up" {
		line.WakesUp = true
	} else {
		matches = guardRE.FindStringSubmatch(line.Message)
		if matches == nil {
			return LogLine{}, fmt.Errorf("guard change string \"%s\" does not match pattern", line.Message)
		}

		guard, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			return LogLine{}, fmt.Errorf("parsing guard ID: %s", err)
		}

		line.Guard = guard
	}

	return line, nil
}

func readInput(filename string) ([]LogLine, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file %s for reading: %s", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := make([]LogLine, 0)

	for scanner.Scan() {
		line, err := parseLogLine(scanner.Text())
		if err != nil {
			log.Printf("Error parsing line: %s", err)
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input file: %s", err)
	}

	return lines, nil
}

func calculateSleepiestGuard(lines []LogLine) (sleepiestGuard uint64, totalSleep time.Duration) {
	sleepCounter := make(map[uint64]time.Duration)

	var currentGuard uint64
	var fellAsleep time.Time
	for _, line := range lines {
		if line.Guard != 0 {
			currentGuard = line.Guard
		} else if line.FallsAsleep {
			fellAsleep = line.Time
		} else if line.WakesUp && !fellAsleep.IsZero() {
			durationAsleep := line.Time.Sub(fellAsleep)
			sleepCounter[currentGuard] += durationAsleep
			fellAsleep = time.Time{}
		}
	}

	for guard, durationAsleep := range sleepCounter {
		if durationAsleep > totalSleep {
			sleepiestGuard = guard
			totalSleep = durationAsleep
		}
	}

	return
}

func calculateSleepiestMinute(lines []LogLine, guard uint64) (sleepiestMinute int) {
	minuteTracker := make(map[int]uint64)

	var currentGuard uint64
	var fellAsleep time.Time
	for _, line := range lines {
		if line.Guard != 0 {
			currentGuard = line.Guard
		} else if line.FallsAsleep {
			fellAsleep = line.Time
		} else if line.WakesUp && !fellAsleep.IsZero() {
			if currentGuard == guard {
				for minute := fellAsleep.Minute(); minute < line.Time.Minute(); minute++ {
					minuteTracker[minute]++
				}
			}
			fellAsleep = time.Time{}
		}
	}

	var mostTimesAsleep uint64
	for minute, timesAsleep := range minuteTracker {
		if timesAsleep > mostTimesAsleep {
			sleepiestMinute = minute
			mostTimesAsleep = timesAsleep
		}
	}

	return
}

func calculateMostFrequentlyAsleepOnSameMinute(lines []LogLine) (targetGuard uint64, targetMinute int, timesAsleep uint64) {
	minuteTrackers := make(map[uint64]map[int]uint64)

	var currentGuard uint64
	var fellAsleep time.Time
	for _, line := range lines {
		if line.Guard != 0 {
			currentGuard = line.Guard
			if _, ok := minuteTrackers[currentGuard]; !ok {
				minuteTrackers[currentGuard] = make(map[int]uint64)
			}
		} else if line.FallsAsleep {
			fellAsleep = line.Time
		} else if line.WakesUp && !fellAsleep.IsZero() {
			tracker := minuteTrackers[currentGuard]
			for minute := fellAsleep.Minute(); minute < line.Time.Minute(); minute++ {
				tracker[minute]++
			}
			fellAsleep = time.Time{}
		}
	}

	for guard, tracker := range minuteTrackers {
		var sleepiestMinute int
		var mostTimesAsleep uint64
		for minute, timesAsleepDuringMinute := range tracker {
			if timesAsleepDuringMinute > mostTimesAsleep {
				sleepiestMinute = minute
				mostTimesAsleep = timesAsleepDuringMinute
			}
		}

		if mostTimesAsleep > timesAsleep {
			targetGuard = guard
			targetMinute = sleepiestMinute
			timesAsleep = mostTimesAsleep
		}
	}

	return
}

func main() {
	filename := "input.txt"

	lines, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input values from %s: %s\n", filename, err)
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].Time.Before(lines[j].Time)
	})

	sleepiestGuard, totalSleep := calculateSleepiestGuard(lines)
	fmt.Printf("Sleepiest guard is %d with %v spent asleep\n", sleepiestGuard, totalSleep)

	sleepiestMinute := calculateSleepiestMinute(lines, sleepiestGuard)
	fmt.Printf("Sleepiest minute for guard %d is %d\n", sleepiestGuard, sleepiestMinute)

	targetGuard, targetMinute, timesAsleep := calculateMostFrequentlyAsleepOnSameMinute(lines)
	fmt.Printf("Guard %d spend minute %d asleep %d times\n", targetGuard, targetMinute, timesAsleep)
}
