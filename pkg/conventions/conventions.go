package conventions

import (
	"context"
	"encoding/json"
	"log"

	corev1 "k8s.io/api/core/v1"

	"github.com/x95castle1/convention-server-framework/pkg/convention"
)

const (
	Prefix                  = "x95castle1.org"
	ReadinessId             = Prefix + "-readiness"
	ReadinessAnnotation     = Prefix + "/readinessProbe"
	LivenessId              = Prefix + "-liveness"
	LivenessAnnotation      = Prefix + "/livenessProbe"
	StartupId               = Prefix + "-startup"
	StartUpAnnotation       = Prefix + "/startupProbe"
	ArgsId                  = Prefix + "-args"
	ArgsAnnotation          = Prefix + "/args"
	StorageId               = Prefix + "-storage"
	StorageAnnotation       = Prefix + "/storage"
	TolerationId            = Prefix + "-toleration"
	TolerationAnnotation    = Prefix + "/toleration"
	WorkloadNameId          = Prefix + "-carto-run-workload-name"
	WorkloadNameLabel       = "carto.run/workload-name"
	WorkloadNameEnvVariable = "CARTO_RUN_WORKLOAD_NAME"
)

var Conventions = []convention.Convention{
	&convention.BasicConvention{
		Id: ReadinessId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, ReadinessAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			readinessProbe := getAnnotation(target, ReadinessAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				if c.ReadinessProbe == nil {
					p, err := getProbe(readinessProbe)
					if err != nil {
						return err
					}
					log.Printf("Adding ReadinessProbe %+v", p)
					c.ReadinessProbe = p
				}
			}
			return nil
		},
	},

	&convention.BasicConvention{
		Id: LivenessId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, LivenessAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			livenessProbe := getAnnotation(target, LivenessAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				if c.LivenessProbe == nil {
					p, err := getProbe(livenessProbe)
					if err != nil {
						return err
					}
					log.Printf("Adding LivenessProbe %+v", p)
					c.LivenessProbe = p
				}
			}
			return nil
		},
	},

	&convention.BasicConvention{
		Id: StartupId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, StartUpAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			startupProbe := getAnnotation(target, StartUpAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				if c.StartupProbe == nil {
					p, err := getProbe(startupProbe)
					if err != nil {
						return err
					}
					log.Printf("Adding StartupProbe %+v", p)
					c.StartupProbe = p
				}
			}
			return nil
		},
	},

	&convention.BasicConvention{
		Id: WorkloadNameId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getLabel(target, WorkloadNameLabel) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			value := getLabel(target, WorkloadNameLabel)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]
				addEnvVar(c, corev1.EnvVar{
					Name:  WorkloadNameEnvVariable,
					Value: value,
				})
			}

			return nil
		},
	},

	&convention.BasicConvention{
		Id: ArgsId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, ArgsAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			arguments := getAnnotation(target, ArgsAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				a, err := getArguments(arguments)
				if err != nil {
					return err
				}

				c.Args = append(c.Args, a...)
			}
			return nil
		},
	},

	&convention.BasicConvention{
		Id: StorageId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, StorageAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			storage := getAnnotation(target, StorageAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				s, err := getStorage(storage)
				if err != nil {
					return err
				}

				c.VolumeMounts = append(c.VolumeMounts, s.VolumeMounts...)
				target.Spec.Volumes = append(target.Spec.Volumes, s.Volumes...)
			}
			return nil
		},
	},
}

// getArguments parse the arguments into a string array
func getArguments(arguments string) ([]string, error) {
	var a []string
	err := json.Unmarshal([]byte(arguments), &a)
	return a, err
}

// getStorage parse the arguments into a storage struct
func getStorage(arguments string) (Storage, error) {
	var s Storage
	err := json.Unmarshal([]byte(arguments), &s)
	return s, err
}

// getAnnotation gets the annotation on PodTemplateSpec
func getAnnotation(pts *corev1.PodTemplateSpec, key string) string {
	if pts.Annotations == nil || len(pts.Annotations[key]) == 0 {
		return ""
	}
	return pts.Annotations[key]
}

// getLabel gets the label on PodTemplateSpec
func getLabel(pts *corev1.PodTemplateSpec, key string) string {
	if pts.Labels == nil || len(pts.Labels[key]) == 0 {
		return ""
	}
	return pts.Labels[key]
}

func getProbe(config string) (*corev1.Probe, error) {
	probe := corev1.Probe{}
	err := json.Unmarshal([]byte(config), &probe)
	return &probe, err
}

func addEnvVar(container *corev1.Container, envvar corev1.EnvVar) bool {
	for _, e := range container.Env {
		if e.Name == envvar.Name {
			return false
		}
	}
	container.Env = append(container.Env, envvar)
	return true
}

type Storage struct {
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	Volumes      []corev1.Volume      `json:"volumes,omitempty"`
}
