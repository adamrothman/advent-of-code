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

var logLineRegexp = regexp.MustCompile(`^\[(\d{4}\-\d{2}\-\d{2} \d{2}:\d{2})\] ([\w #]+)$`)
var guardRegexp = regexp.MustCompile(`^Guard #(\d+) begins shift$`)

func parseLogLine(raw string) (LogLine, error) {
	matches := logLineRegexp.FindStringSubmatch(raw)
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
		matches = guardRegexp.FindStringSubmatch(line.Message)
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
		return nil, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	lines := make([]LogLine, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line, err := parseLogLine(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("parsing line: %s", err)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input file: %s", err)
	}

	return lines, nil
}

type timesAsleepPerMinute map[int]uint
type timesAsleepPerMinutePerGuard map[uint64]timesAsleepPerMinute

func countTimesAsleepPerMinutePerGuard(lines []LogLine) timesAsleepPerMinutePerGuard {
	timesAsleepByGuard := make(timesAsleepPerMinutePerGuard)

	var currentGuard uint64
	var fellAsleep time.Time
	for _, line := range lines {
		if line.Guard != 0 {
			currentGuard = line.Guard
			if _, ok := timesAsleepByGuard[currentGuard]; !ok {
				timesAsleepByGuard[currentGuard] = make(timesAsleepPerMinute)
			}
		} else if line.FallsAsleep {
			fellAsleep = line.Time
		} else if line.WakesUp && !fellAsleep.IsZero() {
			perMinute := timesAsleepByGuard[currentGuard]
			for minute := fellAsleep.Minute(); minute < line.Time.Minute(); minute++ {
				perMinute[minute]++
			}
			fellAsleep = time.Time{}
		}
	}

	return timesAsleepByGuard
}

func calculateSleepiestGuard(counts timesAsleepPerMinutePerGuard) (sleepiestGuard uint64, minutesAsleep uint, sleepiestMinute int) {
	for guard, minuteCounts := range counts {
		var guardSleepiestMinute int
		var guardMinutesAsleep, guardMaxTimesAsleep uint

		for minute, timesAsleepDuringMinute := range minuteCounts {
			if timesAsleepDuringMinute > guardMaxTimesAsleep {
				guardSleepiestMinute = minute
				guardMaxTimesAsleep = timesAsleepDuringMinute
			}
			guardMinutesAsleep += timesAsleepDuringMinute
		}

		if guardMinutesAsleep > minutesAsleep {
			sleepiestGuard = guard
			minutesAsleep = guardMinutesAsleep
			sleepiestMinute = guardSleepiestMinute
		}
	}

	return
}

func calculateTargetGuardAndMinute(counts timesAsleepPerMinutePerGuard) (targetGuard uint64, targetMinute int, timesAsleep uint) {
	for guard, minuteCounts := range counts {
		var guardSleepiestMinute int
		var guardMaxTimesAsleep uint

		for minute, timesAsleepDuringMinute := range minuteCounts {
			if timesAsleepDuringMinute > guardMaxTimesAsleep {
				guardSleepiestMinute = minute
				guardMaxTimesAsleep = timesAsleepDuringMinute
			}
		}

		if guardMaxTimesAsleep > timesAsleep {
			targetGuard = guard
			targetMinute = guardSleepiestMinute
			timesAsleep = guardMaxTimesAsleep
		}
	}

	return
}

func main() {
	filename := "input.txt"

	lines, err := readInput(filename)
	if err != nil {
		log.Fatalf("Error reading input from %s: %s\n", filename, err)
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].Time.Before(lines[j].Time)
	})

	timesAsleepPerMinutePerGuard := countTimesAsleepPerMinutePerGuard(lines)

	sleepiestGuard, minutesAsleep, sleepiestMinute := calculateSleepiestGuard(timesAsleepPerMinutePerGuard)
	fmt.Printf("Sleepiest guard is %d with %d minutes spent asleep\n", sleepiestGuard, minutesAsleep)
	fmt.Printf("Sleepiest minute for guard %d is %d\n", sleepiestGuard, sleepiestMinute)
	product := sleepiestGuard * uint64(sleepiestMinute)
	fmt.Printf("\t%d * %d = %d\n", sleepiestGuard, sleepiestMinute, product)

	targetGuard, targetMinute, timesAsleep := calculateTargetGuardAndMinute(timesAsleepPerMinutePerGuard)
	fmt.Printf("Guard %d spent minute %d asleep %d times\n", targetGuard, targetMinute, timesAsleep)
	product = targetGuard * uint64(targetMinute)
	fmt.Printf("\t%d * %d = %d\n", targetGuard, targetMinute, product)
}
