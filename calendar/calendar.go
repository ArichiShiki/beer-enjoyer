package calendar

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CalendarInfo struct {
	Keyboard         tgbotapi.InlineKeyboardMarkup
	State            calendarKeyboardState
	Year, Day int
	Month time.Month
}

func (c *CalendarInfo) Update(callbackData string) {
	switch c.State {
	case YearKeyboard:
		switch callbackData {
		case PrevData:
			c.Year -= 9
			c.Keyboard = newPickYearCalendarKeyboard(c.Year)
		case NextData:
			c.Year += 9
			c.Keyboard = newPickYearCalendarKeyboard(c.Year)
		default:
			c.Year, _ = strconv.Atoi(callbackData)
			c.Keyboard = newPickMonthCalendarKeyboard(c.Year)
			c.State = MonthKeyboard
		}
	case MonthKeyboard:
		switch callbackData {
		case PrevData:
			c.Year--
			c.Keyboard = newPickMonthCalendarKeyboard(c.Year)
		case NextData:
			c.Year++
			c.Keyboard = newPickMonthCalendarKeyboard(c.Year)
		case GoToYearData:
			c.State = YearKeyboard
			c.Keyboard = newPickYearCalendarKeyboard(c.Year)
		default:
			m, _ := strconv.Atoi(callbackData)
			c.Month = time.Month(m)
			c.Keyboard = newPickDayCalendarKeyboard(c.Year, c.Month)
			c.State = DayKeyboard
		}
	case DayKeyboard:
		switch callbackData {
		case PrevData:
			if c.Month != 1 {
				c.Month--
			} else {
				c.Month = 12
				c.Year--
			}
			c.Keyboard = newPickDayCalendarKeyboard(c.Year, c.Month)
		case NextData:
			if c.Month != 12 {
				c.Month++
			} else {
				c.Month = 1
				c.Year++
			}
			c.Keyboard = newPickDayCalendarKeyboard(c.Year, c.Month)
		case GoToMonthData:
			c.State = MonthKeyboard
			c.Keyboard = newPickMonthCalendarKeyboard(c.Year)
		default:
			c.Day, _ = strconv.Atoi(callbackData)
			c.State = Finish
		}
	}
}

func (c *CalendarInfo) String() string {
	return fmt.Sprintf("%02d/%02d/%04d", c.Day, int(c.Month), c.Year)
}

func NewCalendar(state calendarKeyboardState) (*CalendarInfo, error) {
	res := &CalendarInfo{
		State:    state,
	}
	now := time.Now()
	switch state {
	case YearKeyboard:
		res.Keyboard = newPickYearCalendarKeyboard(now.Year())
	case MonthKeyboard:
		res.Keyboard = newPickMonthCalendarKeyboard(now.Year())
		res.Year = now.Year()
	case DayKeyboard:
		res.Keyboard = newPickDayCalendarKeyboard(now.Year(), now.Month())
		res.Year = now.Year()
		res.Month = now.Month()
	default:
		return nil, errors.New("unexpected state for new calendar")
	}
	return res, nil
}
