package scrn

import (
	"fmt"
	"os"
	"sync"
	"unicode"

	"golang.org/x/sys/unix"
)

type Screen struct{}

func (s *Screen) Open() {
	os.Stdout.WriteString("\033[?1049h")
}

func (s *Screen) Close() {
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

func GetSize() (width, height int, err error) {
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		return -1, -1, err
	}
	return int(ws.Col), int(ws.Row), nil
}

type KeyCode int

const (
	KeyESC = KeyCode(int('\033'))
)

var (
	KeyPressCaseSensitive bool
	watchKeyPressOnce     sync.Once
	c                     chan KeyCode
	mask                  [2]uint64
)

func KeyPress(code ...KeyCode) <-chan KeyCode {
	watchKeyPressOnce.Do(func() {
		c = make(chan KeyCode, 1)
		termios, err := unix.IoctlGetTermios(0, unix.TCGETS)
		if err != nil {
			panic(err)
		}
		termios.Lflag &^= unix.ECHO | unix.ICANON
		termios.Lflag |= unix.ISIG
		termios.Cc[unix.VMIN] = 1
		termios.Cc[unix.VTIME] = 0
		if err := unix.IoctlSetTermios(0, unix.TCSETS, termios); err != nil {
			panic(err)
		}
		go func() {
			var buf [1]byte
			for {
				if n, _ := unix.Read(0, buf[:]); n == 0 {
					continue
				}
				if contains(int(buf[0])) {
					c <- KeyCode(buf[0])
				}
			}
		}()
	})
	if len(code) == 0 {
		var max uint64 = 1<<64 - 1
		for i := range mask {
			mask[i] = max
		}
	} else {
		reset()
		for _, c := range code {
			add(c)
		}
	}
	return c
}

func add(k KeyCode) {
	if uint32(k) > unicode.MaxASCII {
		return
	}
	set(int(k))
	r := rune(k)
	if KeyPressCaseSensitive && !unicode.IsLetter(r) {
		return
	}
	if unicode.IsLower(r) {
		set(int(unicode.ToUpper(r)))
		return
	}
	if unicode.IsUpper(r) {
		set(int(unicode.ToLower(r)))
		return
	}
}

func contains(code int) bool {
	return (mask[code/64]>>uint(code&63))&1 == 1
}

func set(code int) {
	mask[code/64] |= 1 << uint(code&63)
}

func reset() {
	var m [2]uint64
	mask = m
}
