package values

import "github.com/shota3506/gostlc/internal/environment"

type Rho = environment.Environment[Value]

func NewRho() *Rho {
	return environment.NewEnvironment[Value]()
}
