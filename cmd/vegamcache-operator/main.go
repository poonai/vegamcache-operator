package main

import (
	"fmt"

	vegamclient "github.com/sch00lb0y/vegamcache-operator/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

const configFile = "/home/schoolgirl/.kube/config"

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		fmt.Print(err)
	}
	vegamClient, err := vegamclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	list, _ := vegamClient.VegamcacheoperatorV1alpha1().VegamCaches("").List(metav1.ListOptions{})
	fmt.Print(list)
}
