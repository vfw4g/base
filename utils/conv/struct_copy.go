package conv

import "github.com/mohae/deepcopy"

func Copy(orig interface{}) interface{} {
	return deepcopy.Copy(orig)
}
