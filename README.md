# aws-okta

`aws-okta` allows you to authenticate with AWS using your Okta credentials.

## Installing

[Download a binary from our release page](https://github.com/FiveAI/aws-okta/releases)

### Windows

See [docs/windows.md](docs/windows.md) for information on getting this working
with Windows.

## Usage

### Adding Okta credentials

```bash
$ aws-okta add
```

This will prompt you for your Okta organization, custom domain, region,
username, and password. These credentials will then be stored in your keyring
for future use.

### Configuring your AWS config

The most seamless way to use `aws-okta` is via AWS SDK config credential
process. To do so you will need to setup your AWS config as `aws-okta` assumes
that your base role is one that has been configured for Okta's SAML integration
by your Okta admin.

Okta provides a guide for setting up that integration
[here](https://help.okta.com/en/prod/Content/Topics/DeploymentGuides/AWS/aws-deployment.htm).
During that configuration, your admin should be able to grab the AWS App Embed
URL from the General tab of the AWS application in your Okta org.  You will need
to set that value in your `~/.aws/config` file, for example:

```ini
[okta]
aws_saml_url = home/amazon_aws/0ac4qfegf372HSvKF6a3/965
```

Next, you need to set up your base Okta role. This will be one your admin
created while setting up the integration. It should be specified similar to
to any other aws profile. However, instead of `role_arn` we will use
`okta_aws_role_arn`:

```
[profile okta-dev]
# okta_aws_role_arn is used here instead of role_arn because
# role_arn requires either source_profile or credential_source to be set
okta_aws_role_arn = arn:aws:iam::<account-id>:role/<okta-role-name>
credential_process = aws-okta cred-process okta-dev
```

Your setup may require additional roles to be configured if your admin has set
up a more complicated role scheme like cross account roles.  For more details on
the authentication process, see the internals section.

### cred-process

```bash
aws-okta cred-process profile
```

`cred-process` generates the required JSON output for the AWS CLI external
process sourcing. This is set in the `credential_process` field in the AWS
config file.

If your organization requires MFA it's advices to use [Okta
Verify](https://help.okta.com/en/prod/end-user/Content/Topics/end-user/ov-new-install.htm)
as the AWS CLI may hide output required to completing the MFA process. A fuller
example of the `credential_process` command would be:

```ini
[profile okta-dev]
okta_aws_role_arn = arn:aws:iam::<account-id>:role/<okta-role-name>
credential_process = aws-okta cred-process --mfa-factor-type push --mfa-provider okta okta-dev
```

```
Usage:
  aws-okta cred-process <profile> [flags]

Examples:
[profile foo]
credential_process = aws-okta cred-process profile

Flags:
  -a, --assume-role-ttl duration   Expiration time for assumed role (default 1h0m0s)
  -h, --help                       help for cred-process
  -p, --pretty                     Pretty print display
  -t, --session-ttl duration       Expiration time for okta role session (default 1h0m0s)

Global Flags:
  -b, --backend string              Secret backend to use [keychain pass file]
  -d, --debug                       Enable debug logging
      --mfa-duo-device string       Device to use phone1, phone2, u2f or token (default "phone1")
      --mfa-factor-type string      MFA Factor Type to use (eg push, token:software:totp)
      --mfa-provider string         MFA Provider to use (eg DUO, OKTA, GOOGLE)
      --session-cache-single-item   (alpha) Enable single-item session cache; aka AWS_OKTA_SESSION_CACHE_SINGLE_ITEM
```


#### A more complex example

The `aws_saml_url` can be set in the "okta" ini section, or on a per profile
basis. This is useful if, for example, your organization has several Okta Apps
(i.e. one for dev/qa and one for prod, or one for internal use and one for
integrations with third party providers). For example:

```ini
[okta]
# This is the "default" Okta App
aws_saml_url = home/amazon_aws/cuZGoka9dAIFcyG0UllG/214

[profile okta-dev]
okta_aws_role_arn = arn:aws:iam::<account-id>:role/<okta-role-name>
credential_process = aws-okta cred-process okta-dev

[profile integrations-auth]
# This is a distinct Okta App
aws_saml_url = home/amazon_aws/woezQTbGWUaLSrYDvINU/214
okta_aws_role_arn = arn:aws:iam::<account-id>:role/<okta-role-name>
credential_process = aws-okta cred-process integrations-auth

[profile vendor]
# This profile uses the "integrations-auth" Okta app combined with secondary role assumption
source_profile = integrations-auth
role_arn = arn:aws:iam::<account-id>:role/<secondary-role-name>

[profile testaccount]
# This stores the Okta session in a separate item in the Keyring.
# This is useful if the Okta session is used or modified by other applications
# and needs to be isolated from other sessions. It is also useful for
# development versions or multiple versions of aws-okta running.
okta_session_cookie_key = okta-session-cookie-test
role_arn = arn:aws:iam::<account-id>:role/<okta-role-name>
```

The configuration above means that you can use multiple Okta Apps at the same time and switch between them easily.

#### Multiple Okta accounts
setup accounts:
```ini
aws-okta add --account=account-a
aws-okta add --account=account-b
```

define keyring key for each profile:
```ini
[profile account-a]
# This is a distinct Okta App
aws_saml_url = home/amazon_aws/woezQTbGWUaLSrYDvINU/214
okta_aws_rolerole_arn = arn:aws:iam::<account-id>:role/<okta-role-name>
okta_account_name = account-a

[profile account-b]
aws_saml_url = home/amazon_aws/woezQTbGaDAA4rYDvINU/123
okta_aws_role_arn = arn:aws:iam::<account-id>:role/<okta-role-name>
okta_account_name = account-b
```

#### Configuring Okta assume role and AWS assume role TTLs

The default TTLs for both the initial SAML assumed role and secondary AWS assumed roles are 1 hour.  This means that AWS credentials will expire every hour.

* *session-ttl*: Duration of initial role assumed by Okta
* *assume-role-ttl*: Duration of second role assumed

In addition to specifying session and AWS assume role TTLs with command-line flags, they can be set using environment variables.

```bash
export AWS_SESSION_TTL=1h
export AWS_ASSUME_ROLE_TTL=1h
```

The AWS assume role TTL can also be set per-profile in the aws config:

```ini
# Example with an initial and secondary role that are configured with a max session duration of 12 hours
[profile ttldemo]
aws_saml_url = home/amazon_aws/cuZGoka9dAIFcyG0UllG/214
okta_aws_role_arn = arn:aws:iam::<account-id>:role/<okta-role-name>
session_ttl = 12h

[profile ttldemo-role]
source_profile = ttldemo
okta_aws_role_arn = arn:aws:iam::<account-id>:role/<secondary-role-name>
assume_role_ttl = 12h
```

#### Multi-factor Authentication (MFA) configuration

If you have a single MFA factor configured, that factor will be automatically
selected.  By default, if you have multiple available MFA factors, then you will
be prompted to select which one to use.  However, if you have multiple factors
and want to specify which factor to use, you can do one of the following:

* Specify on the command line with `--mfa-provider` and `--mfa-factor-type`
* Specify with environment variables `AWS_OKTA_MFA_PROVIDER` and `AWS_OKTA_MFA_FACTOR_TYPE`
* Specify in your aws config with `mfa_provider` and `mfa_factor_type`

### Exec

When using applications that do not integrate the AWS SDK and need credentials
to be in the environment you can used `aws-okta exec`

```bash
$ aws-okta exec <profile> -- <command>
```

Exec will assume the role specified by the given aws config profile and execute
a command with the proper environment variables set.  This command is a drop-in
replacement for `aws-vault exec` and accepts all of the same command line flags:

```bash
$ aws-okta help exec
exec will run the command specified with aws credentials set in the environment

Usage:
  aws-okta exec <profile> -- <command>

Flags:
  -a, --assume-role-ttl duration   Expiration time for assumed role (default 1h0m0s)
  -h, --help                       help for exec
  -t, --session-ttl duration       Expiration time for okta role session (default 1h0m0s)

Global Flags:
  -b, --backend string   Secret backend to use [kwallet secret-service file] (default "file")
  -d, --debug            Enable debug logging
```

### Exec for EKS and Kubernetes

`aws-okta` can also be used to authenticate `kubectl` to your AWS EKS cluster.
Assuming you have [installed
`kubectl`](https://docs.aws.amazon.com/eks/latest/userguide/install-kubectl.html),
[setup your
kubeconfig](https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html).
you can now access your EKS cluster with `kubectl`. Note that on a new cluster,
your Okta CLI user needs to be using the same assumed role as the one who
created the cluster. Otherwise, your cluster needs to have been configured to
allow your assumed role.

```bash
$ kubectl version --short
```

### Shell completion

`aws-okta` provides shell completion support for BASH and ZSH via the `aws-okta
completion` command.

## Backends

We use 99design's keyring package that they use in `aws-vault`.  Because of
this, you can choose between different pluggable secret storage backends just
like in `aws-vault`.  You can either set your backend from the command line as a
flag, or set the `AWS_OKTA_BACKEND` environment variable.

For Linux / Ubuntu add the following to your bash config / zshrc etc:

```
export AWS_OKTA_BACKEND=secret-service
```

## --session-cache-single-item aka AWS_OKTA_SESSION_CACHE_SINGLE_ITEM (alpha)

This flag enables a new secure session cache that stores all sessions in the
same keyring item. For macOS users, this means drastically fewer authorization
prompts when upgrading or running local builds.

No provision is made to migrate sessions between session caches.

Implemented in [https://github.com/segmentio/aws-okta/issues/146](#146).

## Local Development

If you're developing in Linux, you'll need to get `libusb`. For Ubuntu, install
the libusb-1.0-0-dev or use the `Dockerfile` provided in the repo.

## Running Tests

`make test`

## Internals

### Authentication process

We use the following multiple step authentication:

- Step 1 : Basic authentication against Okta
- Step 2 : MFA challenge if required
- Step 3 : Get AWS SAML assertion from Okta
- Step 4 : Assume base okta role from profile with the SAML Assertion
- Step 5 : Assume the requested AWS Role from the targeted AWS account to generate STS credentials
