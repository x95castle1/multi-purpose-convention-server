/*
Copyright 2020 VMware Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"

	"github.com/garethjevans/simple-conventions/pkg/handler"

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

	c := func(template *corev1.PodTemplateSpec, images []webhook.ImageConfig) ([]string, error) {
		return handler.AddConventions(logger, template, images)
	}

	http.HandleFunc("/", webhook.ConventionHandler(ctx, c))

	logger.Fatal(webhook.NewConventionServer(ctx, fmt.Sprintf(":%s", port)))
}
