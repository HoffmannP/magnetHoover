package history

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type History struct {
	c *sql.DB
	k string
}

func New(n string) (*History, error) {
	db, err := sql.Open("sqlite3", n)
	if err != nil {
		return nil, err
	}
	return &History{db, ""}, nil
}
func (h *History) Close() {
	h.c.Close()
}
func (h *History) Select(k string) *History {
	_, err := h.c.Exec("CREATE TABLE IF NOT EXISTS `" + k + "` (id string not null primary key)")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	h.k = k
	return h
}
func (h *History) Exists(id string) bool {
	var count int
	err := h.c.QueryRow("SELECT COUNT(*) FROM `" + h.k + "` WHERE id = \"" + id + "\"").Scan(&count)
	if err != nil {
		fmt.Printf("SQL Count Error »%s«: %v\n", id, err)
	}
	return count > 0
}
func (h *History) Add(id string) error {
	_, err := h.c.Exec("INSERT INTO `" + h.k + "` VALUES (\"" + id + "\")")
	return err
}
