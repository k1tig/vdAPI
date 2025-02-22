package main

import (
	"time"
)

// Group timeout limit for memory life
// timer set short for testing
var maxTime time.Duration = 30 * time.Second

func groupTimeout() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for index, group := range racerGroups {
			if time.Since(group.Lifetime) > maxTime {
				//remove expired group
				racerGroups = append(racerGroups[:index], racerGroups[index+1:]...)
				break
			}
		}
	}
}
