// Flag
package main

import (
	"flag"
)

var Port = flag.Int("P", 9877, "listen port")

func init() {
	flag.Parse()
}
