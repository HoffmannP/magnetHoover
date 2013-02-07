package history

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const torrents_table = "Torrents"
const torrent_buffer = 15

type History struct {
	db *sql.DB
	q  chan string
}

func New(n string) (*History, error) {
	db, err := sql.Open("sqlite3", n)
	if err != nil {
		fmt.Println("SQL Database access error »%s«: %v\n", n, err)
		return nil, err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS `" + torrents_table + "` (url string not null primary key)"); err != nil {
		fmt.Println("SQL Create Error »%s«: %v\n", torrents_table, err)
		return nil, err
	}

	h := &History{db, make(chan string)}
	go h.adder()
	return h, nil
}

func (h *History) Exists(url string) bool {
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM `" + torrents_table + "` WHERE url = \"" + url + "\"").Scan(&count)
	if err != nil {
		fmt.Printf("SQL Count Error »%s«: %v\n", url, err)
	}
	return count > 0
}

func (h *History) adder() {
	for url := range h.q {
		_, err := h.db.Exec("INSERT INTO `" + torrents_table + "` VALUES (\"" + url + "\")")
		if err != nil {
			fmt.Println("SQL Insert Error »%s«: %v\n", url, err)
		}
	}
}

func (h *History) Add(url string) {
	h.q <- url
}

func (h *History) Close() {
	close(h.q)
	h.db.Close()
}
