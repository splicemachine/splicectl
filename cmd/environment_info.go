package cmd

import (
	"context"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getEnvironmentName() string {

	cfg, err := restConfig()
	if err != nil {
		logrus.WithError(err).Fatal("could not get config")
		os.Exit(1)
	}
	if cfg == nil {
		return "default"
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("could not create client from config")
		return "default"
	}

	secretResource, secerr := client.CoreV1().Secrets("splice-system").Get(context.TODO(), "vault-key-store", v1.GetOptions{})
	if secerr != nil {
		logrus.WithError(secerr).Error("Secret Not Found vault-key-store")
		return "default"
	}

	return string(secretResource.Data["ENVIRONMENT"][:])

}

func getIngressDetail() string {
	cfg, err := restConfig()
	if err != nil {
		logrus.WithError(err).Fatal("could not get config")
		os.Exit(1)
	}
	if cfg == nil {
		return ""
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("could not create client from config")
		os.Exit(1)
	}

	ingressResult, err := client.NetworkingV1().Ingresses("splice-system").Get(context.TODO(), "splicectl-api", v1.GetOptions{})
	if err != nil {
		logrus.WithError(err).Warn("could not read from ingress: splicectl-api")
		return ""
	}
	return fmt.Sprintf("https://%s", ingressResult.Spec.Rules[0].Host)

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
