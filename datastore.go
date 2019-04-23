package querybuilder

import (
	"google.golang.org/appengine/datastore"
)

type QueryFilter func(*datastore.Query) *datastore.Query
