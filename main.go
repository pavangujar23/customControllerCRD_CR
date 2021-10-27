package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	klient "customControllerCRD_CR/pkg/client/clientset/versioned"
	kInfFac "customControllerCRD_CR/pkg/client/informers/externalversions"
	"customControllerCRD_CR/pkg/controller"
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
		//Code to run Kluster inside POD
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("error %s building inclusterconfig", err.Error())
		}
	}

	klientset, err := klient.NewForConfig(config)
	if err != nil {
		log.Printf("getting klient set %s\n", err.Error())
	}

	//To Print num of nodes basic
	// klusters, err := klientset.PavangujarV1alpha1().Klusters("").List(context.Background(), metav1.ListOptions{})
	// if err != nil {
	// 	log.Printf("listing Klusters %s\n", err.Error())
	// }
	// fmt.Printf("length of klusters is %d and name is %s\n", len(klusters.Items), klusters.Items[0].Name)

	// klientset.CoreV1().Pods("").Get()
	//Watch() method of POD interface

	infoFactory := kInfFac.NewSharedInformerFactory(klientset, 20*time.Minute)
	ch := make(chan struct{})
	c := controller.NewController(klientset, infoFactory.Pavangujar().V1alpha1().Klusters())
	infoFactory.Start(ch)

	if err := c.Run(ch); err != nil {
		log.Printf("error running controller %s\n", err.Error())
	}

}

// func toIntializeInformer() {
// infoFactory := kInfFac.NewSharedInformerFactory(klientset, 20*time.Minute) // HElps in resync
// //it will get resources from all name spaces
// kInfFac.NewFilteredSharedInformerFactory(klientset, 20*time.Minute, "default", func (to *metav1.ListOptions) {
// 	//it will get resources from default name spaces
// 	//we call define list options as well
// })
// podInformer := infoFactory.Pavangujar().V1alpha1().Pods()
// podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs) {
// 	AddEvent:
// 	UpdateEvent:
// 	DeleteEvent:
// }
// infoFactory.Start(wait.NeverStop)
// infoFactory.WaitForCacheSync(wait.NeverStop) // it will call the API Server using list() then
// it will get all the respective data
// and store in
//incache memory

//from the above line we get data, Using that lister will get/List the respective data
// pod, err := podInformer.lister().PushP("default").Get("default") //to get pod from default namespace

// fmt.Printf(pod)
// }

// specify the global tags forward Api

//controll the bheavior of code generator
// global  and local

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
