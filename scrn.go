package scrn

import (
        "fmt"
        "os"
)

type Screen struct {}

func (s *Screen) Save() {
        os.Stdout.WriteString("\033[?1049h")
}

func (s *Screen) Restore() {
        os.Stdout.WriteString("\033[?1049l")
}

func (s *Screen) Clear() {
        os.Stdout.WriteString("\033[2J")
}

func (s *Screen) Move(col, row int) {
        os.Stdout.WriteString(fmt.Sprintf("\033[%d;%dH", row, col))
}

func (s *Screen) Print(text string) (n int, err error) {
        return os.Stdout.WriteString(text)
}

func (s *Screen) Printf(format string, a ...interface{}) (n int, err error) {
        text := fmt.Sprintf(format, a...)
        return s.Print(text)
}
