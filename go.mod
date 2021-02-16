module github.com/splicemachine/splicectl

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/evanphx/json-patch v4.2.0+incompatible
	github.com/go-resty/resty/v2 v2.2.0
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/maahsome/gron v0.1.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/tidwall/sjson v1.1.1
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.4 // indirect
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v0.18.4
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66 // indirect
	sigs.k8s.io/yaml v1.2.0
	k8s.io/kubernetes v1.13.0
)
