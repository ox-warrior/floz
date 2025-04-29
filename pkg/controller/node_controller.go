// Copyright 2025 ox-warrior
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

package controller

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeController struct {
	client kubernetes.Interface
}

func NewNodeController(client kubernetes.Interface) *NodeController {
	return &NodeController{client: client}
}

func (c *NodeController) CordonNode(nodeName string) error {
	node, err := c.client.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if node.Spec.Unschedulable {
		return nil
	}

	node.Spec.Unschedulable = true
	_, err = c.client.CoreV1().Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	hlog.Infof("Node %s has been cordoned", nodeName)
	return nil
}
