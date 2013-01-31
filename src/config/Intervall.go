package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Intervall time.Duration

func (t *Intervall) Set(value string) error {
	var intervallF float64
	var err error

	unit := strings.TrimLeft(value, "WDHMSwdhms")
	intervallS := strings.TrimRight(value, ".,0123456789")
	if len(intervallS)+len(unit) < len(value) {
		return errors.New("supported Units are: [W]eeks, [D]ays, [H]ours, [M]inutes, [S]econds (only one unit allowed")
	}
	if len(intervallS) == 0 {
		intervallF = 1
	} else {
		intervallS = strings.Replace(intervallS, ",", ".", 1)
		intervallF, err = strconv.ParseFloat(intervallS, 32)
		if err != nil {
			return err
		}
	}
	if len(unit) == 0 {
		unit = "M"
	}
	if len(unit) > 0 {
		return errors.New("supported Units are: [W]eeks, [D]ays, [H]ours, [M]inutes, [S]econds")
	}
	switch strings.ToUpper(unit) {
	case "W":
		intervallF *= 7
		unit = "D"
		fallthrough
	case "D":
		intervallF *= 24
		unit = "H"
		fallthrough
	case "H":
		*t = Intervall(int64(intervallF * float64(time.Hour)))
	case "M":
		*t = Intervall(int64(intervallF * float64(time.Minute)))
	case "S":
		*t = Intervall(int64(intervallF * float64(time.Second)))
	}
	return nil
}
func (t *Intervall) String() string {
	return fmt.Sprint(*t)
}
func (t Intervall) After() chang Time {
	return time.After(time.Duration(t))
}
func DefaultIntervall() Intervall {
	return Intervall(5 * time.Minute)
}
