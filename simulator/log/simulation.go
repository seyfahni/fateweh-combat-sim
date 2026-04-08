package log

type Simulation interface {
	PrintTo(Printer) error
}
