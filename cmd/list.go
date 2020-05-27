package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/segmentio/aws-okta/lib"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list will show you the profiles currently configured",
	RunE:  listRun,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func listProfiles() (lib.Profiles, error) {
	config, err := lib.NewConfigFromEnv()
	if err != nil {
		return nil, err
	}

	profiles, err := config.Parse()
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func listProfileNames(profiles lib.Profiles) []string {
	// Let's sort this list of profiles so we can have some more deterministic output:
	var profileNames []string

	for profile := range profiles {
		profileNames = append(profileNames, profile)
	}

	sort.Strings(profileNames)

	return profileNames
}

func listRun(cmd *cobra.Command, args []string) error {
	profiles, err := listProfiles()
	if err != nil {
		return err
	}

	profileNames := listProfileNames(profiles)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "PROFILE\tARN\tSOURCE_ROLE\t")
	for _, profile := range profileNames {
		v := profiles[profile]
		if role, exist := v["role_arn"]; exist {
			fmt.Fprintf(w, "%s\t%s\t%s\n", profile, role, v["source_profile"])
		}
	}
	w.Flush()

	return nil
}
