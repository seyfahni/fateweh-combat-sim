package log

type Indent struct {
	Log Simulation
}

func (i Indent) PrintTo(printer Printer) error {
	return i.Log.PrintTo(indenter{printer})
}

type indenter struct {
	printer Printer
}

func (i indenter) Print(line string) error {
	return i.printer.Print("  " + line)
}
