package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func findStopByName(stops []Stop, stopName string) *Stop {
	stopNameLower := strings.ToLower(stopName)
	for i, stop := range stops {
		if strings.ToLower(stop.StopName) == stopNameLower {
			return &stops[i]
		}
	}
	return nil
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
