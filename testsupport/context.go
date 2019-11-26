package testsupport

import (
	"context"
	"testing"
)

func WithAEContext(t *testing.T, f func(context.Context) error) {
	if err := f(context.Background()); err != nil {
		t.Fatal(err)
	}
}
