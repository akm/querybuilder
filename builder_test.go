package querybuilder

import (
	"context"
	"testing"

	"google.golang.org/appengine/datastore"

	"github.com/stretchr/testify/assert"

	"github.com/akm/querybuilder/testsupport"
)

type EnumA int

const (
	EnumA0 EnumA = 0 // usually not used but used for test
	EnumA1 EnumA = 1
	EnumA2 EnumA = 2
	EnumA3 EnumA = 3
)

type Entity4Test struct {
	Int1  int
	Int2  int
	Str1  string
	Str2  string
	EnumA EnumA
}

const Kind4Test = "entity4test"

var Entities = []*Entity4Test{
	{Int1: 1, Int2: 1, Str1: "a", Str2: "foo", EnumA: EnumA1},
	{Int1: 2, Int2: 1, Str1: "b", Str2: "bar", EnumA: EnumA2},
	{Int1: 3, Int2: 2, Str1: "c", Str2: "baz", EnumA: EnumA3},
	{Int1: 4, Int2: 3, Str1: "d", Str2: "qux", EnumA: EnumA1},
	{Int1: 5, Int2: 5, Str1: "e", Str2: "quux", EnumA: EnumA2},
	{Int1: 6, Int2: 8, Str1: "f", Str2: "corge", EnumA: EnumA3},
}

func TestBuilder(t *testing.T) {
	testsupport.WithAEContext(t, func(ctx context.Context) error {
		{
			keys := make([]*datastore.Key, len(Entities))
			for i, _ := range keys {
				keys[i] = datastore.NewIncompleteKey(ctx, Kind4Test, nil)
			}
			_, err := datastore.PutMulti(ctx, keys, Entities)
			assert.NoError(t, err)
		}

		{
			b := New("Int1", "Str1", "Str2")
			assert.Equal(t, Strings{"Int1", "Str1", "Str2"}, b.ProjectFields())
			q, _ := b.Build(datastore.NewQuery(Kind4Test))
			var entities []*Entity4Test
			_, err := q.GetAll(ctx, &entities)
			assert.NoError(t, err)
			assert.Equal(t, len(Entities), len(entities))
			for _, entity := range entities {
				assert.Equal(t, 0, entity.Int2)
				assert.Equal(t, EnumA0, entity.EnumA)
			}
		}

		{
			queryValue := 1
			b := New("Int2", "Str1", "Str2")
			b.Eq("Int2", queryValue)
			assert.Equal(t, Strings{"Str1", "Str2"}, b.ProjectFields())
			q, f := b.Build(datastore.NewQuery(Kind4Test))
			var entities []*Entity4Test
			_, err := q.GetAll(ctx, &entities)
			assert.NoError(t, err)
			assert.Equal(t, 2, len(entities))
			for _, entity := range entities {
				assert.Equal(t, 0, entity.Int1)
				assert.Equal(t, 0, entity.Int2)
				assert.Equal(t, EnumA0, entity.EnumA)
			}
			assert.NoError(t, f.AssignAll(entities))
			for _, entity := range entities {
				assert.Equal(t, 0, entity.Int1)
				assert.Equal(t, queryValue, entity.Int2) // assigned by f returned from Build
				assert.Equal(t, EnumA0, entity.EnumA)
			}
		}

		{
			rangeLow := 2
			rangeHigh := 5 // not included
			b := New("Int1", "Str1", "EnumA")
			b.Gte("Int1", rangeLow)
			b.Lt("Int1", rangeHigh)
			assert.Equal(t, Strings{"Int1", "Str1", "EnumA"}, b.ProjectFields())
			q, _ := b.Build(datastore.NewQuery(Kind4Test))
			var entities []*Entity4Test
			_, err := q.GetAll(ctx, &entities)
			assert.NoError(t, err)
			assert.Equal(t, 3, len(entities))

			int1s := []int{}
			for _, entity := range entities {
				int1s = append(int1s, entity.Int1)
			}
			assert.Equal(t, []int{2, 3, 4}, int1s)
		}

		{
			b := New("Int1", "Str1", "Str2", "EnumA")
			b.Starts("Str2", "ba") // "bar" and "baz"
			assert.Equal(t, Strings{"Int1", "Str1", "Str2", "EnumA"}, b.ProjectFields())
			assert.Equal(t, Strings{"Str2"}, b.SortFields())
			q, _ := b.Build(datastore.NewQuery(Kind4Test))
			var entities []*Entity4Test
			_, err := q.GetAll(ctx, &entities)
			assert.NoError(t, err)
			assert.Equal(t, 2, len(entities))

			int1s := []int{}
			for _, entity := range entities {
				int1s = append(int1s, entity.Int1)
			}
			assert.Equal(t, []int{2, 3}, int1s)
		}

		{
			b := New("Int1", "Str1", "Str2", "EnumA")
			b.Starts("Str2", "ba") // "bar" and "baz"
			assert.Equal(t, Strings{"Int1", "Str1", "Str2", "EnumA"}, b.ProjectFields())
			assert.Equal(t, Strings{"Str2"}, b.SortFields())

			var qc *datastore.Query
			{
				qc = b.BuildForCount(datastore.NewQuery(Kind4Test))
				c, err := qc.Count(ctx)
				assert.NoError(t, err)
				assert.Equal(t, 2, c)
			}
			{
				q, f := b.BuildForList(qc)
				var entities []*Entity4Test
				_, err := q.GetAll(ctx, &entities)
				assert.NoError(t, err)
				assert.Equal(t, 2, len(entities))
				assert.NoError(t, f.AssignAll(entities))
			}
		}

		{
			b := New()
			b.Asc("Int1")
			b.Offset(2)
			b.Limit(3)
			assert.Equal(t, Strings{}, b.ProjectFields())
			q, _ := b.Build(datastore.NewQuery(Kind4Test))
			var entities []*Entity4Test
			_, err := q.GetAll(ctx, &entities)
			assert.NoError(t, err)
			assert.Equal(t, 3, len(entities))
			int1s := []int{}
			for _, entity := range entities {
				int1s = append(int1s, entity.Int1)
			}
			assert.Equal(t, []int{3, 4, 5}, int1s)
		}

		return nil
	})
}

