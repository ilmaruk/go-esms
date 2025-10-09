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
	Team   Team
	Player TeamPlayer
}

func NewEvent(minute int, team Team, player TeamPlayer) Event {
	return Event{
		Minute: minute,
		Team:   team,
		Player: player,
	}
}

type ChanceEvent struct {
	Event
	Assister *TeamPlayer
	Tackler  *TeamPlayer
	Outcome  string
}

func NewChanceEvent(minute int, team Team) ChanceEvent {
	return ChanceEvent{
		Event: NewEvent(minute, team, TeamPlayer{}),
	}
}

func (e *ChanceEvent) WithShooter(p TeamPlayer) {
	e.Player = p
}

func (e *ChanceEvent) WithAssister(p TeamPlayer) {
	e.Assister = &p
}

func (e *ChanceEvent) WithTackler(p TeamPlayer) {
	e.Tackler = &p
}

func (e *ChanceEvent) WithOutcome(o string) {
	e.Outcome = o
}

func (e *ChanceEvent) String() string {
	teamColor := color.New(colorToBg(e.Team.Colors[0]), colorToFg(e.Team.Colors[1]), color.Bold).SprintFunc()
	home := color.New(color.FgBlue, color.Bold).SprintFunc()

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s [%2dm] %c ", teamColor(strings.ToUpper(e.Team.Name)), e.Minute, outcomeToEmoji(e.Outcome)))

	b.WriteString(fmt.Sprintf("Chance for %s", home(e.Player.Name.String())))
	if e.Assister != nil {
		b.WriteString(fmt.Sprintf(", assisted by %s", home(e.Assister.Name.String())))
	}
	b.WriteString(":")

	if e.Tackler != nil {
		away := color.New(color.FgCyan, color.Bold).SprintFunc()
		b.WriteString(fmt.Sprintf(" tackled by %s!", away(e.Tackler.Name.String())))
		return b.String()
	}

	switch e.Outcome {
	case "OFFTARGET":
		b.WriteString(" the shot is off-target!")
	case "SAVE":
		b.WriteString(" the shot is saved!")
	}

	return b.String()
}

type BookingEvent struct {
	Event
	Outcome string
}

func NewBookingEvent(minute int, team Team, player TeamPlayer) BookingEvent {
	return BookingEvent{
		Event: NewEvent(minute, team, player),
	}
}

func (e *BookingEvent) WithOutcome(o string) {
	e.Outcome = o
}

func (e BookingEvent) String() string {
	teamColor := color.New(colorToBg(e.Team.Colors[0]), colorToFg(e.Team.Colors[1]), color.Bold).SprintFunc()
	home := color.New(color.FgBlue, color.Bold).SprintFunc()

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s [%2dm] %c ", teamColor(strings.ToUpper(e.Team.Name)), e.Minute, outcomeToEmoji(e.Outcome)))

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

// TODO: use a map instead
func outcomeToEmoji(o string) rune {
	switch o {
	case "GOAL":
		return SOCCER_BALL
	case "YELLOW":
		return YELLOW_CARD
	case "SECONDYELLOW":
		fallthrough
	case "RED":
		return RED_CARD
	default:
		return ' '
	}
}

// TODO: use a map instead
func colorToFg(c string) color.Attribute {
	switch c {
	case "RED":
		return color.FgRed
	case "GREEN":
		return color.FgGreen
	case "WHITE":
		return color.FgWhite
	case "BLUE":
		return color.FgBlue
	case "BLACK":
		fallthrough
	default:
		return color.FgBlack
	}
}

// TODO: use a map instead
func colorToBg(c string) color.Attribute {
	switch c {
	case "RED":
		return color.BgRed
	case "GREEN":
		return color.BgGreen
	case "BLACK":
		return color.BgBlack
	case "BLUE":
		return color.BgBlue
	case "WHITE":
		fallthrough
	default:
		return color.BgWhite
	}
}
