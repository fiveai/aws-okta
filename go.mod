module github.com/fiveai/aws-okta

go 1.15

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

retract v0.19.4

require (
	github.com/99designs/keyring v1.1.6
	github.com/alessio/shellescape v1.4.1
	github.com/aws/aws-sdk-go v1.37.11
	github.com/marshallbrekka/go-u2fhost v0.0.0-20210111072507-3ccdec8c8105
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.7.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.3.0
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/ini.v1 v1.62.0
)
