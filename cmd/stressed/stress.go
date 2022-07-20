package stressed

import (
	"fmt"

	"github.com/rsuww-load-reaper/pkg/stressed"
	"github.com/spf13/cobra"
)

var stressCmd = &cobra.Command{
	Use:     "stresstest",
	Aliases: []string{"stress"},
	Short:   "stress test a dev-ditto API endpoint for end-users only",
	Args:    cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		res := stressed.StressTest(args[0], args[1])
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(stressCmd)
}
