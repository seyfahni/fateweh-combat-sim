package log

type nothing struct{}

func (nothing) PrintTo(Printer) error {
	return nil
}

var Nothing = nothing{}
