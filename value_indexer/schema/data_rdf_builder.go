package schema

import (
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"github.com/ssor/zlog"
	"github.com/tidwall/gjson"
	"strings"
)

func NewRDFBuilder(schemas ...CommonSchema) *RDFBuilder {
	if schemas == nil {
		return nil
	}

	builder := &RDFBuilder{}

	nodes, edges := SplitByType(schemas...)
	if len(nodes) != 1 {
		zlog.Warnf("builder should has just one node schema")
		return nil
	}
	builder.node = nodes[0]
	builder.edges = edges
	return builder
}

func SplitByType(schemas ...CommonSchema) (nodes NodeSchemas, edges EdgeSchemas) {
	for _, cs := range schemas {
		switch cs.Type() {
		case SchemaTypeEdge:
			schema := cs.(*EdgeSchema)
			edges = append(edges, schema)
		case SchemaTypeNode:
			schema := cs.(*NodeSchema)
			nodes = append(nodes, schema)
		}
	}
	return
}

type RDFBuilder struct {
	node  *NodeSchema
	edges EdgeSchemas
}

type UIDCache interface {
	GetUidInCache(key string) (string, bool)
}

type RDFBuildResult struct {
	Mutations dgraph_helper.Mutations
	UID       string
	QueryKey  string
}

func (result *RDFBuildResult) AddMutations(mutations dgraph_helper.Mutations) {
	result.Mutations = result.Mutations.Add(mutations...)
}

func (builder *RDFBuilder) Build(raw string, uidCache UIDCache) (result RDFBuildResult, ok bool) {
	uid, _, dataMutations, success := builder.buildNodeRDF(raw)
	if success == false {
		zlog.Warnf("build node [%s] failed", builder.node.Name)
		return
	}
	result.Mutations = result.Mutations.Add(dataMutations...)
	result.UID = uid
	result.QueryKey = uid

	edges, success := builder.buildEdges(raw, uidCache, uid)
	if success == false {
		zlog.Warnf("build edges failed")
		return
	}
	result.Mutations = result.Mutations.Add(edges...)
	ok = true
	return
}

func (builder *RDFBuilder) buildEdges(raw string, uidCache UIDCache, uid string) (mutations dgraph_helper.Mutations, ok bool) {
	result := gjson.Parse(raw)

	for _, edge := range builder.edges {
		field, success := builder.node.FindField(edge.from.PropertyName)
		if success == false {
			zlog.Warnf("cannot find field for edge %s in node %s", edge.from.PropertyName, builder.node.Name)
			return
		}

		list := result.Get(field.Name).Array()
		var values []string
		for _, result := range list {
			values = append(values, result.String())
		}
		nquads := edge.ToNQuads(uid, values, uidCache)
		for _, nquad := range nquads {
			var from, to dgraph_helper.MutationItem
			if isUID(nquad.Subject) {
				from = dgraph_helper.MutationItemUid(nquad.Subject)
			} else {
				from = dgraph_helper.MutationItemEmpty(nquad.Subject)
			}

			if isUID(nquad.Object) {
				to = dgraph_helper.MutationItemUid(nquad.Object)
			} else {
				to = dgraph_helper.MutationItemEmpty(nquad.Object)
			}
			mutations = mutations.Add(
				dgraph_helper.NewMutation(from, to, dgraph_helper.MutationPredict(nquad.Predict)),
			)
		}
	}
	ok = true
	return
}

func isUID(s string) bool {
	return strings.HasPrefix(s, "0x")
}

func (builder *RDFBuilder) buildNodeRDF(raw string) (nodeUID, value string, mutations dgraph_helper.Mutations, ok bool) {
	result := gjson.Parse(raw)
	primaryField, success := builder.node.PrimaryField()
	if success == false {
		zlog.Warnf("no primary key for schema [%s]", builder.node.Name)
		return
	}
	specName := builder.node.LowerName()

	value = result.Get(primaryField.Name).String()
	uid := builder.node.FormatCacheQueryKey(value)
	nodeUID = uid
	mutations = mutations.Add(dgraph_helper.NewMutation(
		dgraph_helper.MutationItemEmpty(uid),
		dgraph_helper.MutationItemValue(value),
		dgraph_helper.MutationPredict(primaryField.FormatPredict(specName))))

	for _, field := range builder.node.fields {
		if field.IsPrimaryKey() || field.IsListType() {
			continue
		}

		predict := field.FormatPredict(specName)
		if field.IsUidsType() {
			if !result.Get(field.Name).IsArray() {
				continue
			}
			values := result.Get(field.Name).Array()
			for _, value := range values {
				v := strings.Replace(value.String(), `"`, "", -1)
				if len(v) <= 0 {
					continue
				}
				mutations = mutations.Add(dgraph_helper.NewMutation(
					dgraph_helper.MutationItemEmpty(uid),
					dgraph_helper.MutationItemUid(v),
					dgraph_helper.MutationPredict(predict)))
			}
		} else if field.IsUidType() {
			value := result.Get(field.Name).String()
			value = strings.Replace(value, `"`, "", -1)
			if len(value) <= 0 {
				continue
			}
			mutations = mutations.Add(dgraph_helper.NewMutation(
				dgraph_helper.MutationItemEmpty(uid),
				dgraph_helper.MutationItemUid(value),
				dgraph_helper.MutationPredict(predict)))
		} else {
			value := result.Get(field.Name).String()
			value = strings.Replace(value, `"`, "", -1)
			if len(value) <= 0 {
				continue
			}
			mutations = mutations.Add(dgraph_helper.NewMutation(
				dgraph_helper.MutationItemEmpty(uid),
				dgraph_helper.MutationItemValue(value),
				dgraph_helper.MutationPredict(predict)))
		}
	}

	ok = true
	return
}
