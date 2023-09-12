package convention

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestBasicConvention_GetId(t *testing.T) {
	type fields struct {
		Id         string
		Applicable func(ctx context.Context, target *corev1.PodTemplateSpec, metadata ImageMetadata) bool
		Apply      func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata ImageMetadata, imageName string) error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "gets id",
			fields: fields{
				Id: "convention identification",
			},
			want: "convention identification",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &BasicConvention{
				Id:         tt.fields.Id,
				Applicable: tt.fields.Applicable,
				Apply:      tt.fields.Apply,
			}
			if got := o.GetId(); got != tt.want {
				t.Errorf("BasicConvention.GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicConvention_IsApplicable(t *testing.T) {
	type fields struct {
		Id         string
		Applicable func(ctx context.Context, target *corev1.PodTemplateSpec, metadata ImageMetadata) bool
		Apply      func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata ImageMetadata, imageName string) error
	}
	type args struct {
		ctx      context.Context
		metadata ImageMetadata
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "is applicable null",
			fields: fields{
				Id:         "applicable convention",
				Applicable: nil,
			},
			args: args{},
			want: true,
		},
		{
			name: "is applicable",
			fields: fields{
				Id:         "applicable convention",
				Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata ImageMetadata) bool { return true },
			},
			args: args{},
			want: true,
		},
		{
			name: "is NOT applicable",
			fields: fields{
				Id:         "not applicable convention",
				Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata ImageMetadata) bool { return false },
			},
			args: args{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &BasicConvention{
				Id:         tt.fields.Id,
				Applicable: tt.fields.Applicable,
				Apply:      tt.fields.Apply,
			}
			if got := o.IsApplicable(tt.args.ctx, nil, tt.args.metadata); got != tt.want {
				t.Errorf("BasicConvention.IsApplicable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicConvention_ApplyConvention(t *testing.T) {
	type fields struct {
		Id         string
		Applicable func(ctx context.Context, target *corev1.PodTemplateSpec, metadata ImageMetadata) bool
		Apply      func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata ImageMetadata, imageName string) error
	}
	type args struct {
		ctx          context.Context
		target       *corev1.PodTemplateSpec
		containerIdx int
		metadata     ImageMetadata
		imageName    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "no error appply",
			fields: fields{
				Id: "no-error-convention",
				Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata ImageMetadata, imageName string) error {
					return nil
				},
			},
			args:    args{},
			wantErr: false,
		},
		{
			name: "error on appply",
			fields: fields{
				Id: "error-convention",
				Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata ImageMetadata, imageName string) error {
					return fmt.Errorf("Error")
				},
			},
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &BasicConvention{
				Id:         tt.fields.Id,
				Applicable: tt.fields.Applicable,
				Apply:      tt.fields.Apply,
			}
			if err := o.ApplyConvention(tt.args.ctx, tt.args.target, tt.args.containerIdx, tt.args.metadata, tt.args.imageName); (err != nil) != tt.wantErr {
				t.Errorf("BasicConvention.ApplyConvention() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setAnnotation(t *testing.T) {
	p := &corev1.PodTemplateSpec{}
	p.Annotations = make(map[string]string)
	type args struct {
		pts   *corev1.PodTemplateSpec
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "no annotations exist",
			args: args{
				pts:   &corev1.PodTemplateSpec{},
				key:   "annotation-key",
				value: "annotation-value",
			},
		},
		{
			name: "no annotations exist",
			args: args{
				pts:   p,
				key:   "annotation-key",
				value: "annotation-value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setAnnotation(tt.args.pts, tt.args.key, tt.args.value)
			if _, ok := tt.args.pts.Annotations[tt.args.key]; !ok {
				t.Errorf("BasicConvention.setAnnotation() = %v, want %v", ok, tt.args.value)
			}
		})
	}
}

func Test_setLabel(t *testing.T) {
	p := &corev1.PodTemplateSpec{}
	p.Labels = make(map[string]string)
	type args struct {
		pts   *corev1.PodTemplateSpec
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "no annotations exist",
			args: args{
				pts:   &corev1.PodTemplateSpec{},
				key:   "annotation-key",
				value: "annotation-value",
			},
		},
		{
			name: "no annotations exist",
			args: args{
				pts:   p,
				key:   "annotation-key",
				value: "annotation-value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setLabel(tt.args.pts, tt.args.key, tt.args.value)
			if _, ok := tt.args.pts.Labels[tt.args.key]; !ok {
				t.Errorf("BasicConvention.setLabel() = %v, want %v", ok, tt.args.value)
			}
		})
	}
}

func Test_findEnvVar(t *testing.T) {
	c := corev1.Container{
		Env: []corev1.EnvVar{
			{
				Name:  "env-var1",
				Value: "value1",
			},
			{
				Name:  "env-var2",
				Value: "value2",
			},
		},
	}
	type args struct {
		container corev1.Container
		name      string
	}
	tests := []struct {
		name string
		args args
		want *corev1.EnvVar
	}{
		{
			name: "env var exist",
			args: args{
				container: c,
				name:      "env-var2",
			},
			want: &corev1.EnvVar{
				Name:  "env-var2",
				Value: "value2",
			},
		},
		{
			name: "env var not exist",
			args: args{
				container: c,
				name:      "env-var3",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findEnvVar(tt.args.container, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findContainerPort(t *testing.T) {
	ps := corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name: "container1",
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: 1,
					},
					{
						ContainerPort: 2,
					},
				},
			},
			{
				Name: "container2",
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: 3,
					},
					{
						ContainerPort: 4,
					},
				},
			},
		},
	}
	type args struct {
		ps   corev1.PodSpec
		port int32
	}
	tests := []struct {
		name              string
		args              args
		wantName          string
		wantContainerPort *corev1.ContainerPort
	}{
		{
			name: "port found",
			args: args{
				ps:   ps,
				port: 4,
			},
			wantName: "container2",
			wantContainerPort: &corev1.ContainerPort{
				ContainerPort: 4,
			},
		},
		{
			name: "port not found",
			args: args{
				ps:   ps,
				port: 5,
			},
			wantName:          "",
			wantContainerPort: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findContainerPort(tt.args.ps, tt.args.port)
			if got != tt.wantName {
				t.Errorf("findContainerPort() got = %v, want %v", got, tt.wantName)
			}
			if !reflect.DeepEqual(got1, tt.wantContainerPort) {
				t.Errorf("findContainerPort() got1 = %v, want %v", got1, tt.wantContainerPort)
			}
		})
	}
}

