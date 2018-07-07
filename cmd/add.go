package cmd

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"

	"github.com/99designs/keyring"
	"github.com/fiveai/aws-okta/lib"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add your okta credentials",
	RunE:  add,
}

func init() {
	RootCmd.AddCommand(addCmd)
}

func AddCredentials(kr keyring.Keyring) error {
	// Ask username password from prompt
	server, err := lib.Prompt("Okta Region (emea/us)", false)
	if err != nil {
		return err
	}

	organization, err := lib.Prompt("Okta organization", false)
	if err != nil {
		return err
	}

	creds := lib.OktaCreds{
		Server:       server,
		Organization: organization,
	}

	encoded, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	item := keyring.Item{
		Key:   "okta-creds",
		Data:  encoded,
		Label: "okta credentials",
		KeychainNotTrustApplication: false,
	}

	if err := kr.Set(item); err != nil {
		return ErrFailedToSetCredentials
	}
	return nil
}

func add(cmd *cobra.Command, args []string) error {
	var allowedBackends []keyring.BackendType
	if backend != "" {
		allowedBackends = append(allowedBackends, keyring.BackendType(backend))
	}
	kr, err := lib.OpenKeyring(allowedBackends)

	if err != nil {
		log.Fatal(err)
	}

	if err := AddCredentials(kr); err != nil {
		return err
	}

	log.Infof("Added credentials for user %s", username)
	return nil
}
