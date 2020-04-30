package errors

var NotFound = &innerError{code: "400"}

//type notFound struct {
//	Error
//}
//
//func (e *notFound) Retry() bool {
//	return true
//}
