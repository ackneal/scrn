package main

import (
        "time"

        "github.com/ackneal/scrn"
)

func main() {
        var s scrn.Screen
        s.Open() // open alternate screen
        defer s.Close() // close

        go func() {
                for {
                        s.Move(1,1) // moves the position of a cursor
                        s.Printf("Current Time: %s", time.Now().Format(time.RFC1123Z))
                        time.Sleep(time.Second)
                }
        }()

        // listens keyPress for q and ESC
        quitKey := []scrn.KeyCode{
                scrn.KeyCode(int('q')),
                scrn.KeyESC,
        }

        for {
                select {
                case <-scrn.KeyPress(quitKey...):
                        return
                }
        }
}
