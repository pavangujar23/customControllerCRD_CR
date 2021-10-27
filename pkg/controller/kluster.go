package controller

import (
	"context"
	"customControllerCRD_CR/pkg/apis/pavangujar.dev/v1alpha1"
	klientset "customControllerCRD_CR/pkg/client/clientset/versioned"
	kinf "customControllerCRD_CR/pkg/client/informers/externalversions/pavangujar.dev/v1alpha1"
	klister "customControllerCRD_CR/pkg/client/listers/pavangujar.dev/v1alpha1"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const statemessage = "SUCCEED"
const statusmessage = "updated successfully"

type Controller struct {
	// clientset for custom resource kluster
	klient klientset.Interface
	// kluster has synced
	klusterSynced cache.InformerSynced
	// queue
	wq workqueue.RateLimitingInterface
	//lister
	kLister klister.KlusterLister
}

func NewController(klient klientset.Interface, klusterInformer kinf.KlusterInformer) *Controller {
	c := &Controller{
		klient:        klient,
		klusterSynced: klusterInformer.Informer().HasSynced,
		kLister:       klusterInformer.Lister(),
		wq:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "kluster"),
	}

	klusterInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)
	return c
}

func (c *Controller) Run(ch chan struct{}) error {

	if ok := cache.WaitForCacheSync(ch, c.klusterSynced); !ok {
		log.Println("cache was not synced")
	}

	go wait.Until(c.worker, time.Second, ch)
	<-ch
	return nil
}

func (c *Controller) worker() {
	for c.processNextItem() {

	}
}

func (c *Controller) processNextItem() bool {
	item, shutDown := c.wq.Get()
	if shutDown {
		// logs as well
		return false
	}

	defer c.wq.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		log.Printf("err %s calling Namespace key func on cache for item", err.Error())
	}
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.Printf("splitting key into namespace and name, error %s\n", err.Error())
		return false
	}

	kluster, err := c.kLister.Klusters(ns).Get(name)
	if err != nil {
		log.Printf("error %s, Getting the kluster resource from lister\n", err.Error())
		return false
	}

	err = c.updateStatus(kluster)
	if err != nil {
		klog.Fatal(err)
	}
	log.Printf("kluster spec that we have is %+v\n", kluster.Spec)
	return true
}

func (c *Controller) updateStatus(kluster *v1alpha1.Kluster) error {
	// get the latest version of kluster
	k, err := c.klient.PavangujarV1alpha1().Klusters(kluster.Namespace).Get(context.Background(), kluster.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	statusCopy := kluster.DeepCopy()
	statusCopy.Status.State = statemessage
	k.Status.Progress = statusmessage
	_, err = c.klient.PavangujarV1alpha1().Klusters(kluster.Namespace).UpdateStatus(context.Background(), k, metav1.UpdateOptions{})
	return err
}

func (c *Controller) handleAdd(obj interface{}) {
	log.Println("handleAdd was called")
	c.wq.Add(obj)
}

func (c *Controller) handleDel(obj interface{}) {
	log.Println("handleDel was called")
	c.wq.Add(obj)
}
