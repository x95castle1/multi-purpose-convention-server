package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	webhookv1alpha1 "github.com/vmware-tanzu/cartographer-conventions/webhook/api/v1alpha1"

	"github.com/google/go-cmp/cmp"
	"github.com/vmware-tanzu/cartographer-conventions/webhook"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/x95castle1/probes-convention-service/pkg/resources"
)

const imageDefault = "sample-accelerators/tanzu-java-web-app"

func Test_addConventions(t *testing.T) {
	testdataPath := "testdata"
	var emptyAppliedConventions []string
	l := zap.NewNop().Sugar()
	type args struct {
		logger   *zap.SugaredLogger
		template *corev1.PodTemplateSpec
		images   []webhook.ImageConfig
	}
	tests := []struct {
		skip               bool
		name               string
		args               args
		want               []string
		wantErr            bool
		validateTemplate   bool
		wantedTemplateFile string
	}{
		{
			name: "no images",
			args: args{
				logger:   l,
				template: getMockTemplate(""),
				images:   make([]webhook.ImageConfig, 0),
			},
			want:    emptyAppliedConventions,
			wantErr: false,
		},
		{
			name: "readinessProbe",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/readinessProbe", resources.Prefix), "{\"exec\":{\"command\":[\"cat\",\"/tmp/healthy\"]},\"initialDelaySeconds\":5,\"periodSeconds\":5}"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-readiness", resources.Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "readinessProbe.json",
		},
		{
			name: "livenessProbe",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/livenessProbe", resources.Prefix), "{\"exec\":{\"command\":[\"cat\",\"/tmp/healthy\"]},\"initialDelaySeconds\":5,\"periodSeconds\":5}"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-liveness", resources.Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "livenessProbe.json",
		},
		{
			name: "startupProbe",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/startupProbe", resources.Prefix), "{\"exec\":{\"command\":[\"cat\",\"/tmp/healthy\"]},\"initialDelaySeconds\":5,\"periodSeconds\":5}"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-startup", resources.Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "startupProbe.json",
		},
		{
			name: "cartoRunWorkloadName",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndLabel("", "carto.run/workload-name", "my-workload"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-carto-run-workload-name", resources.Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "cartoRunWorkloadName.json",
		},
		{
			name: "args",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/args", resources.Prefix), "[\"one\",\"two\",\"three\"]"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-args", resources.Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "args.json",
		},
		{
			name: "storage",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/storage", resources.Prefix), "{\"volumeMounts\":[{\"mountPath\":\"/test\",\"name\":\"test\"}],\"volumes\":[{\"name\":\"test\",\"emptyDir\":{}}]}"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-storage", resources.Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "storage.json",
		},
	}

	for _, tt := range tests {
		if tt.skip {
			t.Logf("skipping test %s", tt.name)
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddConventions(tt.args.logger, tt.args.template, tt.args.images)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddConventions() = unwanted error: %v", err)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("%s() = (-expected, +actual): %s", tt.name, diff)
			}

			if tt.validateTemplate {
				if tt.wantedTemplateFile == "" {
					t.Errorf("AddConventions(): [%s] Template JSON file not provided", tt.name)
				} else {
					pts := &corev1.PodTemplateSpec{}
					tf := getFileBytes(testdataPath + "/templates/" + tt.wantedTemplateFile)

					err = json.Unmarshal(tf, pts)
					if err != nil {
						t.Errorf("AddConventions(): [%s] %v", tt.name, err)
					}
					if diff := cmp.Diff(tt.args.template, pts); diff != "" {
						t.Errorf("%s() = (-expected, +actual): %s", tt.name, diff)
					}
				}
			}
		})
	}
}

func getFileBytes(file string) []byte {
	b, _ := os.ReadFile(file)
	return b
}

func getMockTemplate(img string) *corev1.PodTemplateSpec {
	if img == "" {
		img = imageDefault
	}
	return &corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "workload",
					Image: img,
					Ports: make([]corev1.ContainerPort, 0),
				},
			},
		},
	}
}

func getMockTemplateWithImageAndAnnotation(img string, key string, value string) *corev1.PodTemplateSpec {
	if img == "" {
		img = imageDefault
	}
	return &corev1.PodTemplateSpec{
		ObjectMeta: v1.ObjectMeta{Annotations: map[string]string{key: value}},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "workload",
					Image: img,
					Ports: make([]corev1.ContainerPort, 0),
				},
			},
		},
	}
}

func getMockTemplateWithImageAndLabel(img string, key string, value string) *corev1.PodTemplateSpec {
	if img == "" {
		img = imageDefault
	}
	return &corev1.PodTemplateSpec{
		ObjectMeta: v1.ObjectMeta{Labels: map[string]string{key: value}},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "workload",
					Image: img,
					Ports: make([]corev1.ContainerPort, 0),
				},
			},
		},
	}
}
