package server

var ServerFlag serverCmdFlag

type serverCmdFlag struct {
	Port    string
	Store   string
	Urlpath string
}
