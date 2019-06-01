package animation

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Animation struct {
	Frame []string
	Width int
}

func (this Animation) Progress(out io.Writer) func() {
	done := make(chan struct{})
	go func() {
		backspace := strings.Repeat("\b", this.Width)
		erase := strings.Repeat(" ", this.Width)

		fmt.Fprint(out, this.Frame[0])

		ticker := time.NewTicker(time.Second / 2)
		i := 1
		for {
			select {
			case <-done:
				ticker.Stop()
				close(done)
				fmt.Fprint(out, backspace)
				fmt.Fprint(out, erase)
				fmt.Fprint(out, backspace)
				return
			case <-ticker.C:
				if i >= len(this.Frame) {
					i = 0
				}
				fmt.Fprint(out, backspace)
				fmt.Fprint(out, this.Frame[i])
				i++
			}
		}
	}()

	return func() {
		done <- struct{}{}
	}
}

func Progress() func() {
	return Animation{
		Frame: []string{" /", " -", " \u2216", " |"},
		Width: 2,
	}.Progress(os.Stdout)
}
