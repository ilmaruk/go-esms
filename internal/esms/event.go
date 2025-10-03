package esms

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Emojis
const (
	YELLOW_CARD = '\U0001F7E8'
	RED_CARD    = '\U0001F7E5'
	SOCCER_BALL = '\U000026BD'
	CRUTCHES    = '\U0001F9F9'
)

type Event struct {
	Minute int
	Type   string
	Team   string
	Player Player
}

func NewEvent(minute int, team string, player Player) Event {
	return Event{
		Minute: minute,
		Team:   team,
		Player: player,
	}
}

type ChanceEvent struct {
	Event
	Assister *Player
	Tackler  *Player
	Outcome  string
}

func NewChanceEvent(minute int, team string) ChanceEvent {
	return ChanceEvent{
		Event: NewEvent(minute, team, Player{}),
	}
}

func (e *ChanceEvent) WithShooter(p Player) {
	e.Player = p
}

func (e *ChanceEvent) WithAssister(p Player) {
	e.Assister = &p
}

func (e *ChanceEvent) WithTackler(p Player) {
	e.Tackler = &p
}

func (e *ChanceEvent) WithOutcome(o string) {
	e.Outcome = o
}

func (e *ChanceEvent) String() string {
	home := color.New(color.FgBlue, color.Bold).SprintFunc()

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("[%2d - %s] ", e.Minute, e.Team))

	b.WriteString(fmt.Sprintf("Chance for %s", home(e.Player.Name.String())))
	if e.Assister != nil {
		b.WriteString(fmt.Sprintf(", assisted by %s:", home(e.Assister.Name.String())))
	}

	if e.Tackler != nil {
		away := color.New(color.FgRed, color.Bold).SprintFunc()
		b.WriteString(fmt.Sprintf(" tackled by %s!", away(e.Tackler.Name.String())))
		return b.String()
	}

	switch e.Outcome {
	case "OFFTARGET":
		b.WriteString(" the shot is off-target!")
	case "SAVE":
		b.WriteString(" the shot is saved!")
	default:
		b.WriteString(fmt.Sprintf(" it's in goal!!! %c", SOCCER_BALL))
	}

	return b.String()
}

type BookingEvent struct {
	Event
	Outcome string
}

func NewBookingEvent(minute int, team string, player Player) BookingEvent {
	return BookingEvent{
		Event: NewEvent(minute, team, player),
	}
}

func (e *BookingEvent) WithOutcome(o string) {
	e.Outcome = o
}

func (e BookingEvent) String() string {
	home := color.New(color.FgBlue, color.Bold).SprintFunc()

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("[%2d - %s] ", e.Minute, e.Team))

	switch e.Outcome {
	case "YELLOW":
		b.WriteString(fmt.Sprintf("%s is booked %c", home(e.Player.Name.String()), YELLOW_CARD))
	case "SECONDYELLOW":
		b.WriteString(fmt.Sprintf("%s is booked again and sent off %c%c", home(e.Player.Name.String()), YELLOW_CARD, RED_CARD))
	case "RED":
		b.WriteString(fmt.Sprintf("%s is sent off %c", home(e.Player.Name.String()), RED_CARD))
	}

	return b.String()
}
