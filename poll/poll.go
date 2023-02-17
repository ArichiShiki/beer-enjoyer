package poll

import (
	"fmt"
	"time"

	"arichi/beerEnjoyer/timeHelper"
)

type Poll struct {
	Type     PollType
	Question string
	Options  []string
}

type PollType int

const (
	CurWeekPoll PollType = iota
	NextWeekPoll
)

const (
	curWeekQuestion  string = "Могу играть на этой неделе в:"
	nextWeekQuestion string = "Могу играть на следующей неделе в:"
	neverOption      string = "не могу в эти даты"
)

func getOptionsFromDateTillEndOfWeek(from time.Time) []string {
	cur := from

	weekday := timeHelper.Weekday(cur)

	options := make([]string, 0, timeHelper.DaysInWeek-weekday+1)

	for i := weekday; i <= timeHelper.DaysInWeek; i++ {
		s := fmt.Sprintf("%d %s, %s", cur.Day(), timeHelper.MonthAsGenitive[cur.Month()], timeHelper.DayOfWeek[i])
		options = append(options, s)
		cur = cur.Add(timeHelper.Day)
	}

	options = append(options, neverOption)

	return options
}

func CreatePoll(pollType PollType) Poll {
	var options []string
	var question string

	switch pollType {
	case CurWeekPoll:
		question = curWeekQuestion
		options = getOptionsFromDateTillEndOfWeek(time.Now())
	case NextWeekPoll:
		question = nextWeekQuestion
		options = getOptionsFromDateTillEndOfWeek(timeHelper.NextMonday(time.Now()))
	}

	return Poll{
		Type:     pollType,
		Question: question,
		Options:  options,
	}
}
