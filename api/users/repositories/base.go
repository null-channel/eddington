package repositories

type Seedable interface {
	Seed() error
}
