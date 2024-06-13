package main

import (
	answercmd "github.com/apache/incubator-answer/cmd"
	// remote plugins
	_ "github.com/apache/incubator-answer-plugins/connector-github"
	_ "github.com/apache/incubator-answer-plugins/connector-wallet"
	// local plugins
)

func main() {
	answercmd.Main()
}
