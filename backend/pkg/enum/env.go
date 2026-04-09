package enum

type Env string

const (
	EnvDev  Env = "debug"
	EnvStag Env = "test"
	EnvProd Env = "release"
)

func (e Env) String() string {
	return string(e)
}
