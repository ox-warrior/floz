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

package queue

import (
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type PodReplacement struct {
	OriginalPod    string // 原 Pod 名称
	ReplacementPod string // 替代 Pod 名称
	Namespace      string
	Labels         map[string]string
}


func GetDeploymentKeyFromLabels(labels map[string]string, namespace string) string {

	// K8s 会在 Pod 上设置 pod-template-hash 标签识别版本
	// 我们移除这个标签来获取基础标签集
	baseLabels := make(map[string]string)
	for k, v := range labels {
		if k != "pod-template-hash" {
			baseLabels[k] = v
		}
	}

	// 生成唯一键
	key := namespace
	for k, v := range baseLabels {
		key += "/" + k + "=" + v
	}

	return key
}

type PodReplacementMap struct {
	data sync.Map
}

func NewPodReplacementMap() *PodReplacementMap {
	return &PodReplacementMap{}
}

func (m *PodReplacementMap) Store(key string, replacement PodReplacement) {
	m.data.Store(key, replacement)
	hlog.Infof("Added replacement to map with key %s: %s -> %s", key, replacement.OriginalPod, replacement.ReplacementPod)
}

func (m *PodReplacementMap) Load(key string) (PodReplacement, bool) {
	value, ok := m.data.Load(key)
	if !ok {
		return PodReplacement{}, false
	}

	replacement, ok := value.(PodReplacement)
	return replacement, ok
}

func (m *PodReplacementMap) Delete(key string) {
	m.data.Delete(key)
}
