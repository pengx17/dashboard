// Copyright 2017 The Kubernetes Authors.
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

package application

import (
	"log"

	applicationApi "github.com/kubernetes-sigs/application/pkg/apis/app/v1alpha1"
	applicationAlphaClient "github.com/kubernetes-sigs/application/pkg/client/clientset/versioned"
	"github.com/kubernetes/dashboard/src/app/backend/api"
	"github.com/kubernetes/dashboard/src/app/backend/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplicationGenericComponentList struct {
	ListMeta api.ListMeta `json:"listMeta"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// ApplicationDetail is a presentation layer view of Kubernetes-sigs Application resource.
type ApplicationDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// Label selector of the application.
	Selector map[string]string `json:"selector"`

	// Specs of the application.
	Descriptor          applicationApi.Descriptor `json:"descriptor"`
	Info                []applicationApi.InfoItem `json:"info"`
	ComponentGroupKinds []metav1.GroupKind        `json:"componentGroupKinds"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// GetApplicationDetail returns model object of application and error, if any.
func GetApplicationDetail(client applicationAlphaClient.Interface, namespace string,
	applicationName string) (*ApplicationDetail, error) {

	log.Printf("Getting details of %s application in %s namespace", applicationName, namespace)

	application, err := client.AppV1alpha1().Applications(namespace).Get(applicationName, metav1.GetOptions{})

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return &ApplicationDetail{
		ObjectMeta:          api.NewObjectMeta(application.ObjectMeta),
		TypeMeta:            api.NewTypeMeta(api.ResourceKindApplication),
		Descriptor:          application.Spec.Descriptor,
		Selector:            application.Spec.Selector.MatchLabels,
		Info:                application.Spec.Info,
		ComponentGroupKinds: application.Spec.ComponentGroupKinds,
		Errors:              nonCriticalErrors,
	}, nil
}
