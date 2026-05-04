package main

import (
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/spf13/cobra"
)

const (
	VERSION         = "(untracked)"
	CONFIG_FILENAME = "kaeruconfig.yml"
)

var cfgPath string

var rootCmd = &cobra.Command{
	Use:   "kaeru",
	Short: "An open source, self-hosted server management tool in a single binary",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", CONFIG_FILENAME, "configuration file path")
}

// @title			Kaeru API
// @version		0.1
// @description	An open source, self-hosted server management tool
// @termsOfService	https://kaeru.net/tos
//
// @contact.name	Gabriel Lake
// @contact.url	https://lakelimbo.com
// @contact.email	gabriel@lakelimbo.com
//
// @license.name	GPLv3
// @license.url	https://www.gnu.org/licenses/gpl-3.0.html
func main() {
	rootCmd.AddCommand(serveCmd)

	rootCmd.Execute()
}

func loadConfigErr() {
	logger.Fatal("Error loading config. Could not start the server")
}
