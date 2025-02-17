package main

import "time"

var groupCounter = 1

type racer struct {
	Name        string      `json:"racername"`
	QualifyTime int         `json:"racerQaulifyTime"` //Maybe for server side sorting later?
	Racetimes   [10]float32 `json:"racerRaceTimes"`
}

type raceGroup struct {
	GroupId     int       `json:"rgGroupId"`
	GroupPhrase string    `json:"rgGroupPhrase"`
	GroupRev    int       `json:"rgGroupRev"`
	Racers      []racer   `json:"rgRacer"`
	Livetime    time.Time `json:"rgTime"`
}

// Lifetime counter for in program data
var maxTime time.Duration = 1 * time.Minute

type racerGroupResponse struct {
	Groups []raceGroup `json:"raceGroups"`
}

var racerGroups []raceGroup

func main() {

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if racerGroups != nil {
					for index, item := range racerGroups {
						if time.Since(item.Livetime) > maxTime {
							racerGroups = append(racerGroups[:index], racerGroups[index+1:]...)
						}
					}
				}
			}
		}
	}()
	server := NewAPIServer((":8080"))
	server.Run()
}

/*

{
  "rgGroupPhrase": "banana",
  "rgGroupRev": 1,
  "rgRacer":[
    {"racerName": "Eedok",
    "racerQualifytime": 69.420,
    "racerRacetimes": [0,0,0,0,0,0,0,0,0,0]
    },
    {"racerName": "MrE",
    "racerQualifytime": 69.113,
    "racerRacetimes": [0,0,0,0,0,0,0,0,0,0]
    },
     {"racername": "JonE4",
    "racerQualifytime": 42.069,
    "racerRacetimes": [0,0,0,0,0,0,0,0,0,0]
    }

  ]

}

*/
