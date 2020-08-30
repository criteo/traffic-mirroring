package mirror

type Module interface {
	SetInput(<-chan Request)
	Output() <-chan Request
}
