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
	_%s_
	*%s* - *%s*

	`

	title := fmt.Sprintf("%s\n", m.League)
	str += title

	if len(m.Games) == 0 {
		str += "\n❌No games today❌\n\n"
		str += "----------------------\n"
		return str
	}

	for _, v := range m.Games {
		str += fmt.Sprintf(t, v.Time.Format("3:04 PM"), v.Home, v.Away)
	}

	str += "----------------------\n"
	return str
}

type GameModel struct {
	Home string
	Away string
	Time time.Time
}
