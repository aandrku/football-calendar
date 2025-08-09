package main

import (
	"fmt"
	"time"
)

type Model struct {
	League string
	Games  []GameModel
}

func (m Model) Format() string {
	var str string
	t := `
	%s
	%s - %s

	`
	title := m.League
	str += title

	if len(m.Games) == 0 {
		str += "\nNo games today"
		return str
	}

	for _, v := range m.Games {
		str += fmt.Sprintf(t, v.Time.Format("3:04 PM"), v.Home, v.Away)
	}

	return str
}

type GameModel struct {
	Home string
	Away string
	Time time.Time
}
