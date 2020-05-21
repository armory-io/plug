/*
Copyright © Armory, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/armory-io/plug/pkg/serve"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a plugin bundle produced by `gradle releaseBundle`",
	Run: func(cmd *cobra.Command, args []string) {
		pd, err := cmd.Flags().GetString("plugin-dir")
		if err != nil {
			log.Fatalf("Could not get plugin directory from args: %v", err)
		}
		if pd == "" {
			log.Fatal("Must provide plugin directory")
		}

		s, err := serve.New("http://localhost", 9001, pd)
		if err != nil {
			log.Fatalf("Could not build server: %v", err)
		}

		go func() {
			if err := s.Server.ListenAndServe(); err != nil {
				log.Fatalf("Could not start server: %v", err)
			}
		}()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		log.Print("Shutting down...")
		s.Server.Shutdown(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("plugin-dir", "p", "", "path to plugin directory")
}
