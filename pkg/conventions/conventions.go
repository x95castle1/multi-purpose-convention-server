package conventions

import (
	"context"
	"encoding/json"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"

	"github.com/x95castle1/convention-server-framework/pkg/convention"
)

var Prefix = os.Getenv("ANNOTATION_PREFIX")
var ReadinessId = Prefix + "-readiness"
var ReadinessAnnotation = Prefix + "/readinessProbe"
var LivenessId = Prefix + "-liveness"
var LivenessAnnotation = Prefix + "/livenessProbe"
var StartupId = Prefix + "-startup"
var StartUpAnnotation = Prefix + "/startupProbe"
var ArgsId = Prefix + "-args"
var ArgsAnnotation = Prefix + "/args"
var StorageId = Prefix + "-storage"
var StorageAnnotation = Prefix + "/storage"
var TolerationId = Prefix + "-tolerations"
var TolerationAnnotation = Prefix + "/tolerations"
var NodeSelectorId = Prefix + "-nodeSelector"
var NodeSelectorAnnotation = Prefix + "/nodeSelector"
var AffinityId = Prefix + "-affinity"
var AffinityAnnotation = Prefix + "/affinity"

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

				log.Printf("Adding Args %+v", a)
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

				log.Printf("Adding Volumes %+v", s.Volumes)
				target.Spec.Volumes = append(target.Spec.Volumes, s.Volumes...)
				log.Printf("Adding VolumeMounts %+v", s.VolumeMounts)
				c.VolumeMounts = append(c.VolumeMounts, s.VolumeMounts...)
			}
			return nil
		},
	},

	&convention.BasicConvention{
		Id: TolerationId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, TolerationAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			tolerations := getAnnotation(target, TolerationAnnotation)

			s, err := getTolerations(tolerations)
			if err != nil {
				return err
			}

			log.Printf("Adding Tolerations %+v", s)
			target.Spec.Tolerations = append(target.Spec.Tolerations, s...)

			return nil
		},
	},

	&convention.BasicConvention{
		Id: NodeSelectorId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, NodeSelectorAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			nodeSelector := getAnnotation(target, NodeSelectorAnnotation)

			parsedSelector, err := parseMap(nodeSelector)
			if err != nil {
				return err
			}

			log.Printf("Adding NodeSelector %+v", parsedSelector)
			target.Spec.NodeSelector = parsedSelector

			return nil
		},
	},

	&convention.BasicConvention{
		Id: AffinityId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, AffinityAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			affinity := getAnnotation(target, AffinityAnnotation)

			parsedAffinity, err := parseAffinity(affinity)
			if err != nil {
				return err
			}

			log.Printf("Adding Affinity %+v", parsedAffinity)
			target.Spec.Affinity = parsedAffinity

			return nil
		},
	},
}

func parseAffinity(arguments string) (*corev1.Affinity, error) {
	var a corev1.Affinity
	err := json.Unmarshal([]byte(arguments), &a)
	return &a, err
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
func getProbe(config string) (*corev1.Probe, error) {
	probe := corev1.Probe{}
	err := json.Unmarshal([]byte(config), &probe)
	return &probe, err
}

func getTolerations(config string) ([]corev1.Toleration, error) {
	tolerations := []corev1.Toleration{}
	err := json.Unmarshal([]byte(config), &tolerations)
	return tolerations, err
}

func parseMap(arguments string) (map[string]string, error) {
	var a map[string]string
	err := json.Unmarshal([]byte(arguments), &a)
	return a, err
}

type Storage struct {
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	Volumes      []corev1.Volume      `json:"volumes,omitempty"`
}
