package dgraph_helper

import (
	"fmt"
	"strings"
)

type MutationItemValue string

func (mue MutationItemValue) Raw() string {
	return string(mue)
}

func (mue MutationItemValue) Value() string {
	if mue == "*" {
		return string(mue)
	}
	return fmt.Sprintf(`"%s"`, string(mue))
}

type MutationItemEmpty string

func (mue MutationItemEmpty) Raw() string {
	return string(mue)
}
func (mue MutationItemEmpty) Value() string {
	return fmt.Sprintf("_:%s", mue)
}

type MutationItemUid string

func (miu MutationItemUid) Raw() string {
	return string(miu)
}

func (miu MutationItemUid) Value() string {
	return fmt.Sprintf("<%s>", miu)
}

type MutationItem interface {
	Value() string
	Raw() string
}

type MutationPredict string

func (mp MutationPredict) Raw() string {
	return string(mp)
}

func (mp MutationPredict) String() string {
	if mp == "*" {
		return string(mp)
	}
	return fmt.Sprintf("<%s>", string(mp))
}

func NewMutation(subject, object MutationItem, predict MutationPredict) *Mutation {
	return &Mutation{
		Subject: subject,
		Object:  object,
		Predict: predict,
	}
}

type Mutation struct {
	Subject   MutationItem
	Object    MutationItem
	Predict   MutationPredict
	isPrimary bool
}

func (mu *Mutation) SetPrimary(b bool) {
	mu.isPrimary = b
}

func (mu *Mutation) SetUid(uid string) {
	mu.Subject = MutationItemUid(uid)
}

func (mu *Mutation) String() string {
	return fmt.Sprintf("%s %s %s .", mu.Subject.Value(), mu.Predict.String(), mu.Object.Value())
}

type Mutations []*Mutation

func (mutations Mutations) Add(newMutations ...*Mutation) Mutations {
	return append(mutations, newMutations...)
}

func (mutations Mutations) Primary() *Mutation {
	for _, mu := range mutations {
		if mu.isPrimary {
			return mu
		}
	}
	return nil
}
func (mutations Mutations) IsEmpty() bool {
	return mutations == nil || len(mutations) <= 0
}

func (mutations Mutations) Assembly() string {
	var builder strings.Builder
	for _, mu := range mutations {
		builder.WriteString(mu.String())
		builder.WriteString("\n")
	}
	//builder.WriteString("}}")
	return builder.String()
}
