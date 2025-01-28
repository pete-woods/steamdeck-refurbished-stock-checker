package fetch

import (
	"slices"
)

type cleaner struct {
	list []func() error
}

func (l *cleaner) Add(f func() error) {
	l.list = append(l.list, f)
}

func (l *cleaner) Cleanup() error {
	for _, f := range slices.Backward(l.list) {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}
