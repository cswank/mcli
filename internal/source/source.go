package source

type Result struct {
	Value string
	URL   string
}

type Source interface {
	FindArtist(string) ([]Result, error)
	FindAlbum(string) ([]Result, error)
	FindTrack(string) ([]Result, error)
}
