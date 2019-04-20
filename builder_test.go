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

var Entities = []*Entity4Test{
	{Int1: 1, Int2: 1, Str1: "a", Str2: "foo", EnumA: EnumA1},
	{Int1: 2, Int2: 1, Str1: "b", Str2: "bar", EnumA: EnumA2},
	{Int1: 3, Int2: 2, Str1: "c", Str2: "baz", EnumA: EnumA3},
	{Int1: 4, Int2: 3, Str1: "d", Str2: "qux", EnumA: EnumA1},
	{Int1: 5, Int2: 5, Str1: "e", Str2: "quux", EnumA: EnumA2},
	{Int1: 6, Int2: 8, Str1: "f", Str2: "corge", EnumA: EnumA3},
}

const Kind4Test = "entity4test"

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

		return nil
	})
}
