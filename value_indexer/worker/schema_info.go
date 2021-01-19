package worker

import (
	"encoding/json"
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/value_indexer/schema"
	"github.com/RoseRocket/xerrs"
	"github.com/elliotchance/pie/pie"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
)

func ParseSchemaInfo(raw string) (SchemaInfo, error) {
	result := gjson.Parse(raw)
	ledger := result.Get("ledger").String()
	if len(ledger) < 6 {
		return SchemaInfo{}, fmt.Errorf("ledger required")
	}

	associateAccount := result.Get("associate_account").String()
	if len(associateAccount) < 6 {
		return SchemaInfo{}, fmt.Errorf("associate_account required")
	}
	content := result.Get("content").String()
	if len(content) <= 0 {
		return SchemaInfo{}, fmt.Errorf("content required")
	}
	nodeSchema, err := schema.NewSchemaParser().FirstNodeSchema(content)
	if err != nil {
		logger.Infof("failed get node schema: %s", err)
		return SchemaInfo{}, xerrs.Mask(err, fmt.Errorf("parse Schema failed: %s", err))
	}
	_, b := nodeSchema.PrimaryField()
	if b == false {
		return SchemaInfo{}, fmt.Errorf("no primay field in Schema")
	}

	id := fmt.Sprintf("%s-%s-%s", nodeSchema.LowerName(), ledger[:6], associateAccount[:6])
	schemaInfo, err := NewSchemaInfo(id, ledger, associateAccount, content)
	if err != nil {
		return SchemaInfo{}, err
	}
	schemaInfo.nodeSchema = nodeSchema

	return schemaInfo, nil
}

func NewSchemaInfo(id, ledger, account, content string) (info SchemaInfo, e error) {
	info = SchemaInfo{
		ID:               id,
		Ledger:           ledger,
		AssociateAccount: account,
		Content:          content,
	}
	ns, err := schema.NewSchemaParser().FirstNodeSchema(content)
	if err != nil {
		e = err
		return
	}
	info.nodeSchema = ns
	return
}

type SchemaInfo struct {
	ID               string `json:"id"`
	Ledger           string `json:"ledger"`
	AssociateAccount string `json:"associate_account"`
	Content          string `json:"content"`
	nodeSchema       *schema.NodeSchema
}

func (info SchemaInfo) String() string {
	bs, err := json.Marshal(info)
	if err != nil {
		return err.Error()
	}
	return string(pretty.Pretty(bs))
}

func (info SchemaInfo) UniqueMutationName() string {
	return fmt.Sprintf("schemainfo-%s", info.ID)
}

func (info SchemaInfo) PredictTo(name string) string {
	return fmt.Sprintf("%s-%s", "schemainfo", name)
}

func (info SchemaInfo) CreateMutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(info.UniqueMutationName()),
			dgraph_helper.MutationItemValue(info.ID),
			dgraph_helper.MutationPredict(info.PredictTo("id")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(info.UniqueMutationName()),
			dgraph_helper.MutationItemValue(info.Ledger),
			dgraph_helper.MutationPredict(info.PredictTo("ledger")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(info.UniqueMutationName()),
			dgraph_helper.MutationItemValue(info.AssociateAccount),
			dgraph_helper.MutationPredict(info.PredictTo("associate_account")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(info.UniqueMutationName()),
			dgraph_helper.MutationItemValue(info.Content),
			dgraph_helper.MutationPredict(info.PredictTo("content")),
		),
	)
	return
}

//go:generate  pie SchemaInfos
type SchemaInfos []SchemaInfo

// All will return true if all callbacks return true. It follows the same logic
// as the all() function in Python.
//
// If the list is empty then true is always returned.
func (ss SchemaInfos) All(fn func(value SchemaInfo) bool) bool {
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
func (ss SchemaInfos) Any(fn func(value SchemaInfo) bool) bool {
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
func (ss SchemaInfos) Append(elements ...SchemaInfo) SchemaInfos {
	return append(ss, elements...)
}

// Contains returns true if the element exists in the slice.
//
// When using slices of pointers it will only compare by address, not value.
func (ss SchemaInfos) Contains(lookingFor SchemaInfo) bool {
	for _, s := range ss {
		if s == lookingFor {
			return true
		}
	}

	return false
}

// Extend will return a new slice with the slices of elements appended to the
// end.
//
// It is acceptable to provide zero arguments.
func (ss SchemaInfos) Extend(slices ...SchemaInfos) (ss2 SchemaInfos) {
	ss2 = ss

	for _, slice := range slices {
		ss2 = ss2.Append(slice...)
	}

	return ss2
}

// First returns the first element, or zero. Also see FirstOr().
func (ss SchemaInfos) First() SchemaInfo {
	return ss.FirstOr(SchemaInfo{})
}

// FirstOr returns the first element or a default value if there are no
// elements.
func (ss SchemaInfos) FirstOr(defaultValue SchemaInfo) SchemaInfo {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[0]
}

// JSONString returns the JSON encoded array as a string.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array.
func (ss SchemaInfos) JSONString() string {
	if ss == nil {
		return "[]"
	}

	// An error should not be possible.
	data, _ := json.Marshal(ss)

	return string(data)
}

// Last returns the last element, or zero. Also see LastOr().
func (ss SchemaInfos) Last() SchemaInfo {
	return ss.LastOr(SchemaInfo{})
}

// LastOr returns the last element or a default value if there are no elements.
func (ss SchemaInfos) LastOr(defaultValue SchemaInfo) SchemaInfo {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[len(ss)-1]
}

// Len returns the number of elements.
func (ss SchemaInfos) Len() int {
	return len(ss)
}

// Reverse returns a new copy of the slice with the elements ordered in reverse.
// This is useful when combined with Sort to get a descending sort order:
//
//   ss.Sort().Reverse()
//
func (ss SchemaInfos) Reverse() SchemaInfos {
	// Avoid the allocation. If there is one element or less it is already
	// reversed.
	if len(ss) < 2 {
		return ss
	}

	sorted := make([]SchemaInfo, len(ss))
	for i := 0; i < len(ss); i++ {
		sorted[i] = ss[len(ss)-i-1]
	}

	return sorted
}

// Select will return a new slice containing only the elements that return
// true from the condition. The returned slice may contain zero elements (nil).
//
// Unselect works in the opposite way as Select.
func (ss SchemaInfos) Select(condition func(SchemaInfo) bool) (ss2 SchemaInfos) {
	for _, s := range ss {
		if condition(s) {
			ss2 = append(ss2, s)
		}
	}

	return
}

// ToStrings transforms each element to a string.
func (ss SchemaInfos) ToStrings(transform func(SchemaInfo) string) pie.Strings {
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

// Transform will return a new slice where each element has been transformed.
// The number of element returned will always be the same as the input.
//
// Be careful when using this with slices of pointers. If you modify the input
// value it will affect the original slice. Be sure to return a new allocated
// object or deep copy the existing one.
func (ss SchemaInfos) Transform(fn func(SchemaInfo) SchemaInfo) (ss2 SchemaInfos) {
	if ss == nil {
		return nil
	}

	ss2 = make([]SchemaInfo, len(ss))
	for i, s := range ss {
		ss2[i] = fn(s)
	}

	return
}

// Unselect works the same as Select, with a negated condition. That is, it will
// return a new slice only containing the elements that returned false from the
// condition. The returned slice may contain zero elements (nil).
func (ss SchemaInfos) Unselect(condition func(SchemaInfo) bool) (ss2 SchemaInfos) {
	for _, s := range ss {
		if !condition(s) {
			ss2 = append(ss2, s)
		}
	}

	return
}
