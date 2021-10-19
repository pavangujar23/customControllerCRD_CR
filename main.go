package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	klient "customControllerCRD_CR/pkg/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// kInfFac "customControllerCRD_CR/pkg/client/informers/"
	// "customControllerCRD_CR/pkg/controller"
)

func main() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Printf("Building config from flags failed, %s, trying to build inclusterconfig", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("error %s building inclusterconfig", err.Error())
		}
	}

	klientset, err := klient.NewForConfig(config)
	if err != nil {
		log.Printf("getting klient set %s\n", err.Error())
	}

	fmt.Println(klientset)

	klusters, err := klientset.PavangujarV1alpha1().Klusters("").List(context.Background(), metav1.ListOptions{})

	if err != nil {
		log.Printf("listing Klusters %s\n", err.Error())
	}

	fmt.Printf("length of klusters is %d and name is %s\n", len(klusters.Items), klusters.Items[0].Name)

	// client, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	log.Printf("getting std client %s\n", err.Error())
	// }

	// infoFactory := kInfFac.NewSharedInformerFactory(klientset, 20*time.Minute)
	// ch := make(chan struct{})
	// c := controller.NewController(client, klientset, infoFactory.Viveksingh().V1alpha1().Klusters())

	// infoFactory.Start(ch)
	// if err := c.Run(ch); err != nil {
	// 	log.Printf("error running controller %s\n", err.Error())
	// }
}

// specify the global tags forward Api

//controll the bheavior of code generator
// global  and local
