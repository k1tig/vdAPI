// TODO
// - Middleware: Rate limit general calls + createGroups POST
// - Tidy up form validation responses
// - Make sure all responses follow standard resp
// - Strip passphrase from group calls

package main

import (
	"github.com/k1tig/vdAPI/middleware"
)

func main() {

	go groupTimeout() // removes groups that haven't been active from memory
	middleware.GetKeys()
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
