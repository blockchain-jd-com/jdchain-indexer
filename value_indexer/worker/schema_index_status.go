package worker

import (
	"encoding/json"
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"github.com/elliotchance/pie/pie"
	"math/rand"
	"strconv"
)

type SchemaStatus int

func (status SchemaStatus) String() string {
	switch status {
	case SchemaStatusDefault:
		return "SchemaStatusDefault"
	case SchemaStatusRunning:
		return "SchemaStatusRunning"
	case SchemaStatusClearing:
		return "SchemaStatusClearing"
	case SchemaStatusCleared:
		return "SchemaStatusCleared"
	case SchemaStatusStopped:
		return "SchemaStatusStopped"
	default:
		return "unknown status"
	}
}

const (
	SchemaStatusDefault  SchemaStatus = 0 //schema just exists, but no index data
	SchemaStatusRunning  SchemaStatus = 1 // schema is indexing data
	SchemaStatusStopped  SchemaStatus = 2 // schema is indexing data
	SchemaStatusClearing SchemaStatus = 3 // schema is clearing indexed data
	SchemaStatusCleared  SchemaStatus = 4 // schema's indexed data is clear
)

func NewSchemaIndexStatusCleared(info SchemaInfo, uid string) SchemaIndexStatus {
	return NewSchemaIndexStatus(uid, info, SchemaStatusCleared, -1)
}

func NewSchemaIndexStatusClearing(info SchemaInfo, uid string) SchemaIndexStatus {
	return NewSchemaIndexStatus(uid, info, SchemaStatusClearing, -1)
}

func NewSchemaIndexStatusRunning(info SchemaInfo, uid string) SchemaIndexStatus {
	return NewSchemaIndexStatus(uid, info, SchemaStatusRunning, -1)
}

func NewSchemaIndexStatusDefault(info SchemaInfo, uid string) SchemaIndexStatus {
	return NewSchemaIndexStatus(uid, info, SchemaStatusDefault, -1)
}

func NewSchemaIndexStatusStopped(info SchemaInfo, uid string) SchemaIndexStatus {
	return NewSchemaIndexStatus(uid, info, SchemaStatusStopped, -1)
}

func NewSchemaIndexStatus(uid string, info SchemaInfo, status SchemaStatus, progress int64) SchemaIndexStatus {
	return SchemaIndexStatus{
		uid:      uid,
		Schema:   info,
		Status:   status,
		from:     0,
		to:       -1,
		Progress: progress,
	}
}

type SchemaIndexStatus struct {
	uid      string
	Schema   SchemaInfo   `json:"schema"`
	Status   SchemaStatus `json:"status"`
	Progress int64
	from, to int
}

func (status SchemaIndexStatus) MetaSchemes() (schemes dgraph_helper.Schemas) {
	schemes = schemes.Add(dgraph_helper.NewSchemaStringTermIndex(status.PredictTo("id")))
	return
}

func (status SchemaIndexStatus) IsRunning() bool {
	return status.Status == SchemaStatusRunning
}

func (status SchemaIndexStatus) UniqueMutationName() string {
	return fmt.Sprintf("schemainfo-%s", status.Schema.ID)
}

func (status SchemaIndexStatus) PredictTo(name string) string {
	return fmt.Sprintf("%s-%s", "schemainfo", name)
}

func (status SchemaIndexStatus) ProgressMutations(progress int64) (mutations dgraph_helper.Mutations) {
	if len(status.uid) <= 0 {
		return
	}
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(status.uid),
			dgraph_helper.MutationItemValue(strconv.FormatInt(progress, 10)),
			dgraph_helper.MutationPredict(status.PredictTo("progress")),
		),
	)
	return
}

func (status SchemaIndexStatus) RangeMutations() (mutations dgraph_helper.Mutations) {
	if len(status.uid) <= 0 {
		return
	}
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(status.uid),
			dgraph_helper.MutationItemValue(strconv.Itoa(status.from)),
			dgraph_helper.MutationPredict(status.PredictTo("from")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(status.uid),
			dgraph_helper.MutationItemValue(strconv.Itoa(status.to)),
			dgraph_helper.MutationPredict(status.PredictTo("to")),
		),
	)
	return
}

func (status SchemaIndexStatus) UpdateStatusMutations() (mutations dgraph_helper.Mutations) {
	if len(status.uid) <= 0 {
		return
	}
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(status.uid),
			dgraph_helper.MutationItemValue(strconv.Itoa(int(status.Status))),
			dgraph_helper.MutationPredict(status.PredictTo("status")),
		),
	)
	return
}

func (status SchemaIndexStatus) UpdateContentMutations() (mutations dgraph_helper.Mutations) {
	if len(status.uid) <= 0 {
		return
	}
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(status.uid),
			dgraph_helper.MutationItemValue(status.Schema.Content),
			dgraph_helper.MutationPredict(status.PredictTo("content")),
		),
	)
	return
}

func (status SchemaIndexStatus) DeleteMutations() (mutations dgraph_helper.Mutations) {
	if len(status.uid) <= 0 {
		return
	}
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(status.uid),
			dgraph_helper.MutationItemValue("*"),
			dgraph_helper.MutationPredict("*"),
		),
	)
	return
}

