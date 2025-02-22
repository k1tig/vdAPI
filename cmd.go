package main

import (
	"fmt"
	"time"
)

// Group timeout limit for memory life
// timer set short for testing
var maxTime time.Duration = 30 * time.Second

func groupTimeout() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	shutUpForError := make(chan bool)

	for {
		select {
		case <-ticker.C:
			for index, group := range racerGroups {
				if time.Since(group.Lifetime) > maxTime {
					//remove expired group
					racerGroups = append(racerGroups[:index], racerGroups[index+1:]...)
					break
				}
			}
		//This has to be here to shut up the for/while loop not liking there only being one case
		case <-shutUpForError:
			fmt.Println("All this to shut up the Error")
		}
	}
}
