package postgres

type Config struct {
	DSN         string `json:"dsn" yaml:"dsn"`
	MaxIdle     int    `json:"max_idle" yaml:"max_idle"`
	MaxOpen     int    `json:"max_open" yaml:"max_open"`
	MaxLifeTime int    `json:"max_life_time" yaml:"max_life_time"` // 最大连接时间，单位秒
	MaxIdleTime int    `json:"max_idle_time" yaml:"max_idle_time"` // 最大空闲时间，单位秒
	IsDebug     bool   `json:"is_debug" yaml:"is_debug"`
	Schema      string `json:"schema" yaml:"schema"`
}
