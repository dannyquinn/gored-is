import (
	"fmt"
	"iter"
)

type Slice []int

func (s Slice) Map(transform func(int) int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for _,v := range s {
			if !yield(transform(v)) {
				return
			}
		}
	}
}
func main() {
	numbers := Slice{0,1,2,3,4,5}

	doubled:= numbers.Map(func (i int) int {
		return i * 2
	})

	for n:=range doubled {
		fmt.Println(n)
	}
}
