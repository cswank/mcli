package player

type History interface {
	Save(Result) error
	Fetch(int, int, Sort) (*Results, error)
}
