package tmplfunc

type Depender interface {
	Depend(string) error
}
