module github.com/fiveai/aws-okta

go 1.15

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

require (
	github.com/99designs/keyring v1.0.0
	github.com/alessio/shellescape v0.0.0-20190409004728-b115ca0f9053
	github.com/aws/aws-sdk-go v1.25.25
	github.com/marshallbrekka/go-u2fhost v0.0.0-20200114212649-cc764c209ee9
	github.com/mitchellh/go-homedir v1.1.0
	github.com/segmentio/aws-okta v1.0.2
	github.com/sirupsen/logrus v1.4.1
	github.com/skratchdot/open-golang v0.0.0-20160302144031-75fb7ed4208c
	github.com/spf13/cobra v0.0.0-20170621173259-31694f19adee
	github.com/stretchr/testify v1.3.0
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
	gopkg.in/ini.v1 v1.42.0
)
