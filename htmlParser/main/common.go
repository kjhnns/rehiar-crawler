package main

import "time"

type Debugging int

const (
	DEBUG_NONE      Debugging = 0
	DEBUG_VERBOSE   Debugging = 1
	DEBUG_DUMMY     Debugging = 2
	DEBUG_DEBUGGING Debugging = 4
	DEBUG_REPAIR    Debugging = 8
)

type Settings struct {
	StartTime time.Time
	Debugging Debugging
}

var Configuration Settings
