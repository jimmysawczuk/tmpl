package config

type Block struct {
	In     string `json:"in"`
	Out    string `json:"out"`
	Format string `json:"format"`

	Options Options `json:"options"`
}

type Options struct {
	Minify bool                   `json:"minify"`
	Env    map[string]string      `json:"env"`
	Params map[string]interface{} `json:"params"`
}
