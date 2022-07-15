package client

var ClientFlag clientCmdFlag

type clientCmdFlag struct {
	Host       string
	Uploadfile string
	SpecPath   string
}
