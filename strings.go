package querybuilder

type Strings []string

func (s Strings) Has(v string) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}

func (s Strings) Except(v Strings) Strings {
	r := Strings{}
	for _, i := range s {
		if !v.Has(i) {
			r = append(r, i)
		}
	}
	return r
}
