package cmaestro_ingest

type Function struct {
	Name       string
	Args       []Argument
	ReturnType Primitive
}

func NewFunction(name string, returnType string) Function {
	return Function{Name: name, Args: []Argument{}, ReturnType: Primitive(returnType)}
}

func (f *Function) AddArgument(name string, returnType string) *Function {
	f.Args = append(f.Args, Argument{Name: name, Type: Primitive(returnType)})
	return f
}
