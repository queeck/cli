package state

import "time"

type State interface {
	Time() time.Time
}

type state struct {
	//
}

func New() State {
	return &state{
		//
	}
}

func (s *state) Time() time.Time {
	return time.Now()
}
