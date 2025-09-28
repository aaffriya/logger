package config

type LoggerConfig struct {
	Stack         StackConfig      `yaml:"stack"             json:"stack"`
	Level         string           `yaml:"level"             json:"level"` // debug, info, warn, error
	DefaultFields DefaultFieldInfo `yaml:"default_fields"    json:"default_fields"`
	Pretty        PrettyConfig     `yaml:"pretty"           json:"pretty"`
}

type StackConfig struct {
	Enabled bool        `yaml:"enabled" json:"enabled"`
	Skip    int         `yaml:"skip"    json:"skip"`
	Depth   StackDepths `yaml:"depth"   json:"depth"`
}

type StackDepths struct {
	Error int `yaml:"error" json:"error"`
	Debug int `yaml:"debug" json:"debug"`
	Info  int `yaml:"info"  json:"info"`
	Warn  int `yaml:"warn"  json:"warn"`
}

type DefaultFieldInfo struct {
	Service string `yaml:"service" json:"service"`
	Version string `yaml:"version" json:"version"`
}
