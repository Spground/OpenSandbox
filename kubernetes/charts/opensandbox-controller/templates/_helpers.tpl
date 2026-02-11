{{/*
Expand the name of the chart.
*/}}
{{- define "opensandbox.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "opensandbox.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "opensandbox.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "opensandbox.labels" -}}
helm.sh/chart: {{ include "opensandbox.chart" . }}
{{ include "opensandbox.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "opensandbox.selectorLabels" -}}
app.kubernetes.io/name: {{ include "opensandbox.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
control-plane: controller-manager
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "opensandbox.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (printf "%s-controller-manager" (include "opensandbox.fullname" .)) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get the namespace to use
*/}}
{{- define "opensandbox.namespace" -}}
{{- if .Values.namespaceOverride }}
{{- .Values.namespaceOverride }}
{{- else }}
{{- printf "%s-system" (include "opensandbox.name" .) }}
{{- end }}
{{- end }}

{{/*
Controller image
*/}}
{{- define "opensandbox.controllerImage" -}}
{{- $tag := .Values.controller.image.tag | default .Chart.AppVersion }}
{{- printf "%s:%s" .Values.controller.image.repository $tag }}
{{- end }}

{{/*
Create the name for the leader election role
*/}}
{{- define "opensandbox.leaderElectionRoleName" -}}
{{- printf "%s-leader-election-role" (include "opensandbox.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create the name for the manager role
*/}}
{{- define "opensandbox.managerRoleName" -}}
{{- printf "%s-manager-role" (include "opensandbox.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Return the appropriate apiVersion for RBAC APIs
*/}}
{{- define "opensandbox.rbac.apiVersion" -}}
{{- if .Capabilities.APIVersions.Has "rbac.authorization.k8s.io/v1" }}
{{- print "rbac.authorization.k8s.io/v1" }}
{{- else }}
{{- print "rbac.authorization.k8s.io/v1beta1" }}
{{- end }}
{{- end }}

{{/*
Return image pull policy
*/}}
{{- define "opensandbox.imagePullPolicy" -}}
{{- .Values.controller.image.pullPolicy | default "IfNotPresent" }}
{{- end }}
