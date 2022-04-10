package api

import "os"

//      Warning!!!!!!!!!!!This file is readonly!Don't modify this file!

const pkg = "github.com/chenjie199234/config"
const Name = "config"

var Group = os.Getenv("GROUP")

func init() {
	if Group == "" || Group == "<GROUP>" {
		panic("missing GROUP env")
	}
}
