/*

mybot - Illustrative Slack bot in Go

Copyright (c) 2015 RapidLoop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	//"encoding/csv"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: mybot slack-bot-token aemet-API-Key\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	fmt.Println("mybot ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) == 3 && parts[1] == "temperatura" {
				// looks good, get the quote and reply with the result
				go func(m Message) {
					//m.Text = getQuote(parts[2])
					url := fmt.Sprintf("https://opendata.aemet.es/opendata/api/observacion/convencional/datos/estacion/B228/?api_key=%v",os.Args[2])
					//m.Text = recuperaDatos("https://opendata.aemet.es/opendata/api/observacion/convencional/datos/estacion/B228/?api_key=eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzcGJyYXZvQGdtYWlsLmNvbSIsImp0aSI6IjgyNGM1Y2UzLWFmNzYtNDk0NS1hZDBmLTdhMDk1ZTIyMzJkZCIsImlzcyI6IkFFTUVUIiwiaWF0IjoxNTA5NTU2OTY5LCJ1c2VySWQiOiI4MjRjNWNlMy1hZjc2LTQ5NDUtYWQwZi03YTA5NWUyMjMyZGQiLCJyb2xlIjoiIn0.XAEUT7p_9sXrkavMunL9CtQySwUWOicCIbGfsYxdVZk")
					m.Text = recuperaDatos(url)
					postMessage(ws, m)
				}(m)
				// NOTE: the Message object is copied, this is intentional
			} else {
				// huh?
				m.Text = fmt.Sprintf("sorry, that does not compute\n")
				postMessage(ws, m)
			}
		}
	}
}

// Get the quote via Yahoo. You should replace this method to something
// relevant to your team!
func getQuote(sym string) string {
	sym = strings.ToUpper(sym)
	//url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=nsl1op&e=.csv", sym)
	url := fmt.Sprintf("https://reloj-alarma.es/temporizador/ano-nuevo/?s0%s", sym)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	//rows, err := csv.NewReader(resp.Body).ReadAll()
	
	rows, err := ioutil.ReadAll(resp.Body)
	if err != nil { 
		return fmt.Sprintf("error: %v", err)
	}
	//if len(rows) >= 1 && len(rows[0]) == 5 {
		if len(rows) >= 1 {
		return fmt.Sprintf("La temperatura es %s", rows)
	}
	return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}