func (status SchemaIndexStatus) CreateMutations() (mutations dgraph_helper.Mutations) {
	mutations = status.Schema.CreateMutations().
		Add(
			dgraph_helper.NewMutation(
				dgraph_helper.MutationItemEmpty(status.UniqueMutationName()),
				dgraph_helper.MutationItemValue(strconv.Itoa(int(status.Status))),
				dgraph_helper.MutationPredict(status.PredictTo("status")),
			),
			dgraph_helper.NewMutation(
				dgraph_helper.MutationItemEmpty(status.UniqueMutationName()),
				dgraph_helper.MutationItemValue(strconv.FormatInt(status.Progress, 10)),
				dgraph_helper.MutationPredict(status.PredictTo("progress")),
			),
		)
	return
}

//go:generate  pie SchemaIndexStatusList.*
type SchemaIndexStatusList []SchemaIndexStatus

// All will return true if all callbacks return true. It follows the same logic
// as the all() function in Python.
//
// If the list is empty then true is always returned.
func (ss SchemaIndexStatusList) All(fn func(value SchemaIndexStatus) bool) bool {
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
func (ss SchemaIndexStatusList) Any(fn func(value SchemaIndexStatus) bool) bool {
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
func (ss SchemaIndexStatusList) Append(elements ...SchemaIndexStatus) SchemaIndexStatusList {
	return append(ss, elements...)
}

// Bottom will return n elements from bottom
//
// that means that elements is taken from the end of the slice
// for this [1,2,3] slice with n == 2 will be returned [3,2]
// if the slice has less elements then n that'll return all elements
// if n < 0 it'll return empty slice.
func (ss SchemaIndexStatusList) Bottom(n int) (top SchemaIndexStatusList) {
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
func (ss SchemaIndexStatusList) Contains(lookingFor SchemaIndexStatus) bool {
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
func (ss SchemaIndexStatusList) Each(fn func(SchemaIndexStatus)) SchemaIndexStatusList {
	for _, s := range ss {
		fn(s)
	}

	return ss
}

// Extend will return a new slice with the slices of elements appended to the
// end.
//
// It is acceptable to provide zero arguments.
func (ss SchemaIndexStatusList) Extend(slices ...SchemaIndexStatusList) (ss2 SchemaIndexStatusList) {
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
func (ss SchemaIndexStatusList) Filter(condition func(SchemaIndexStatus) bool) (ss2 SchemaIndexStatusList) {
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
func (ss SchemaIndexStatusList) FilterNot(condition func(SchemaIndexStatus) bool) (ss2 SchemaIndexStatusList) {
	for _, s := range ss {
		if !condition(s) {
			ss2 = append(ss2, s)
		}
	}

	return
}

// First returns the first element, or zero. Also see FirstOr().
func (ss SchemaIndexStatusList) First() SchemaIndexStatus {
	return ss.FirstOr(SchemaIndexStatus{})
}

// FirstOr returns the first element or a default value if there are no
// elements.
func (ss SchemaIndexStatusList) FirstOr(defaultValue SchemaIndexStatus) SchemaIndexStatus {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[0]
}

// JSONString returns the JSON encoded array as a string.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array.
func (ss SchemaIndexStatusList) JSONString() string {
	if ss == nil {
		return "[]"
	}

	// An error should not be possible.
	data, _ := json.Marshal(ss)

	return string(data)
}

// Last returns the last element, or zero. Also see LastOr().
func (ss SchemaIndexStatusList) Last() SchemaIndexStatus {
	return ss.LastOr(SchemaIndexStatus{})
}

// LastOr returns the last element or a default value if there are no elements.
func (ss SchemaIndexStatusList) LastOr(defaultValue SchemaIndexStatus) SchemaIndexStatus {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[len(ss)-1]
}

// Len returns the number of elements.
func (ss SchemaIndexStatusList) Len() int {
	return len(ss)
}

// Map will return a new slice where each element has been mapped (transformed).
// The number of elements returned will always be the same as the input.
//
// Be careful when using this with slices of pointers. If you modify the input
// value it will affect the original slice. Be sure to return a new allocated
// object or deep copy the existing one.
func (ss SchemaIndexStatusList) Map(fn func(SchemaIndexStatus) SchemaIndexStatus) (ss2 SchemaIndexStatusList) {
	if ss == nil {
		return nil
	}

	ss2 = make([]SchemaIndexStatus, len(ss))
	for i, s := range ss {
		ss2[i] = fn(s)
	}

	return
}

// Random returns a random element by your rand.Source, or zero
func (ss SchemaIndexStatusList) Random(source rand.Source) SchemaIndexStatus {
	n := len(ss)

	// Avoid the extra allocation.
	if n < 1 {
		return SchemaIndexStatus{}
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
func (ss SchemaIndexStatusList) Reverse() SchemaIndexStatusList {
	// Avoid the allocation. If there is one element or less it is already
	// reversed.
	if len(ss) < 2 {
		return ss
	}

	sorted := make([]SchemaIndexStatus, len(ss))
	for i := 0; i < len(ss); i++ {
		sorted[i] = ss[len(ss)-i-1]
	}

	return sorted
}

// Top will return n elements from head of the slice
// if the slice has less elements then n that'll return all elements
// if n < 0 it'll return empty slice.
func (ss SchemaIndexStatusList) Top(n int) (top SchemaIndexStatusList) {
	for i := 0; i < len(ss) && n > 0; i++ {
		top = append(top, ss[i])
		n--
	}

	return
}

// ToStrings transforms each element to a string.
func (ss SchemaIndexStatusList) ToStrings(transform func(SchemaIndexStatus) string) pie.Strings {
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
