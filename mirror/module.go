package mirror

type ModuleContext struct {
	Name string
}

type Module interface {
	SetInput(<-chan Request)
	Output() <-chan Request
}