type SubEntity struct {
	I1 int
	S1 string
}

type ComplicatedEntity4Test struct {
	ID      int
	Name    string
	Strings []string
	Ints    []int
	Sub1    SubEntity
	Subs    []SubEntity
}

const ComplicatedKind4Test = "complicated4test"

var ComplicatedEntities = []*ComplicatedEntity4Test{
	{ID: 1, Name: "Foo"},
	{ID: 2, Name: "Bar", Strings: []string{"a"}, Ints: []int{1, 2, 3}},
	{ID: 3, Name: "Baz", Strings: []string{"a", "b"}, Ints: []int{2}},
	{ID: 4, Name: "Qux", Strings: []string{"d"}, Sub1: SubEntity{I1: 1, S1: "A"}},
	{ID: 5, Name: "Quux", Strings: []string{"d"}, Subs: []SubEntity{{I1: 1, S1: "A"}, {I1: 3, S1: "C"}}},
	{ID: 6, Name: "Corge", Strings: []string{"b"}, Subs: []SubEntity{{I1: 2, S1: "B"}, {I1: 3, S1: "C"}}},
	{ID: 7, Name: "Grault", Strings: []string{"c"}, Subs: []SubEntity{{I1: 2, S1: "C"}, {I1: 4, S1: "D"}}},
	{ID: 8, Name: "Garply", Strings: []string{"d"}, Sub1: SubEntity{I1: 2, S1: "B"}},
	{ID: 9, Name: "Waldo", Strings: []string{"e"}, Sub1: SubEntity{I1: 2, S1: "B"}},
}

