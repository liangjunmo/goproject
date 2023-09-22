package configtemplate

type Environment string

const (
	EnvironmentDev    Environment = "dev"
	EnvironmentTest   Environment = "test"
	EnvironmentOnline Environment = "online"
)

func (env Environment) IsOnline() bool {
	return env == EnvironmentOnline
}
