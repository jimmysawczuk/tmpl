package tmplfunc

import (
	"github.com/pkg/errors"
)

func Seq(start, end int) ([]int, error) {
	if end <= start {
		return nil, errors.New("end can't be < start")
	}

	v := make([]int, 0, end-start)
	for i := start; i <= end; i++ {
		v = append(v, i)
	}

	return v, nil
}
