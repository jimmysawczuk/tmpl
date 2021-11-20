package tmplfunc

type Watcher interface {
	Watch(string) error
}
