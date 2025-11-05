{{- define "system-education.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "system-education.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "system-education.chart" -}}
{{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end -}}

{{- define "system-education.labels" -}}
helm.sh/chart: {{ include "system-education.chart" . }}
{{ include "system-education.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "system-education.selectorLabels" -}}
app.kubernetes.io/name: {{ include "system-education.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
