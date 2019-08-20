// Copyright © 2018 Andreas Fritzler <andreas.fritzler@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"encoding/json"
	"fmt"

	kutype "github.com/afritzler/kube-universe/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespaceType = "namespace"
	podType       = "pod"
	nodeType      = "node"
)

// GetGraph returns the rendered dependency graph
func GetGraph(kubeconfig string) ([]byte, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset for kubeconfig: %s", err)
	}

	nodes := make(map[string]*kutype.Node)
	links := make([]kutype.Link, 0)

	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %s", err)
	}
	for _, n := range namespaces.Items {
		key := fmt.Sprintf("%s-%s", namespaceType, n.Name)
		nodes[key] = &kutype.Node{Id: key, Name: n.Name, Type: namespaceType, Namespace: n.Namespace}
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %s", err)
	}
	for _, p := range pods.Items {
		podKey := fmt.Sprintf("%s-%s", p.Namespace, p.Name)
		namespaceKey := fmt.Sprintf("%s-%s", namespaceType, p.Namespace)
		nodeKey := fmt.Sprintf("%s-%s", nodeType, p.Spec.NodeName)
		nodes[podKey] = &kutype.Node{
			Id:            podKey,
			Name:          p.Name,
			Type:          podType,
			Namespace:     p.Namespace,
			Status:        string(p.Status.Phase),
			StatusMessage: p.Status.Message}
		links = append(links, kutype.Link{Source: namespaceKey, Target: podKey, Value: 0})
		links = append(links, kutype.Link{Source: podKey, Target: nodeKey, Value: 0})
	}

	clusterNodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %s", err)
	}
	for _, n := range clusterNodes.Items {
		key := fmt.Sprintf("%s-%s", nodeType, n.Name)
		nodes[key] = &kutype.Node{Id: key, Name: n.Name, Type: nodeType, Namespace: n.Namespace, Status: string(n.Status.Phase)}
	}

	data, err := json.MarshalIndent(kutype.Graph{Nodes: values(nodes), Links: &links}, "", "	")
	if err != nil {
		return nil, fmt.Errorf("JSON marshaling failed: %s", err)
	}
	return data, nil
}

func values(nodes map[string]*kutype.Node) *[]kutype.Node {
	array := []kutype.Node{}
	for _, n := range nodes {
		array = append(array, *n)
	}
	return &array
}
