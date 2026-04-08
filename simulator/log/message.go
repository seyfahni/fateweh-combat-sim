package log

import "fmt"

type Message string

func (m Message) PrintTo(p Printer) error {
	return p.Print(string(m))
}

func MessageF(format string, args ...any) Message {
	return Message(fmt.Sprintf(format, args...))
}
