package values

type Rho struct {
	name   string
	value  Value
	parent *Rho
}

func NewRho() Rho {
	return Rho{}
}

func (r Rho) Lookup(name string) (Value, bool) {
	if name == "" {
		return nil, false
	}
	if r.name == name {
		return r.value, true
	}
	if r.parent != nil {
		return r.parent.Lookup(name)
	}
	return nil, false
}

func (r Rho) Bind(name string, value Value) Rho {
	return Rho{
		name:   name,
		value:  value,
		parent: &r,
	}
}
