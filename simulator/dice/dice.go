package dice

type Pouch struct {
	Random
}

func (p *Pouch) RollD6() int {
	return p.Intn(6) + 1
}
