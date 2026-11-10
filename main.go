// Command bjhello is a sample hosted integration for ButlerJames: one tool,
// hello, which makes one brokered GET to https://example.com/ and says what
// came back. It speaks the hosted integration line framing on stdin and
// stdout (one JSON value per line, echoing each line's exchange id, "x")
// with no dependencies, so it builds with the Go template (an empty
// vendor/modules.txt). Its butlerjames.json allowlists example.com.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type response struct {
	Status int    `json:"status"`
	Body   []byte `json:"body"`
}

type step struct {
	Tool     string          `json:"tool"`
	Args     json.RawMessage `json:"args"`
	N        int             `json:"n"`
	Response *response       `json:"response"`
	List     bool            `json:"list"`
	X        string          `json:"x"` // the exchange id: every answer echoes it
}

var tools = json.RawMessage(`[{"name":"hello","description":"Says hello, with the title of https://example.com/ fetched through the broker.","inputSchema":{"type":"object","properties":{"name":{"type":"string"}}}}]`)

var title = regexp.MustCompile(`<title>([^<]{0,200})</title>`)

func text(s string, isError bool) map[string]any {
	return map[string]any{"result": map[string]any{"content": []map[string]string{{"type": "text", "text": s}}, "isError": isError}}
}

func answer(s step) map[string]any {
	switch {
	case s.List:
		return map[string]any{"tools": tools}
	case s.Tool != "hello":
		return text("unknown tool", true)
	case s.Response == nil:
		return map[string]any{"need": map[string]any{"method": "GET", "url": "https://example.com/"}, "state": []byte("1")}
	}
	var a struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(s.Args, &a)
	if a.Name == "" {
		a.Name = "there"
	}
	t := "no title"
	if m := title.FindSubmatch(s.Response.Body); m != nil {
		t = string(m[1])
	}
	return text(fmt.Sprintf("Hello %s, from a hosted integration. example.com answered %d: %q.", a.Name, s.Response.Status, t), false)
}

func main() {
	in := bufio.NewReaderSize(os.Stdin, 4<<20+1)
	out := json.NewEncoder(os.Stdout) // one line per value
	for {
		line, err := in.ReadSlice('\n')
		if len(line) == 0 && err != nil {
			return
		}
		var s step
		if json.Unmarshal(line, &s) != nil {
			os.Exit(2)
		}
		a := answer(s)
		a["x"] = s.X
		if out.Encode(a) != nil {
			os.Exit(2)
		}
	}
}
