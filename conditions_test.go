package querybuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditions(t *testing.T) {
	negatives := []Conditions{
		Conditions{},
		Conditions{{"foo", EQ, 1}},
		Conditions{{"foo", LT, 1}},
		Conditions{{"foo", LTE, 1}},
		Conditions{{"foo", GT, 1}},
		Conditions{{"foo", GTE, 1}},
		Conditions{{"foo", GTE, 1}, {"foo", LT, 5}},
		Conditions{{"foo", GTE, 1}, {"bar", EQ, 100}},
	}
	for _, c := range negatives {
		assert.False(t, c.HasMultipleIneqFields())
	}

	positives := []Conditions{
		Conditions{{"foo", GTE, 1}, {"foo", LT, 5}, {"bar", GT, 5}},
		Conditions{{"foo", GTE, 1}, {"bar", GT, 5}},
	}
	for _, c := range positives {
		assert.True(t, c.HasMultipleIneqFields())
	}

}
