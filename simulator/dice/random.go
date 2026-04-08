package dice

type Random interface {
	Intn(n int) int
}
