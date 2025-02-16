package main

type racer struct {
	Name        string      `json:"racername"`
	QualifyTime int         `json:"racerQaulifyTime"` //Maybe for server side sorting later?
	Racetimes   [10]float32 `json:"racerRaceTimes"`
}

type raceGroup struct {
	GroupId  int     `json:"rgGroupId"`
	GroupRev int     `json:"rgGroupRev"`
	Racers   []racer `json:"rgRacer"`
}

var racerGroups []raceGroup

func main() {
	server := NewAPIServer((":8080"))
	server.Run()
}

/*

{
  "rgGroupId": 1,
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
