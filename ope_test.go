package querybuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpe(t *testing.T) {
	mapping := map[Ope]string{
		LT:  "<",
		LTE: "<=",
		GT:  ">",
		GTE: ">=",
		EQ:  "=",
	}

	for ope, str := range mapping {
		assert.Equal(t, str, ope.String())
		r, ok := OperatorMap[str]
		assert.True(t, ok)
		assert.Equal(t, ope, r)
	}

	invalids := []string{"!=", "<>", "=<", "=>", "=="}
	for _, invalid := range invalids {
		_, ok := OperatorMap[invalid]
		assert.False(t, ok)
	}
}
