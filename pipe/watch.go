package pipe

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

type Watcher struct {
	active bool

	w *fsnotify.Watcher

	pipes map[string]*Pipe
	refs  map[string]*Pipe
}

func New(active bool) (*Watcher, error) {
	if active {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return nil, errors.Wrap(err, "fsnotify: new watcher")
		}

		return &Watcher{
			active: true,
			w:      watcher,
			pipes:  map[string]*Pipe{},
			refs:   map[string]*Pipe{},
		}, nil
	}

	return &Watcher{}, nil
}

func (w *Watcher) Close() error {
	if !w.active {
		return nil
	}

	return w.w.Close()
}

func (w *Watcher) AddPipe(p *Pipe) error {
	if !w.active {
		return nil
	}

	path, err := filepath.Abs(p.In)
	if err != nil {
		return errors.Errorf("filepath: abs (path: %s)", p.In)
	}

	// TODO: check for existence here too?

	w.w.Add(path)

	if _, ok := w.pipes[path]; !ok {
		w.pipes[p.In] = p
	}

	return nil
}

func (w *Watcher) AddRef(ref string, pipe *Pipe) error {
	if !w.active {
		return nil
	}

	path, err := filepath.Abs(ref)
	if err != nil {
		return errors.Errorf("filepath: abs (path: %s)", ref)
	}

	w.w.Add(path)
	w.refs[path] = pipe
	return nil
}

func (w *Watcher) Watch(notify chan string) {
	if !w.active {
		return
	}

	for {
		select {
		case event, ok := <-w.w.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("changed:", event.Name, event.Op)

				for path, pipe := range w.pipes {
					if path == event.Name {
						if err := pipe.Run(); err != nil {
							log.Printf("%s", errors.Wrapf(err, "pipeline (path: %s)", pipe.In))
						}
						log.Println(" --> wrote:", pipe.Out)

						pipe.AttachRefs(w)

						continue
					}
				}

				for ref, pipe := range w.refs {
					if ref == event.Name {
						if err := pipe.Run(); err != nil {
							log.Printf("%s", errors.Wrapf(err, "pipeline (path: %s)", pipe.In))
						}
						log.Println(" --> wrote:", pipe.Out)

						pipe.AttachRefs(w)
					}
				}

				for {
					select {
					case notify <- event.Name:
						log.Println(" --> ! notified listener")
						continue
					default:
						log.Println("nothing to notify!")
					}

					break
				}
			}

		case err, ok := <-w.w.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
