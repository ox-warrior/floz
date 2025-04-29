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

package webhook

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/ox-warrior/floz/pkg/queue"
)

// PodMutatingWebhook handles Pod mutation
type PodMutatingWebhook struct {
	client  client.Client
	decoder admission.Decoder
	podMap  *queue.PodReplacementMap
}

// NewPodMutatingWebhook creates a new PodMutatingWebhook
func NewPodMutatingWebhook(c client.Client, m *queue.PodReplacementMap) *PodMutatingWebhook {
	return &PodMutatingWebhook{
		client: c,
		podMap: m,
	}
}

// Handle handles admission requests.
func (w *PodMutatingWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	if req.Operation != admissionv1.Create {
		return admission.Allowed("not a create operation")
	}

	pod := &corev1.Pod{}
	err := w.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// 从 Pod 标签生成部署标识键
	key := queue.GetDeploymentKeyFromLabels(pod.Labels, pod.Namespace)

	// 检查是否有替换 Pod
	replacement, ok := w.podMap.Load(key)
	if !ok {
		return admission.Allowed("no replacement pod available for " + key)
	}

	// 获取替换 Pod
	replacementPod := &corev1.Pod{}
	err = w.client.Get(ctx, client.ObjectKey{
		Namespace: replacement.Namespace,
		Name:      replacement.ReplacementPod,
	}, replacementPod)
	if err != nil {
		hlog.Errorf("Failed to get replacement pod: %v", err)
		return admission.Errored(http.StatusInternalServerError, err)
	}

	// 更新替换 Pod 的标签
	replacementPod.Labels = replacement.Labels

	// 使用过一次后从 Map 中删除
	w.podMap.Delete(key)

	// 创建 patch
	marshaled, err := json.Marshal(replacementPod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	hlog.Infof("finished pod relacement, new pod name: %s", replacementPod.Name)

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}
