package override

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v2"
)

type (
	// Resource is a resource that a component has that can be overridden
	Resource struct {
		name string
	}
	// Resources is a collection of resource
	Resources []Resource

	// Component is a deployment that has resources that can be overridden
	Component struct {
		name      string
		resources Resources
	}
	// Components is a collection of resources
	Components []Component

	// APIResourceMessage is a response or a request to/from splicectl/api with a resource attached
	APIResourceMessage struct {
		Data string `json:"data"`
	}

	// APIErrorMessage is a response from splicectl api that details why the request did not work
	APIErrorMessage struct {
		Error string `json:"error"`
	}
)

func (r Resource) Name() string {
	return r.name
}

func (r Resources) Resource(name string) (Resource, error) {
	for _, resource := range r {
		if resource.Name() == name {
			return resource, nil
		}
	}
	return Resource{}, fmt.Errorf("no resource with name: %s", name)
}

func (r Resources) Strings() []string {
	rsrcs := make([]string, len(r))
	for i, rsrc := range r {
		rsrcs[i] = rsrc.name
	}
	return rsrcs
}

func (co Component) Name() string {
	return co.name
}

func (co Component) ListResources() []string {
	return co.resources.Strings()
}

func (co Component) PrettyListResources() string {
	return prettyStrings(co.ListResources())
}

func (co Component) Resource(name string) (Resource, error) {
	return co.resources.Resource(name)
}

// GetError - creates error from string Error field that was parsed from api response
func (ae *APIErrorMessage) GetError() error {
	return fmt.Errorf("%s; api request was unsuccesful", ae.Error)
}

func responseToErrorMessage(body []byte) error {
	apiErr := &APIErrorMessage{}
	if err := json.Unmarshal(body, apiErr); err != nil {
		return err
	}
	return nil
}

func responseToYaml(resp *resty.Response, err error) (map[string]interface{}, error) {
	// check that request had no errors
	if err != nil {
		return nil, err
	}

	body := resp.Body()
	fmt.Println(string(body))
	// unmarshal api response into appropriate struct
	apiResp := &APIResourceMessage{}
	if err := json.Unmarshal(body, apiResp); err != nil {
		return nil, responseToErrorMessage(body)
	}

	// unmarshal api response data string into yaml object that represents resource
	resourceMap := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(apiResp.Data), &resourceMap); err != nil {
		return nil, err
	}

	// assuming all data has been unmarshalled properly, return it with no error
	return resourceMap, nil
}

// GetDefaultResource - gets the default resource used by the operator from
// splicectl/api for the specified component/resource.
func (co Component) GetDefaultResource(name string) (interface{}, error) {
	resource, err := co.Resource(name)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("splicectl/v1/defaults/components/%s/resources/%s", co.name, resource.name)
	fmt.Println(fmt.Sprintf("%s/%s", c.ApiServer, uri))
	data, err := responseToYaml(c.RestyWithHeaders().Get(fmt.Sprintf("%s/%s", c.ApiServer, uri)))
	return data, nil
}

// GetOverrideResource - gets the override resource used by the operator from
// splicectl/api for the specified component/resource.
func (co Component) GetOverrideResource(name string) (interface{}, error) {
	resource, err := co.Resource(name)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("splicectl/v1/overrides/components/%s/resources/%s", co.name, resource.name)
	data, err := responseToYaml(c.RestyWithHeaders().Get(fmt.Sprintf("%s/%s", c.ApiServer, uri)))
	return data, nil
}

// PutOverrideResource - puts (updates) the override resource used by the
// operator from splicectl/api for the specified component/resource.
func (co Component) PutOverrideResource(name string, data interface{}) error {
	resource, err := co.Resource(name)
	if err != nil {
		return err
	}
	overrideData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	apiReq := &APIResourceMessage{
		Data: string(overrideData),
	}
	uri := fmt.Sprintf("splicectl/v1/overrides/components/%s/resources/%s", co.name, resource.name)
	_, err = c.RestyWithHeaders().SetBody(apiReq).Put(fmt.Sprintf("%s/%s", c.ApiServer, uri))
	return err
}

func (co Components) Component(name string) (Component, error) {
	for _, component := range co {
		if component.Name() == name {
			return component, nil
		}
	}
	return Component{}, fmt.Errorf("no component with name: %s", name)
}

func (co Components) Strings() []string {
	comps := make([]string, len(co))
	for i, comp := range co {
		comps[i] = comp.name
	}
	return comps
}

// Process - is a useful hook that can do work on Components. Main use case to
// do any final processing of global components before being used by override
// commands. Processing could do things like generate aliases, look for
// duplicate names, etc.
func (co Components) Process() Components {
	return co
}

// GetComponent - gets the component with the given name from components,
// otherwise it returns an error.
func GetComponent(name string) (Component, error) {
	return components.Component(name)
}

// ListComponents - lists the components that are available for override
func ListComponents() []string {
	return components.Strings()
}

func prettyStrings(strs []string) string {
	sb := strings.Builder{}
	for _, str := range strs {
		sb.WriteString(fmt.Sprintf("\t%s\n", str))
	}
	return sb.String()
}

// PrettyListComponents - lists the components that are available for override in a nice format
func PrettyListComponents() string {
	return prettyStrings(ListComponents())
}

// below are a list of common resources shared by multiple components
var (
	// xml resources
	fairScheduler = Resource{name: "fairscheduler.xml"}
	coreSite      = Resource{name: "core-site.xml"}
	hdfsSite      = Resource{name: "hdfs-site.xml"}
	hbaseSite     = Resource{name: "hbase-site.xml"}

	// shell resources
	kafkaEnv = Resource{name: "kafka-env.sh"}
	hbaseEnv = Resource{name: "hbase-env.sh"}

	// ini resources
	shiro = Resource{name: "shiro.ini"}

	// yml resources
	promJMXAgent = Resource{name: "prom-jmx-agent-config.yml"}

	// properties resources
	log4j       = Resource{name: "log4j.properties"}
	serverProps = Resource{name: "server.properties"}
)

// components holds all of the components that can be overridden through the
// splicectl override commands.
var components = Components{
	Component{
		name: "splicedb-hbase-config",
		resources: Resources{
			fairScheduler, hbaseEnv, log4j, shiro, hbaseSite,
		},
	},
	Component{
		name: "splicedb-hbase-config-extra",
		resources: Resources{
			hbaseEnv, log4j,
		},
	},
	Component{
		name: "splicedb-kafka-config",
		resources: Resources{
			log4j, promJMXAgent, serverProps, kafkaEnv,
		},
	},
	Component{
		name: "splicedb-olap-config",
		resources: Resources{
			fairScheduler, hbaseEnv, log4j, shiro, hbaseSite,
		},
	},
	Component{
		name: "splicedb-olap-config-extra",
		resources: Resources{
			hbaseEnv, log4j,
		},
	},
	Component{
		name: "splicedb-hadoop-config",
		resources: Resources{
			coreSite, hdfsSite,
		},
	},
	Component{
		name: "splicedb-spark-config",
		resources: Resources{
			hbaseSite,
		},
	},
}.Process()
