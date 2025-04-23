package main

import (
	"log/slog"

	"github.com/lunagic/hephaestus/internal/hephaestus"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use: "hephaestus",
		RunE: func(cmd *cobra.Command, args []string) error {
			return hephaestus.Run()
		},
	}

	if err := root.Execute(); err != nil {
		slog.Default().Error(err.Error())
	}
}
