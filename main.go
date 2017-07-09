package main

import (
	"fmt"
	"syscall"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func enterAlternateScreen(fd int) {
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x3f, 0x01, 0x00, 0x04, 0x09, 0x68, 0x1b, 0x5b, 0x48})
}

func main() {
	in := prompt.NewVT100Parser()
	out := prompt.NewVT100Writer()
	in.Setup()
	defer in.TearDown()
	defer fmt.Println("\nexited!")
	out.SetTitle("はろー")

	out.ScrollDown()
	out.ScrollDown()
	out.ScrollDown()
	out.ScrollDown()
	out.ScrollUp()
	out.ScrollUp()
	out.ScrollUp()
	out.ScrollUp()
	out.Flush()

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)

	buffer := prompt.NewBuffer()

	for {
		b := <-bufCh
		if ac := in.GetASCIICode(b); ac == nil {
			out.EraseDown()
			out.WriteRaw(b)
			buffer.InsertText(string(b), false, true)
		} else if ac.Key == prompt.Enter || ac.Key == prompt.ControlJ {
			buffer.InsertText("\n", false, true)
		} else if ac.Key == prompt.Left {
			l := buffer.CursorLeft(1)
			if l == 0 {
				continue
			}
			out.CursorDown(1)
			out.CursorBackward(1)
			out.EraseDown()
			out.CursorUp(1)
		} else if ac.Key == prompt.Right {
			l := buffer.CursorRight(1)
			if l == 0 {
				continue
			}
			out.CursorDown(1)
			out.EraseDown()
			out.CursorForward(1)
			out.CursorUp(1)
		} else if ac.Key == prompt.Backspace {
			buffer.DeleteBeforeCursor(1)
			out.CursorBackward(1)
			out.EraseDown()
		} else if ac.Key == prompt.Right {
			buffer.CursorRight(1)
		} else if ac.Key == prompt.ControlT {
			enterAlternateScreen(syscall.Stdout)
		} else if ac.Key == prompt.ControlC {
			out.EraseDown()
			out.ClearTitle()
			out.Flush()
			return
		} else {
			out.WriteRaw(b)
			//buffer.InsertText(ac.Key.String(), false, true)
		}

		// Display completions
		if w := buffer.Document().GetWordBeforeCursor(); w != "" {
			out.SetColor("white", "teal")

			out.CursorDown(1)
			out.Write([]byte(" Foo "))
			out.SetColor("white", "darkGray")
			out.Write([]byte(" "))
			out.SetColor("white", "teal")
			out.CursorBackward(len("foo") + 3)

			out.CursorDown(1)
			out.Write([]byte(" Hello "))
			out.SetColor("white", "darkGray")
			out.Write([]byte(" "))
			out.SetColor("white", "teal")
			out.CursorBackward(len("Hello") + 3)

			out.CursorDown(1)
			out.Write([]byte(" World "))
			out.SetColor("white", "darkGray")
			out.Write([]byte(" "))
			out.SetColor("white", "teal")
			out.CursorBackward(len("World") + 3)

			out.CursorUp(3)
			out.SetColor("default", "default")
		}

		out.Flush()
	}
}

func readBuffer(bufCh chan []byte) {
	buf := make([]byte, 1024)

	for {
		if n, err := syscall.Read(syscall.Stdin, buf); err == nil {
			bufCh <- buf[:n]
		}
	}
}