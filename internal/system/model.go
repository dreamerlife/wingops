package system

type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c Config) Valid() bool {
	return c.Key != ""
}
