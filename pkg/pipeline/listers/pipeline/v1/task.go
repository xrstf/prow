/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// TaskLister helps list Tasks.
// All objects returned here must be treated as read-only.
type TaskLister interface {
	// List lists all Tasks in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Task, err error)
	// Tasks returns an object that can list and get Tasks.
	Tasks(namespace string) TaskNamespaceLister
	TaskListerExpansion
}

// taskLister implements the TaskLister interface.
type taskLister struct {
	listers.ResourceIndexer[*v1.Task]
}

// NewTaskLister returns a new TaskLister.
func NewTaskLister(indexer cache.Indexer) TaskLister {
	return &taskLister{listers.New[*v1.Task](indexer, v1.Resource("task"))}
}

// Tasks returns an object that can list and get Tasks.
func (s *taskLister) Tasks(namespace string) TaskNamespaceLister {
	return taskNamespaceLister{listers.NewNamespaced[*v1.Task](s.ResourceIndexer, namespace)}
}

// TaskNamespaceLister helps list and get Tasks.
// All objects returned here must be treated as read-only.
type TaskNamespaceLister interface {
	// List lists all Tasks in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Task, err error)
	// Get retrieves the Task from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Task, error)
	TaskNamespaceListerExpansion
}

// taskNamespaceLister implements the TaskNamespaceLister
// interface.
type taskNamespaceLister struct {
	listers.ResourceIndexer[*v1.Task]
}
