package auth

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/common"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// This is the way
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// Client - Our primary client interface
type Client interface {
	RetrieveTokenBearer() bool
	CheckTokenValidity() bool
	GetTokenBearer() string
	GetSessionID() string
	GetSession() common.SessionData
}

// Info - Our auth properties
type Info struct {
	TokenBearer string
	Session     common.SessionData
	Environment string
}

// NewAuth - return an interface to the Auth routines
func NewAuth(environmentName string, sess common.SessionData) Client {
	return &Info{
		Environment: environmentName,
		Session:     sess,
	}
}

// GetTokenBearer - Return the Token Bearer
func (i *Info) GetTokenBearer() string {
	return i.TokenBearer
}

// GetSessionID - Return the SessionID
func (i *Info) GetSessionID() string {
	return i.Session.SessionID
}

// GetSession - Return Session Info
func (i *Info) GetSession() common.SessionData {
	return i.Session
}

// RetrieveTokenBearer - fetch the token from the K8s secret
func (i *Info) RetrieveTokenBearer() bool {
	cfg, err := restConfig()
	if err != nil {
		logrus.WithError(err).Fatal("could not get config")
		os.Exit(1)
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("could not create client from config")
		os.Exit(1)
	}

	secretResult, err := client.CoreV1().Secrets("splice-system").Get(context.TODO(), "splicectl-api-tokens", v1.GetOptions{})
	if err != nil {
		logrus.WithError(err).Fatal("could not read from secret: splicectl-api-tokens")
		os.Exit(1)
	}

	bearerPath := fmt.Sprintf("%s_token-bearer", i.Session.SessionID)
	i.TokenBearer = strings.TrimSpace(string(secretResult.Data[bearerPath]))
	if i.TokenBearer == "" {
		return false
	}
	return true
}

// CheckTokenValidity - Verify that the token is still good.
func (i *Info) CheckTokenValidity() bool {

	if i.Session.SessionID == "" || i.Session.ValidUntil == "" {
		return false
	}

	notAfter, _ := time.Parse(time.RFC3339, i.Session.ValidUntil)
	timeNow := time.Now().UTC()

	if notAfter.Before(timeNow) {
		return false
	}
	return i.RetrieveTokenBearer()
}

func restConfig() (*rest.Config, error) {
	// We aren't likely to run this INSIDE the K8s cluster, this routine
	// simply picks up the config from the file system of a running POD.
	// kubeCfg, err := rest.InClusterConfig()
	var kubeCfg *rest.Config
	var err error

	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		kubeCfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logrus.Info("No KUBECONFIG ENV")
			return nil, err
		}
	} else {
		// ENV KUBECONFIG not set, check for ~/.kube/config
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		kubeFile := fmt.Sprintf("%s/%s", home, ".kube/config")
		if _, err := os.Stat(kubeFile); err != nil {
			if os.IsNotExist(err) {
				if os.Args[1] != "version" {
					logrus.Info("Could not locate the KUBECONFIG file, normally ~/.kube/config")
					os.Exit(1)
				}
				return nil, nil
			}
		}
		kubeCfg, err = clientcmd.BuildConfigFromFlags("", kubeFile)
	}
	return kubeCfg, nil
}
