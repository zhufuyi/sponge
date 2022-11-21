package window

import (
	"testing"
)

func TestSum(t *testing.T) {
	it := Iterator{
		count:         0,
		iteratedCount: 0,
		cur:           &Bucket{},
	}

	t.Log(Sum(it))
	t.Log(Avg(it))
	t.Log(Max(it))
	t.Log(Count(it))
}
