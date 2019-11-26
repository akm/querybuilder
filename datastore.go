package querybuilder

import (
	"cloud.google.com/go/datastore"
)

type QueryFilter func(*datastore.Query) *datastore.Query
