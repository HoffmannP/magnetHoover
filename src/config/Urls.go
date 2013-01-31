package config

import (
	"fmt"
)

type Urls []string

func (u *Urls) Set(value string) error {
	*u = append(*u, value)
	return nil
}
func (u *Urls) String() string {
	return fmt.Sprint(*u)
}
