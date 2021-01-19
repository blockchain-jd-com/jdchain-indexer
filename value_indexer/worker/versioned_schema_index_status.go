package worker

import (
	"context"
	"encoding/json"
	"github.com/elliotchance/pie/pie"
	"github.com/elliotchance/pie/pie/util"
	"math/rand"
)

const (
	SchemaIndexStatusVersionRemoved = -1
	SchemaIndexStatusVersionDefault = 0
	SchemaIndexStatusVersionUpdated = 1
)

type VersionedSchemaIndexStatus struct {
	schemaIndexStatus SchemaIndexStatus
	version           int
}

//go:generate  pie VersionedSchemaIndexStatusList.*
type VersionedSchemaIndexStatusList []VersionedSchemaIndexStatus

// All will return true if all callbacks return true. It follows the same logic
// as the all() function in Python.
//
// If the list is empty then true is always returned.
func (ss VersionedSchemaIndexStatusList) All(fn func(value VersionedSchemaIndexStatus) bool) bool {
	for _, value := range ss {
		if !fn(value) {
			return false
		}
	}

	return true
}

// Any will return true if any callbacks return true. It follows the same logic
// as the any() function in Python.
//
// If the list is empty then false is always returned.
func (ss VersionedSchemaIndexStatusList) Any(fn func(value VersionedSchemaIndexStatus) bool) bool {
	for _, value := range ss {
		if fn(value) {
			return true
		}
	}

	return false
}

// Append will return a new slice with the elements appended to the end. It is a
// wrapper for the internal append(). It is offered as a function so that it can
// more easily chained.
//
// It is acceptable to provide zero arguments.
func (ss VersionedSchemaIndexStatusList) Append(elements ...VersionedSchemaIndexStatus) VersionedSchemaIndexStatusList {
	return append(ss, elements...)
}

// Bottom will return n elements from bottom
//
// that means that elements is taken from the end of the slice
// for this [1,2,3] slice with n == 2 will be returned [3,2]
// if the slice has less elements then n that'll return all elements
// if n < 0 it'll return empty slice.
func (ss VersionedSchemaIndexStatusList) Bottom(n int) (top VersionedSchemaIndexStatusList) {
	var lastIndex = len(ss) - 1
	for i := lastIndex; i > -1 && n > 0; i-- {
		top = append(top, ss[i])
		n--
	}

	return
}

// Contains returns true if the element exists in the slice.
//
// When using slices of pointers it will only compare by address, not value.
func (ss VersionedSchemaIndexStatusList) Contains(lookingFor VersionedSchemaIndexStatus) bool {
	for _, s := range ss {
		if s == lookingFor {
			return true
		}
	}

	return false
}

// Each is more condensed version of Transform that allows an action to happen
// on each elements and pass the original slice on.
//
//   cars.Each(func (car *Car) {
//       fmt.Printf("Car color is: %s\n", car.Color)
//   })
//
// Pie will not ensure immutability on items passed in so they can be
// manipulated, if you choose to do it this way, for example:
//
//   // Set all car colors to Red.
//   cars.Each(func (car *Car) {
//       car.Color = "Red"
//   })
//
func (ss VersionedSchemaIndexStatusList) Each(fn func(VersionedSchemaIndexStatus)) VersionedSchemaIndexStatusList {
	for _, s := range ss {
		fn(s)
	}

	return ss
}

// Extend will return a new slice with the slices of elements appended to the
// end.
//
// It is acceptable to provide zero arguments.
func (ss VersionedSchemaIndexStatusList) Extend(slices ...VersionedSchemaIndexStatusList) (ss2 VersionedSchemaIndexStatusList) {
	ss2 = ss

	for _, slice := range slices {
		ss2 = ss2.Append(slice...)
	}

	return ss2
}

// Filter will return a new slice containing only the elements that return
// true from the condition. The returned slice may contain zero elements (nil).
//
// FilterNot works in the opposite way of Filter.
func (ss VersionedSchemaIndexStatusList) Filter(condition func(VersionedSchemaIndexStatus) bool) (ss2 VersionedSchemaIndexStatusList) {
	for _, s := range ss {
		if condition(s) {
			ss2 = append(ss2, s)
		}
	}
	return
}

// FilterNot works the same as Filter, with a negated condition. That is, it will
// return a new slice only containing the elements that returned false from the
// condition. The returned slice may contain zero elements (nil).
func (ss VersionedSchemaIndexStatusList) FilterNot(condition func(VersionedSchemaIndexStatus) bool) (ss2 VersionedSchemaIndexStatusList) {
	for _, s := range ss {
		if !condition(s) {
			ss2 = append(ss2, s)
		}
	}

	return
}

// First returns the first element, or zero. Also see FirstOr().
func (ss VersionedSchemaIndexStatusList) First() VersionedSchemaIndexStatus {
	return ss.FirstOr(VersionedSchemaIndexStatus{})
}

