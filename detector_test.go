package otelservicedetector_test

import (
	"context"
	"errors"
	"testing"

	otelservicedetector "github.com/aereal/otel-service-detector"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const (
	svcName    = "svc"
	svcVersion = "latest"
	env        = "production"
)

var (
	wantResource = resource.NewSchemaless(
		semconv.ServiceName(svcName),
		semconv.ServiceVersion(svcVersion),
		semconv.DeploymentEnvironment(env),
	)
	partial     = resource.NewSchemaless(semconv.ServiceVersion(svcVersion))
	stdMetadata = &otelservicedetector.ContainerMetadata{
		Image: "app:latest",
		Labels: otelservicedetector.LabelsMap{
			"otel.service.name":           svcName,
			"otel.deployment.environment": env,
		},
	}
	customMetadata = &otelservicedetector.ContainerMetadata{
		Image: "app:latest",
		Labels: otelservicedetector.LabelsMap{
			"custom.service.name":           svcName,
			"custom.deployment.environment": env,
		},
	}
)

var errOops = errors.New("oops")

func TestDetector_Detect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, deadline)
	}
	defer cancel()

	type testCase struct {
		name    string
		opts    []otelservicedetector.Option
		want    *resource.Resource
		wantErr error
	}
	testCases := []testCase{
		{
			"ok",
			[]otelservicedetector.Option{
				otelservicedetector.WithContainerMetadataProvider(otelservicedetector.StaticContainerMetadataProvider{stdMetadata}),
			},
			wantResource,
			nil,
		},
		{
			"ok/precendence",
			[]otelservicedetector.Option{
				otelservicedetector.WithContainerMetadataProvider(otelservicedetector.StaticContainerMetadataProvider{customMetadata}),
				otelservicedetector.WithServiceNameLabel("custom.service.name"),
				otelservicedetector.WithDeploymentEnvironmentLabel("custom.deployment.environment"),
			},
			wantResource,
			nil,
		},
		{
			"no provider given",
			[]otelservicedetector.Option{},
			nil,
			nil,
		},
		{
			"no matched labels found",
			[]otelservicedetector.Option{
				otelservicedetector.WithContainerMetadataProvider(otelservicedetector.StaticContainerMetadataProvider{customMetadata}),
			},
			partial,
			nil,
		},
		{
			"container label provider failed",
			[]otelservicedetector.Option{
				otelservicedetector.WithContainerMetadataProvider(otelservicedetector.ContainerMetadataProviderFunc(func(context.Context) (*otelservicedetector.ContainerMetadata, error) { return nil, errOops })),
			},
			nil,
			errOops,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := otelservicedetector.New(tc.opts...)
			res, err := resource.New(ctx, resource.WithDetectors(d))
			if !res.Equal(tc.want) {
				t.Errorf("resource: want=%v got=%v", tc.want, res)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("error: want=%v got=%v", tc.wantErr, err)
			}
		})
	}
}
