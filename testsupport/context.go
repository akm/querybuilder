package testsupport

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

func WithAEContext(t *testing.T, f func(context.Context) error) {
	// See https://github.com/golang/appengine/blob/master/aetest/instance.go#L36-L50
	opt := &aetest.Options{StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	assert.NoError(t, err)
	defer inst.Close()

	req, err := inst.NewRequest("GET", "/", nil)
	if !assert.NoError(t, err) {
		inst.Close()
		return
	}

	if err := f(appengine.NewContext(req)); err != nil {
		t.Fatal(err)
	}
}
