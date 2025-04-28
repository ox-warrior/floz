/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	floziov1 "github.com/ox-warrior/floz/pkg/apis/v1"
)

// nolint:unused
// log is for logging in this package.
var grouppolicieslog = logf.Log.WithName("grouppolicies-resource")

// SetupGroupPoliciesWebhookWithManager registers the webhook for GroupPolicies in the manager.
func SetupGroupPoliciesWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&floziov1.GroupPolicies{}).
		WithValidator(&GroupPoliciesCustomValidator{}).
		WithDefaulter(&GroupPoliciesCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-floz-io-floz-io-v1-grouppolicies,mutating=true,failurePolicy=fail,sideEffects=None,groups=floz.io.floz.io,resources=grouppolicies,verbs=create;update,versions=v1,name=mgrouppolicies-v1.kb.io,admissionReviewVersions=v1

// GroupPoliciesCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind GroupPolicies when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type GroupPoliciesCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &GroupPoliciesCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind GroupPolicies.
func (d *GroupPoliciesCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	grouppolicies, ok := obj.(*floziov1.GroupPolicies)

	if !ok {
		return fmt.Errorf("expected an GroupPolicies object but got %T", obj)
	}
	grouppolicieslog.Info("Defaulting for GroupPolicies", "name", grouppolicies.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-floz-io-floz-io-v1-grouppolicies,mutating=false,failurePolicy=fail,sideEffects=None,groups=floz.io.floz.io,resources=grouppolicies,verbs=create;update,versions=v1,name=vgrouppolicies-v1.kb.io,admissionReviewVersions=v1

// GroupPoliciesCustomValidator struct is responsible for validating the GroupPolicies resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type GroupPoliciesCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &GroupPoliciesCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type GroupPolicies.
func (v *GroupPoliciesCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	grouppolicies, ok := obj.(*floziov1.GroupPolicies)
	if !ok {
		return nil, fmt.Errorf("expected a GroupPolicies object but got %T", obj)
	}
	grouppolicieslog.Info("Validation for GroupPolicies upon creation", "name", grouppolicies.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type GroupPolicies.
func (v *GroupPoliciesCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	grouppolicies, ok := newObj.(*floziov1.GroupPolicies)
	if !ok {
		return nil, fmt.Errorf("expected a GroupPolicies object for the newObj but got %T", newObj)
	}
	grouppolicieslog.Info("Validation for GroupPolicies upon update", "name", grouppolicies.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type GroupPolicies.
func (v *GroupPoliciesCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	grouppolicies, ok := obj.(*floziov1.GroupPolicies)
	if !ok {
		return nil, fmt.Errorf("expected a GroupPolicies object but got %T", obj)
	}
	grouppolicieslog.Info("Validation for GroupPolicies upon deletion", "name", grouppolicies.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
