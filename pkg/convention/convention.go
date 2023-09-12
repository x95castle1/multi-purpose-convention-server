package convention

import (
	"context"

	"github.com/vmware-tanzu/cartographer-conventions/webhook"
	corev1 "k8s.io/api/core/v1"
)

// Convention is the interface defining contract of the conventions to be applied
// This is boilerplate code. There shouldn't be a need to change this.

type Convention interface {
	GetId() string
	IsApplicable(ctx context.Context, target *corev1.PodTemplateSpec, metadata ImageMetadata) bool
	ApplyConvention(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata ImageMetadata, imageName string) error
}

type ImageMetadata = map[string]webhook.ImageConfig

var _ Convention = (*BasicConvention)(nil)

// BasicConvention defines a basic convention to be applied
type BasicConvention struct {
	Id         string
	Applicable func(ctx context.Context, target *corev1.PodTemplateSpec, metadata ImageMetadata) bool
	Apply      func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata ImageMetadata, imageName string) error
}

// GetId retrieves the identifier of the convention
func (o *BasicConvention) GetId() string {
	return o.Id
}

// IsApplicable defines if a convention can and will be applied based on the image and container spec
func (o *BasicConvention) IsApplicable(ctx context.Context, target *corev1.PodTemplateSpec, metadata ImageMetadata) bool {
	if o.Applicable == nil {
		return true
	}
	return o.Applicable(ctx, target, metadata)
}

// ApplyConvention Aaplies the convention to the `PodTemplateSpec`
func (o *BasicConvention) ApplyConvention(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata ImageMetadata, imageName string) error {
	return o.Apply(ctx, target, containerIdx, metadata, imageName)
}
