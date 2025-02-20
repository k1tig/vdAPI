package main

import (
	"time"
)

// Group timeout limit for memory life
var maxTime time.Duration = 15 * time.Second

func groupTimeout() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if racerGroups != nil {
				for index, item := range racerGroups {
					if time.Since(item.Livetime) > maxTime {
						racerGroups = append(racerGroups[:index], racerGroups[index+1:]...)
						break
					}
				}
			}
		}
	}
}
