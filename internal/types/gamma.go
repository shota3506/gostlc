package types

import (
	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/environment"
)

type Gamma = environment.Environment[ast.Type]

func NewGamma() *Gamma {
	return environment.NewEnvironment[ast.Type]()
}
