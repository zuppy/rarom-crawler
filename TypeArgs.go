package main

type TypeArgs struct {
	OutputFile string `arg:"-o,required" help:"Full path for the output csv file."`
}

func (TypeArgs) Version() string {
	return "v1.0"
}

func (TypeArgs) Description() string {
	return "RAR Parser"
}
