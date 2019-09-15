// Package main provides ...
package gotapper

type Config struct {
	Tick                int              `yaml:"tick"`
	Conditions          []TestDefinition `yaml:"conditions"`
	CallBackUrlsSuccess []RequestDef     `yaml:"call_back_urls"`
	CallBackUrlsFailure []RequestDef     `yaml:"call_back_urls_failure"`
}

type RequestDef struct {
	Url         string `yaml:"url"`
	ContentType string `yaml:"content_type"`
	Body        string `yaml:"body"`
	Retries     int    `yaml:"retries"`
	Name        string `yaml:"name"`
}

type RequestResult struct {
	StatusCode int    `yaml:"status_codestatus_codeerror"`
	Error      error  `yaml:"error"`
	Name       string `yaml:"name"`
}

type TestDefinition struct {
	URL              string            `yaml:"url"`
	Arguments        map[string]string `yaml:"arguments"`
	Method           string            `yaml:"method"`
	ContentType      string            `yaml:"content_type"`
	Body             string            `yaml:"body"`
	FieldSelector    FieldSelector     `yaml:"field_selector"`
	ExpectedStatus   int               `yaml:"expected_status"`
	ExpectedType     string            `yaml:"expected_type"`
	ExpectedString   string            `yaml:"expected_string,omitempty"`
	ExpectedInt      int               `yaml:"expected_int,omitempty"`
	ExpectedNumber   float64           `yaml:"expected_float,omitempty"`
	ExpectedBool     bool              `yaml:"expected_bool,omitempty"`
	ExpectedLength   int               `yaml:"expected_length.omitenpty"`
	ExpectedDuration int               `yaml:"expected_duration,omitempty"`
}

type FieldSelector struct {
	Field     string `yaml:"field"`
	Separator string `yaml:"separator"`
}
