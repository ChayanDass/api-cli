package cmd

type APIConfig struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Method string `json:"method"`
	Token  string `json:"token"`
	Body   string `json:"body"`
}
