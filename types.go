package main

import "time"

var VDHOLDER string

type racer struct {
	Name         string      `json:"racername"`
	QualifyTime  int         `json:"racerQaulifyTime"` //Maybe for server side sorting later?
	Racetimgooes [10]float32 `json:"racerRaceTimes"`
}

var groupCounter = 1

type raceGroup struct {
	GroupId     int       `json:"rgGroupId"`
	GroupPhrase string    `json:"rgGroupPhrase"`
	GroupRev    int       `json:"rgGroupRev"`
	Racers      []racer   `json:"rgRacer"`
	Lifetime    time.Time `json:"rgTime"`
}

type racerGroupResponse struct {
	Groups []raceGroup `json:"raceGroups"`
}

var racerGroups []raceGroup
