/*
Copyright 2024.

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

package v1alpha1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var jcljoblog = logf.Log.WithName("jcljob-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *JCLJob) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-kzed-io-v1alpha1-jcljob,mutating=true,failurePolicy=fail,sideEffects=None,groups=kzed.io,resources=jcljobs,verbs=create;update,versions=v1alpha1,name=mjcljob.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &JCLJob{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *JCLJob) Default() {
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
//+kubebuilder:webhook:path=/validate-kzed-io-v1alpha1-jcljob,mutating=false,failurePolicy=fail,sideEffects=None,groups=kzed.io,resources=jcljobs,verbs=create;update,versions=v1alpha1,name=vjcljob.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &JCLJob{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *JCLJob) ValidateCreate() (admission.Warnings, error) {
	jcljoblog.Info("validate create", "name", r.Name)

	definedJCLSources := 0
	if r.Spec.DSPath != "" {
		definedJCLSources++
	}
	if r.Spec.JCL != "" {
		definedJCLSources++
	}
	if r.Spec.USSPath != "" {
		definedJCLSources++
	}

	if definedJCLSources > 1 {
		return nil, errors.New("only one of spec.jcl, spec.dsPath or spec.ussPath can be specified")
	}
	if definedJCLSources == 0 {
		return nil, errors.New("one of spec.jcl, spec.dsPath or spec.ussPath must be specified")
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *JCLJob) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *JCLJob) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}
