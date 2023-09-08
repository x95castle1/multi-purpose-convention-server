package handler

import (
	"context"

	"github.com/vmware-tanzu/cartographer-conventions/webhook"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	"github.com/garethjevans/simple-conventions/pkg/resources"
)

func AddConventions(logger *zap.SugaredLogger, template *corev1.PodTemplateSpec, images []webhook.ImageConfig) ([]string, error) {
	imageMap := make(map[string]webhook.ImageConfig)
	for _, config := range images {
		imageMap[config.Image] = config
	}

	var appliedConventions []string
	for i := range template.Spec.Containers {
		container := &template.Spec.Containers[i]
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
