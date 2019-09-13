// Package main provides ...
package gotapper

type config struct {
	URL    string `yaml:"url"`
	Tick   int    `yaml:"tick"`
	Method string `yaml:"method"`
}
