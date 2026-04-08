package log

type Printer interface {
	Print(string) error
}
