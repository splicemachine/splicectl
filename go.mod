module github.com/splicemachine/splicectl

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/blang/semver/v4 v4.0.0
	github.com/go-resty/resty/v2 v2.2.0
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/maahsome/gron v0.1.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/yaml v1.2.0
)
