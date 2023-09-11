package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"

	"github.com/x95castle1/probes-convention-service/pkg/handler"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/vmware-tanzu/cartographer-conventions/webhook"
	"go.uber.org/zap"
)

const (
	logComponentKey  = "component"
	logComponentName = "simple-conventions"
)

func main() {
	ctx := context.Background()
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

	logger.Info("Convention server starting on: %v ...", port)

	// Setting an Anonymous Function to call the handler.AddConventions
	c := func(template *corev1.PodTemplateSpec, images []webhook.ImageConfig) ([]string, error) {
		return handler.AddConventions(logger, template, images)
	}

	// Create a listener on the / root to call. This is the real HTTPServer.
	// This is why you have port defined above.
	http.HandleFunc("/", webhook.ConventionHandler(ctx, c))

	logger.Fatal(webhook.NewConventionServer(ctx, fmt.Sprintf(":%s", port)))
}
