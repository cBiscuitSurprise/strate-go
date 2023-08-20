/*
Copyright © 2023 Casey Boyer <caseyb1101@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/cBiscuitSurprise/strate-go/internal/web"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Startup a websockets server to interact with this game",
	Run: func(cmd *cobra.Command, args []string) {
		origin, err := cmd.Flags().GetString("origin")
		if err != nil {
			log.Error().Err(err).Msgf("failed to launch server, bad option for 'origin'")
		}

		port, _ := cmd.Flags().GetString("port")

		if err != nil {
			log.Error().Err(err).Msgf("failed to launch server, bad option for 'port'")
		}

		if err == nil {
			web.Serve(web.ServerOptions{Origin: fmt.Sprintf("%s:%s", origin, port)})
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().String("origin", "localhost", "The ip address to bind to (e.g. '0.0.0.0')")
	serveCmd.PersistentFlags().String("port", "1300", "The port address to bind to (e.g. '80')")
}
