package get

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/common"
	"k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type namedURL struct {
	name, url string
}

const (
	oauthProxySuffix = "-oauth2-proxy"
	ssNameSpace      = "splice-system"
	displayNameLabel = "displayName"
	defaultNameLabel = "app"
)

var getUrlsCmd = &cobra.Command{
	Use:   "urls",
	Short: "Get a list of core urls.",
	Long: `EXAMPLES
	# Get a list of core urls that are enabled
	splicectl get urls

	# Get a list of core and database urls that are enabled
	splicectl get urls --database-name splicedb

	# If either --build/-b or --prod/-p are specified only those urls are output.
	# The build flag take precedence over prod.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if out := getURLOutput(cmd); out != "" {
			fmt.Println("Your urls are:", out)
		} else {
			fmt.Println("Could not get any valid urls.")
		}
	},
}

func getURLOutput(cmd *cobra.Command) string {
	if buildOnly, _ := cmd.Flags().GetBool("build"); buildOnly {
		return generateBuildURLs()
	} else if prodOnly, _ := cmd.Flags().GetBool("prod"); prodOnly {
		return generateBuildURLs()
	} else {
		fmt.Println(common.DatabaseName(cmd))
		return generateURLsFromNamespaces(ssNameSpace, common.DatabaseName(cmd))
	}
}

func generateBuildURLs() string {
	return `
Kube:                   https://kube.build.splicemachine-dev.io/
Engineering Dashboard:  https://dashboard.build.splicemachine-dev.io/
`
}

func generateProdURLs() string {
	return `
Dashboard:            https://cloud-dashboard.splicemachine.io/
Kibana:               https://cloudadmin.splicemachine.io/kibana
Chronograf:           https://cloudadmin.splicemachine.io/chronograf
Oauth:                https://cloudadmin.splicemachine.io/oauth2
Cloud Manager Admin:  https://cloudadmin.splicemachine.io
Cloud Manager:        https://cloud.splicemachine.io
`
}

func generateURLsFromNamespaces(namespaces ...string) string {
	// Get kube client
	client, err := common.KubeClient()
	if err != nil {
		logrus.WithError(err).Error("could not get kube config to generate core urls")
		return ""
	}

	// Create named urls for outputting
	urls := make([]namedURL, 0)
	for _, namespace := range namespaces {
		if namespace == "" {
			continue
		}

		// Get ingresses
		ings, err := client.
			NetworkingV1beta1().
			Ingresses(namespace).
			List(context.TODO(), v1.ListOptions{})
		if err != nil {
			logrus.WithError(err).Error("could not list ingresses for splice system")
			return ""
		}

		// Get urls from list of ingresses
		urls = append(urls, generateURLsFromIngresses(ings)...)
	}

	// Generate output for console from list of urls
	return generateOutputFromNamedURLs(urls)
}

func generateURLsFromIngresses(ings *v1beta1.IngressList) []namedURL {
	urls := make([]namedURL, 0)

	// Iterate through each ingress
	for _, ing := range ings.Items {
		// Do not include oauth2-proxies
		if strings.Contains(strings.ToLower(ing.Name), oauthProxySuffix) {
			continue
		}

		// Get the preferred name for the ingress from labels
		name := preferredName(ing.Labels)

		for _, rule := range ing.Spec.Rules {
			host := rule.Host

			// Create a named url for each path in the ingress
			for _, path := range rule.HTTP.Paths {
				urls = append(urls, namedURL{
					name: name + ":",
					url:  fmt.Sprintf("https://%s%s", host, path.Path),
				})
			}
		}
	}
	return urls
}

func preferredName(labels map[string]string) string {
	name, ok := labels[displayNameLabel]
	if ok {
		return name
	}
	name = labels[defaultNameLabel]
	return strings.Title(name)
}

func generateOutputFromNamedURLs(urls []namedURL) string {
	// Sort the urls by name
	sort.Slice(urls, func(i, j int) bool {
		return strings.Compare(urls[i].name, urls[j].name) < 0
	})

	// Get the longest name to help right justify urls
	nameLen := -1
	for _, pair := range urls {
		if len(pair.name) > nameLen {
			nameLen = len(pair.url)
		}
	}

	// Concat all names/urls together, one pair per line.
	sb := strings.Builder{}
	sb.WriteString("\n")
	for _, pair := range urls {
		sb.WriteString(fmt.Sprintf("%s %s\n", pair.name, pair.url))
	}

	return sb.String()
}

func init() {
	getUrlsCmd.Flags().BoolP("build", "b", false, "Whether to output build urls")
	getUrlsCmd.Flags().BoolP("prod", "p", false, "Whether to output production urls")

	// add database name and aliases
	getUrlsCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	getUrlsCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	getUrlsCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	getCmd.AddCommand(getUrlsCmd)
}
