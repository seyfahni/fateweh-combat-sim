package log

type Group []Simulation

func (g Group) PrintTo(p Printer) error {
	for _, log := range g {
		err := log.PrintTo(p)
		if err != nil {
			return err
		}
	}
	return nil
}
