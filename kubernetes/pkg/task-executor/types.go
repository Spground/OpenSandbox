// Copyright 2025 Alibaba Group Holding Ltd.
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

package task_executor

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Task represents the internal local task resource (LocalTask)
// It follows the Kubernetes resource model with Metadata, Spec, and Status.
type Task struct {
	Name              string       `json:"name"`
	DeletionTimestamp *metav1.Time `json:"deletionTimestamp,omitempty"`

	// Spec defines the desired behavior of the task.
	// We reuse the v1alpha1.TaskSpec to ensure consistency with the controller API.
	Spec TaskSpec `json:"spec"`

	// Status describes the current state of the task.
	// We reuse the v1alpha1.TaskStatus to ensure consistency with the controller API.
	Status TaskStatus `json:"status"`
}

type TaskSpec struct {
	Process *Process
	Pod     *corev1.PodTemplate
}

type Process struct {
	// Command command
	Command []string `json:"command"`
	// Arguments to the entrypoint.
	Args []string `json:"args,omitempty"`
	// List of environment variables to set in the task.
	Env []corev1.EnvVar `json:"env,omitempty"`
	// WorkingDir task working directory.
	WorkingDir string `json:"workingDir,omitempty"`
}

type TaskStatus struct {
	// Details about the task's current condition.
	// +optional
	State TaskState `json:"state,omitempty"`
}

// TaskState holds a possible state of task.
// Only one of its members may be specified.
// If none of them is specified, the default one is TaskStateWaiting.
type TaskState struct {
	// Details about a waiting task
	// +optional
	Waiting *TaskStateWaiting `json:"waiting,omitempty"`
	// Details about a running task
	// +optional
	Running *TaskStateRunning `json:"running,omitempty"`
	// Details about a terminated task
	// +optional
	Terminated *TaskStateTerminated `json:"terminated,omitempty"`
}

// TaskStateWaiting is a waiting state of a task.
type TaskStateWaiting struct {
	// (brief) reason the task is not yet running.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message regarding why the task is not yet running.
	// +optional
	Message string `json:"message,omitempty"`
}

// TaskStateRunning is a running state of a task.
type TaskStateRunning struct {
	// Time at which the task was last (re-)started
	// +optional
	StartedAt metav1.Time `json:"startedAt,omitempty"`
}

// TaskStateTerminated is a terminated state of a task.
type TaskStateTerminated struct {
	// Exit status from the last termination of the task
	ExitCode int32 `json:"exitCode"`
	// Signal from the last termination of the task
	// +optional
	Signal int32 `json:"signal,omitempty"`
	// (brief) reason from the last termination of the task
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message regarding the last termination of the task
	// +optional
	Message string `json:"message,omitempty"`
	// Time at which previous execution of the task started
	// +optional
	StartedAt metav1.Time `json:"startedAt,omitempty"`
	// Time at which the task last terminated
	// +optional
	FinishedAt metav1.Time `json:"finishedAt,omitempty"`
}