// FirstOr returns the first element or a default value if there are no
// elements.
func (ss VersionedSchemaIndexStatusList) FirstOr(defaultValue VersionedSchemaIndexStatus) VersionedSchemaIndexStatus {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[0]
}

// JSONString returns the JSON encoded array as a string.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array.
func (ss VersionedSchemaIndexStatusList) JSONString() string {
	if ss == nil {
		return "[]"
	}

	// An error should not be possible.
	data, _ := json.Marshal(ss)

	return string(data)
}

// Last returns the last element, or zero. Also see LastOr().
func (ss VersionedSchemaIndexStatusList) Last() VersionedSchemaIndexStatus {
	return ss.LastOr(VersionedSchemaIndexStatus{})
}

// LastOr returns the last element or a default value if there are no elements.
func (ss VersionedSchemaIndexStatusList) LastOr(defaultValue VersionedSchemaIndexStatus) VersionedSchemaIndexStatus {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[len(ss)-1]
}

// Len returns the number of elements.
func (ss VersionedSchemaIndexStatusList) Len() int {
	return len(ss)
}

// Map will return a new slice where each element has been mapped (transformed).
// The number of elements returned will always be the same as the input.
//
// Be careful when using this with slices of pointers. If you modify the input
// value it will affect the original slice. Be sure to return a new allocated
// object or deep copy the existing one.
func (ss VersionedSchemaIndexStatusList) Map(fn func(VersionedSchemaIndexStatus) VersionedSchemaIndexStatus) (ss2 VersionedSchemaIndexStatusList) {
	if ss == nil {
		return nil
	}

	ss2 = make([]VersionedSchemaIndexStatus, len(ss))
	for i, s := range ss {
		ss2[i] = fn(s)
	}

	return
}

// Random returns a random element by your rand.Source, or zero
func (ss VersionedSchemaIndexStatusList) Random(source rand.Source) VersionedSchemaIndexStatus {
	n := len(ss)

	// Avoid the extra allocation.
	if n < 1 {
		return VersionedSchemaIndexStatus{}
	}
	if n < 2 {
		return ss[0]
	}
	rnd := rand.New(source)
	i := rnd.Intn(n)
	return ss[i]
}

// Reverse returns a new copy of the slice with the elements ordered in reverse.
// This is useful when combined with Sort to get a descending sort order:
//
//   ss.Sort().Reverse()
//
func (ss VersionedSchemaIndexStatusList) Reverse() VersionedSchemaIndexStatusList {
	// Avoid the allocation. If there is one element or less it is already
	// reversed.
	if len(ss) < 2 {
		return ss
	}

	sorted := make([]VersionedSchemaIndexStatus, len(ss))
	for i := 0; i < len(ss); i++ {
		sorted[i] = ss[len(ss)-i-1]
	}

	return sorted
}

// Send sends elements to channel
// in normal act it sends all elements but if func canceled it can be less
//
// it locks execution of gorutine
// it doesn't close channel after work
// returns sended elements if len(this) != len(old) considered func was canceled
func (ss VersionedSchemaIndexStatusList) Send(ctx context.Context, ch chan<- VersionedSchemaIndexStatus) VersionedSchemaIndexStatusList {
	for i, s := range ss {
		select {
		case <-ctx.Done():
			return ss[:i]
		default:
			ch <- s
		}
	}

	return ss
}

// Shuffle returns shuffled slice by your rand.Source
func (ss VersionedSchemaIndexStatusList) Shuffle(source rand.Source) VersionedSchemaIndexStatusList {
	n := len(ss)

	// Avoid the extra allocation.
	if n < 2 {
		return ss
	}

	// go 1.10+ provides rnd.Shuffle. However, to support older versions we copy
	// the algorithm directly from the go source: src/math/rand/rand.go below,
	// with some adjustments:
	shuffled := make([]VersionedSchemaIndexStatus, n)
	copy(shuffled, ss)

	rnd := rand.New(source)

	util.Shuffle(rnd, n, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

// Top will return n elements from head of the slice
// if the slice has less elements then n that'll return all elements
// if n < 0 it'll return empty slice.
func (ss VersionedSchemaIndexStatusList) Top(n int) (top VersionedSchemaIndexStatusList) {
	for i := 0; i < len(ss) && n > 0; i++ {
		top = append(top, ss[i])
		n--
	}

	return
}

// ToStrings transforms each element to a string.
func (ss VersionedSchemaIndexStatusList) ToStrings(transform func(VersionedSchemaIndexStatus) string) pie.Strings {
	l := len(ss)

	// Avoid the allocation.
	if l == 0 {
		return nil
	}

	result := make(pie.Strings, l)
	for i := 0; i < l; i++ {
		result[i] = transform(ss[i])
	}

	return result
}
