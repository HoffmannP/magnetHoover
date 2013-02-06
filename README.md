= magnetHoover =
The purpose of this project is to provide a lighwith compiled daemon scanning torrent feeds or likewise for new torrents and adding them

== Build ==
install go on your system, change in the main folder of this project and run
	export GOPATH="$(pwd)"
    go install build
If no errors occure you should find the executabel as `bin/main`

== Config ==
Unfortunately there are no comments in json-files here are the comments for `config.json` (you find an example in `config.json.example`)
* Intervall(String): Waiting time between to polls, valid time units (required) are h(hours), m(inutes), s(econds) (TODO: larger untits like w(eek) and d(ay))
* Database(FilenameString): Sqlite3 database storing already added torrents (will be created if not existent) 
* [Transmission](http://www.transmissionbt.com/):
* * Host(String): IP or Hostname where your Transmission daemon is running
* * Port(Integer): Portnumber of your Transmission daemon is running
* * SSL(Bool): Whether SSL (https) should be used as transport communication with your Transmission daemon
* URIs(Array of Strings): A number of URLs or Identifiers possibliy prefixed with and seperated using the paragraph sign (§) by the name of a parser plugin 

== Parser plugins ==
Parser plugins have to satisfy the `ParserRawFunc` and `UrlFunc` type of the `parser` package and need to be registered using the `parser.register` function

== TODO ==
 * Larger time units
 * Test HTTPS Support
 * Transmission Authentication
 * …?

