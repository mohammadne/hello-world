package main

import (
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/mohammadne/hello-world/cmd"
)

func main() {
	const description = "HelloWorld is an xray utility program which helps you send hello packet to the world"
	root := &cobra.Command{Short: description}

	root.AddCommand(
		cmd.Executer{}.Command(),
	)

	if err := root.Execute(); err != nil {
		log.Fatal("failed to execute root command", zap.Error(err))
	}
}
