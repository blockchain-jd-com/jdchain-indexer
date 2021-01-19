package dgraph_helper

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type ExampleUnit struct {
	Subject MutationItem
	Object  MutationItem
	Predict MutationPredict
	Expect  string
}

var examples = []ExampleUnit{
	{
		Subject: MutationItemEmpty("sub"),
		Object:  MutationItemEmpty("obj"),
		Predict: MutationPredict("hash"),
		Expect:  "_:sub <hash> _:obj .",
	},
	{
		Subject: MutationItemEmpty("sub"),
		Object:  MutationItemValue("obj"),
		Predict: MutationPredict("hash"),
		Expect:  `_:sub <hash> "obj" .`,
	},
	{
		Subject: MutationItemEmpty("sub"),
		Object:  MutationItemUid("obj"),
		Predict: MutationPredict("hash"),
		Expect:  "_:sub <hash> <obj> .",
	},
	{
		Subject: MutationItemUid("sub"),
		Object:  MutationItemUid("obj"),
		Predict: MutationPredict("hash"),
		Expect:  "<sub> <hash> <obj> .",
	},
}

func TestNewMutation(t *testing.T) {
	for _, unit := range examples {
		testNewMutation(t, unit)
	}
}

func testNewMutation(t *testing.T, unit ExampleUnit) {
	mu := NewMutation(unit.Subject, unit.Object, unit.Predict)
	assert.Equal(t, unit.Expect, mu.String())
}

func TestMutations(t *testing.T) {
	var mutations Mutations
	assert.Equal(t, "{set{\n}}", mutations.Assembly())
	for _, unit := range examples {
		mutations = mutations.Add(NewMutation(unit.Subject, unit.Object, unit.Predict))
	}
	assembleResult := mutations.Assembly()
	for _, unit := range examples {
		assert.True(t, strings.Contains(assembleResult, unit.Expect), assembleResult+"\n"+unit.Expect)
	}
}
