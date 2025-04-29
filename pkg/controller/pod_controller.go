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
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/ox-warrior/floz/pkg/queue"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodController struct {
	client kubernetes.Interface
	podMap *queue.PodReplacementMap
}

func NewPodController(client kubernetes.Interface, podMap *queue.PodReplacementMap) *PodController {
	return &PodController{
		client: client,
		podMap: podMap,
	}
}

func (c *PodController) HandleNodeDrain(nodeName string) error {
	pods, err := c.client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		// 排除daemonset的Pod
		// Skip DaemonSet pods
		if pod.OwnerReferences != nil {
			for _, owner := range pod.OwnerReferences {
				if owner.Kind == "DaemonSet" {
					continue
				}
			}
		}
		// 创建替代 Pod
		replacementPod, err := c.createReplacementPod(&pod)
		if err != nil {
			hlog.Errorf("Failed to create replacement pod for %s: %v", pod.Name, err)
			continue
		}
		hlog.Infof("Created replacement pod %s for %s", replacementPod.Name, pod.Name)
		// 生成部署标识键
		key := queue.GetDeploymentKeyFromLabels(pod.Labels, pod.Namespace)

		// 存储到 Map
		c.podMap.Store(key, queue.PodReplacement{
			OriginalPod:    pod.Name,
			ReplacementPod: replacementPod.Name,
			Namespace:      pod.Namespace,
			Labels:         pod.Labels,
		})

		// 60秒后移除原 Pod 标签
		go func(p corev1.Pod, key string) {
			time.Sleep(60 * time.Second)
			hlog.Infof("Removing labels from pod %s", p.Name)
			err := c.removePodLabels(&p)
			if err != nil {
				hlog.Errorf("Failed to remove labels from pod %s: %v", p.Name, err)
			}
		}(pod, key)
	}

	return nil
}

func (c *PodController) createReplacementPod(originalPod *corev1.Pod) (*corev1.Pod, error) {
	// 从原始 Pod 名称中提取命名模式
	// 典型的 Deployment 生成的 Pod 名称格式是: {deployment-name}-{random-suffix}
	podName := originalPod.Name

	// 去掉最后的随机后缀，保留前缀部分
	// 通常后缀是 "-xxxxx" 形式，我们可以找到最后一个"-"
	lastDashIndex := strings.LastIndex(podName, "-")
	podPrefix := podName
	if lastDashIndex > 0 {
		podPrefix = podName[:lastDashIndex]
	}

	// 创建替代 Pod
	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			// 使用相同的命名前缀，让 K8s 生成随机后缀
			GenerateName: podPrefix + "-",
			Namespace:    originalPod.Namespace,
			// 故意不设置标签
			Annotations: originalPod.Annotations,
		},
		Spec: originalPod.Spec,
	}

	// 清除一些不需要的字段
	newPod.Spec.NodeName = ""

	// 移除可能导致问题的字段
	newPod.ResourceVersion = ""
	newPod.UID = ""
	newPod.OwnerReferences = nil // 清除 OwnerReferences 避免控制器干扰

	return c.client.CoreV1().Pods(newPod.Namespace).Create(context.Background(), newPod, metav1.CreateOptions{})
}

func (c *PodController) removePodLabels(pod *corev1.Pod) error {
	pod.Labels = nil
	_, err := c.client.CoreV1().Pods(pod.Namespace).Update(context.Background(), pod, metav1.UpdateOptions{})
	return err
}
