package handler

import (
	"context"

	"github.com/vmware-tanzu/cartographer-conventions/webhook"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	"github.com/x95castle1/probes-convention-service/pkg/resources"
)

// This is more framework that wraps around true Convention. True logic
// is in the resources folder.
// This is boilerplate code as well.

func AddConventions(logger *zap.SugaredLogger, template *corev1.PodTemplateSpec, images []webhook.ImageConfig) ([]string, error) {

	// These come from the PodConventionContextSpec
	// These are the all the images
	imageMap := make(map[string]webhook.ImageConfig)

	for _, config := range images {
		imageMap[config.Image] = config
	}

	var appliedConventions []string

	//Loop through every container defined on the PodTemplateSpec from the PodIntent
	for i := range template.Spec.Containers {
		container := &template.Spec.Containers[i]

		// This is all the images that are part of the containers
		image, ok := imageMap[container.Image]
		if !ok {
			logger.Warnw("image name not defined", "container", container.Name)
			continue
		}

		ctx := context.Background()

		imageName := image.Config.Config.Labels["org.opencontainers.image.title"]

		for _, o := range resources.Conventions {
			if !o.IsApplicable(ctx, template, imageMap) {
				continue
			}
			if err := o.ApplyConvention(ctx, template, i, imageMap, imageName); err != nil {
				logger.Errorw(err.Error(), "convention", o.GetId(), "namespace", template.Namespace, "name", template.Name, "kind", "PodTemplateSpec")
				return nil, err
			}
			appliedConventions = append(appliedConventions, o.GetId())

			logger.Infow("Successfully applied convention", "convention", o.GetId(), "namespace", template.Namespace, "name", template.Name, "kind", "PodTemplateSpec")
		}
	}
	return appliedConventions, nil
}
