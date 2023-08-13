/*
Copyright © 2023 Casey Boyer <caseyb1101@gmail.com>
*/
package cmd

import (
	"github.com/cBiscuitSurprise/strate-go/internal/web"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Startup a websockets server to interact with this game",
	Run: func(cmd *cobra.Command, args []string) {
		web.Serve(web.ServerOptions{Origin: "localhost:12345"})
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