// setAnnotation sets the annotation on PodTemplateSpec
func setAnnotation(pts *corev1.PodTemplateSpec, key, value string) {
	if pts.Annotations == nil {
		pts.Annotations = map[string]string{}
	}
	pts.Annotations[key] = value
}

// setLabel sets the label on PodTemplateSpec
func setLabel(pts *corev1.PodTemplateSpec, key, value string) {
	if pts.Labels == nil {
		pts.Labels = map[string]string{}
	}
	pts.Labels[key] = value
}

// getLabel gets the label on PodTemplateSpec
func getLabel(pts *corev1.PodTemplateSpec, key string) string {
	if pts.Labels == nil || len(pts.Labels[key]) == 0 {
		return ""
	}
	return pts.Labels[key]
}

// findEnvVar search the `name` environment variable in the given `Container`
func findEnvVar(container corev1.Container, name string) *corev1.EnvVar {
	for i := range container.Env {
		e := &container.Env[i]
		if e.Name == name {
			return e
		}
	}
	return nil
}

// findContainerPort check if any `container` in the `PodSpec` has the given `port` and retrives its name and `ContainerPort`
// otherwise return empty string and nill for `ContainerPort`
func findContainerPort(ps corev1.PodSpec, port int32) (string, *corev1.ContainerPort) {
	for _, c := range ps.Containers {
		for _, p := range c.Ports {
			if p.ContainerPort == port {
				return c.Name, &p
			}
		}
	}
	return "", nil
}
