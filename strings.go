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
	return s.Filter(func(_ Strings, i string) bool {
		return !v.Has(i)
	})
}

func (s Strings) Uniq() Strings {
	return s.Filter(func(r Strings, i string) bool {
		return !r.Has(i)
	})
}

func (s Strings) Filter(f func(Strings, string) bool) Strings {
	r := Strings{}
	for _, i := range s {
		if f(r, i) {
			r = append(r, i)
		}
	}
	return r
}