func TestBuilderWithComplicatedEntities(t *testing.T) {
	testsupport.WithAEContext(t, func(ctx context.Context) error {
		{
			keys := make([]*datastore.Key, len(ComplicatedEntities))
			for i, _ := range keys {
				keys[i] = datastore.NewIncompleteKey(ctx, ComplicatedKind4Test, nil)
			}
			_, err := datastore.PutMulti(ctx, keys, ComplicatedEntities)
			assert.NoError(t, err)
		}

		// Single nested field
		{
			queryValue := 2
			b := New("ID", "Name", "Sub1.I1", "Sub1.S1")
			b.Eq("Sub1.I1", queryValue)
			b.Asc("ID")
			assert.Equal(t, Strings{"ID", "Name", "Sub1.S1"}, b.ProjectFields())
			q, f := b.Build(datastore.NewQuery(ComplicatedKind4Test))
			var entities []*ComplicatedEntity4Test
			_, err := q.GetAll(ctx, &entities)
			assert.NoError(t, err)

			type pattern struct {
				ID   int
				Name string
				SubI int
				SubS string
			}

			patterns1 := []pattern{
				{8, "Garply", 0, "B"},
				{9, "Waldo", 0, "B"},
			}
			patterns2 := []pattern{
				{8, "Garply", queryValue, "B"},
				{9, "Waldo", queryValue, "B"},
			}
			impl := func(patterns []pattern) error {
				if assert.Equal(t, len(patterns), len(entities)) {
					for i, ptn := range patterns {
						e := entities[i]
						assert.Equal(t, ptn.ID, e.ID)
						assert.Equal(t, ptn.Name, e.Name)
						assert.Equal(t, ptn.SubI, e.Sub1.I1)
						assert.Equal(t, ptn.SubS, e.Sub1.S1)
					}
				}
				return nil
			}

			impl(patterns1)
			f.AssignAll(&entities)
			impl(patterns2)
		}

		// Slice field
		{
			type sub struct {
				I int
				S string
			}

			type expected struct {
				ID   int
				Name string
				Subs []*sub
			}
			assertExpected := func(ex *expected, e *ComplicatedEntity4Test) {
				assert.Equal(t, ex.ID, e.ID)
				assert.Equal(t, ex.Name, e.Name)
				if assert.Equal(t, len(ex.Subs), len(e.Subs)) {
					for i, sub := range ex.Subs {
						s := e.Subs[i]
						assert.Equal(t, sub.I, s.I1)
						assert.Equal(t, sub.S, s.S1)
					}
				}
			}

			type pattern struct {
				setup  func() ([]*ComplicatedEntity4Test, AssignFuncs)
				before []*expected
				after  []*expected
			}

			assertExpecteds := func(expecteds []*expected, entities []*ComplicatedEntity4Test) {
				if assert.Equal(t, len(expecteds), len(entities)) {
					for i, expected := range expecteds {
						assertExpected(expected, entities[i])
					}
				}
			}

			assertPattern := func(pattern *pattern) {
				entities, f := pattern.setup()
				assertExpecteds(pattern.before, entities)
				f.AssignAll(entities)
				assertExpecteds(pattern.after, entities)
			}

			genSetup := func(otherFields []string, queryValue int, distinction FilterFunc) func() ([]*ComplicatedEntity4Test, AssignFuncs) {
				return func() ([]*ComplicatedEntity4Test, AssignFuncs) {
					var entities []*ComplicatedEntity4Test
					b := New(append([]string{"ID", "Name", "Subs.I1"}, otherFields...)...)
					b.Eq("Subs.I1", queryValue)
					b.Asc("ID")
					assert.Equal(t, append(Strings{"ID", "Name"}, otherFields...), b.ProjectFields())
					q, f := b.Build(datastore.NewQuery(ComplicatedKind4Test))
					if distinction != nil {
						q = distinction(q)
					}
					_, err := q.GetAll(ctx, &entities)
					assert.NoError(t, err)
					return entities, f
				}
			}

			queryValue := 3 // Subs.I

			{
				basePattern := &pattern{
					before: []*expected{
						{5, "Quux", []*sub{{0, "A"}}}, // This SubI data must be set but datastore doesn't set it
						{5, "Quux", []*sub{{0, "C"}}},
						{6, "Corge", []*sub{{0, "B"}}}, // This SubI data must be set but datastore doesn't set it
						{6, "Corge", []*sub{{0, "C"}}},
					},
					after: []*expected{
						{5, "Quux", []*sub{{queryValue, "A"}}}, // This SubI is not value set with SubS "A". It must be 1.
						{5, "Quux", []*sub{{queryValue, "C"}}},
						{6, "Corge", []*sub{{queryValue, "B"}}}, // This SubI is not value set with SubS "B". It must be 2.
						{6, "Corge", []*sub{{queryValue, "C"}}},
					},
				}

				{ // Without Distinct
					ptn := &pattern{
						setup:  genSetup([]string{"Subs.S1"}, queryValue, nil),
						before: basePattern.before,
						after:  basePattern.after,
					}
					assertPattern(ptn)
				}

				{ // With Distinct
					ptn := &pattern{
						setup: genSetup([]string{"Subs.S1"}, queryValue, func(q *datastore.Query) *datastore.Query {
							return q.Distinct()
						}),
						before: basePattern.before,
						after:  basePattern.after,
					}
					assertPattern(ptn)
				}

				{ // With DistinctOn #1
					ptn := &pattern{
						setup: genSetup([]string{"Subs.S1"}, queryValue, func(q *datastore.Query) *datastore.Query {
							return q.DistinctOn("ID", "Name", "Subs.S1")
						}),
						before: basePattern.before,
						after:  basePattern.after,
					}
					assertPattern(ptn)
				}
			}

			{
				queryValueIgnoredPattern := &pattern{
					before: []*expected{
						{5, "Quux", []*sub{{0, "A"}}},  // This Subs.I data must be set but datastore doesn't set it
						{6, "Corge", []*sub{{0, "B"}}}, // This Subs.I data must be set but datastore doesn't set it
					},
					after: []*expected{
						{5, "Quux", []*sub{{queryValue, "A"}}},  // This Subs.I is not value set with Subs.S "A". It must be 1.
						{6, "Corge", []*sub{{queryValue, "B"}}}, // This Subs.I is not value set with Subs.S "B". It must be 2.
					},
				}

				{ // With DistinctOn #2
					ptn := &pattern{
						setup: genSetup([]string{"Subs.S1"}, queryValue, func(q *datastore.Query) *datastore.Query {
							return q.DistinctOn("ID", "Name")
						}),
						before: queryValueIgnoredPattern.before,
						after:  queryValueIgnoredPattern.after,
					}
					assertPattern(ptn)
				}

				{ // With DistinctOn #3
					ptn := &pattern{
						setup: genSetup([]string{"Subs.S1"}, queryValue, func(q *datastore.Query) *datastore.Query {
							return q.DistinctOn("ID", "Name", "Subs.I1")
						}),
						before: queryValueIgnoredPattern.before,
						after:  queryValueIgnoredPattern.after,
					}
					assertPattern(ptn)
				}
			}

			{
				subsIgnoredPattern := &pattern{
					before: []*expected{
						{5, "Quux", []*sub{}},
						{6, "Corge", []*sub{}},
					},
					after: []*expected{
						{5, "Quux", []*sub{}},
						{6, "Corge", []*sub{}},
					},
				}

				{ // Withiout Distinct
					ptn := &pattern{
						setup:  genSetup([]string{}, queryValue, nil),
						before: subsIgnoredPattern.before,
						after:  subsIgnoredPattern.after,
					}
					assertPattern(ptn)
				}

				{ // With Distinct
					ptn := &pattern{
						setup: genSetup([]string{}, queryValue, func(q *datastore.Query) *datastore.Query {
							return q.Distinct()
						}),
						before: subsIgnoredPattern.before,
						after:  subsIgnoredPattern.after,
					}
					assertPattern(ptn)
				}

				{ // With DistinctOn
					ptn := &pattern{
						setup: genSetup([]string{}, queryValue, func(q *datastore.Query) *datastore.Query {
							return q.DistinctOn("ID", "Name")
						}),
						before: subsIgnoredPattern.before,
						after:  subsIgnoredPattern.after,
					}
					assertPattern(ptn)
				}
			}

		}

		return nil
	})
}
