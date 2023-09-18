package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"

	"github.com/x95castle1/convention-server-framework/pkg/handler"

	"github.com/x95castle1/multi-purpose-convention-server/pkg/conventions"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/vmware-tanzu/cartographer-conventions/webhook"
	"go.uber.org/zap"
)

const (
	logComponentKey  = "component"
	logComponentName = "multi-purpose-convention"
)

func main() {
	ctx := context.Background()
	annotationPrefix := os.Getenv("ANNOTATION_PREFIX")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Setup logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	l, err := config.Build()
	if err != nil {
		os.Exit(1)
	}

	logger := l.Sugar().With(logComponentKey, logComponentName)
	ctx = logr.NewContext(ctx, zapr.NewLogger(l))

	logger.Info("Convention server starting on port: ", port, " with prefix: ", annotationPrefix)

	// Setting an Anonymous Function to call the handler.AddConventions
	c := func(template *corev1.PodTemplateSpec, images []webhook.ImageConfig) ([]string, error) {
		return handler.AddConventions(logger, template, images, conventions.Conventions)
	}

	// Create a listener on the / root to call. This is the real HTTPServer.
	// This is why you have port defined above.
	http.HandleFunc("/", webhook.ConventionHandler(ctx, c))

	logger.Fatal(webhook.NewConventionServer(ctx, fmt.Sprintf(":%s", port)))
}
