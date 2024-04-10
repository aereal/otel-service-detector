package otelservicedetector // import otelservicedetector "github.com/aereal/otel-service-detector"

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type Labels interface {
	Get(key string) string
}

type LabelsMap map[string]string

func (m LabelsMap) Get(key string) string {
	return m[key]
}

type ContainerMetadata struct {
	Image  string
	Labels Labels
}

// ContainerMetadataProvider provides container's metadata.
//
// It may communicate with remote endpoints for collecting container's metadata.
type ContainerMetadataProvider interface {
	ProvideContainerMetadata(context.Context) (*ContainerMetadata, error)
}

type ContainerMetadataProviderFunc func(ctx context.Context) (*ContainerMetadata, error)

func (f ContainerMetadataProviderFunc) ProvideContainerMetadata(ctx context.Context) (*ContainerMetadata, error) {
	return f(ctx)
}

type StaticContainerMetadataProvider struct{ *ContainerMetadata }

func (p StaticContainerMetadataProvider) ProvideContainerMetadata(context.Context) (*ContainerMetadata, error) {
	return p.ContainerMetadata, nil
}

type Option func(d *Detector)

func WithContainerMetadataProvider(p ContainerMetadataProvider) Option {
	return func(d *Detector) { d.provider = p }
}

func WithServiceNameLabel(labels ...string) Option {
	return func(d *Detector) { d.labelServiceName = append(d.labelServiceName, labels...) }
}

func WithServiceVersionLabel(labels ...string) Option {
	return func(d *Detector) { d.labelServiceVersion = append(d.labelServiceVersion, labels...) }
}

func WithDeploymentEnvironmentLabel(labels ...string) Option {
	return func(d *Detector) { d.labelDeploymentEnv = append(d.labelDeploymentEnv, labels...) }
}

func New(opts ...Option) *Detector {
	d := &Detector{}
	for _, o := range opts {
		o(d)
	}
	d.labelServiceName = append(d.labelServiceName, "otel.service.name")
	d.labelServiceVersion = append(d.labelServiceVersion, "otel.service.version")
	d.labelDeploymentEnv = append(d.labelDeploymentEnv, "otel.deployment.environment")
	return d
}

type Detector struct {
	provider                                                  ContainerMetadataProvider
	labelServiceName, labelServiceVersion, labelDeploymentEnv []string
}

var _ resource.Detector = (*Detector)(nil)

func (d *Detector) Detect(ctx context.Context) (*resource.Resource, error) {
	if d == nil || d.provider == nil {
		return resource.Empty(), nil
	}
	metadata, err := d.provider.ProvideContainerMetadata(ctx)
	if err != nil {
		return resource.Empty(), &ProvideContainerMetadataError{Err: err}
	}
	attrs := make([]attribute.KeyValue, 0, 3)
	if _, tag, found := strings.Cut(metadata.Image, ":"); found {
		attrs = append(attrs, semconv.ServiceVersion(tag))
	}
	_ = d.add(&attrs, metadata.Labels, d.labelServiceName, semconv.ServiceNameKey)
	_ = d.add(&attrs, metadata.Labels, d.labelDeploymentEnv, semconv.DeploymentEnvironmentKey)
	return resource.NewWithAttributes(semconv.SchemaURL, attrs...), nil
}

func (d *Detector) add(attrs *[]attribute.KeyValue, labels Labels, labelNames []string, key attribute.Key) (found bool) {
	for _, labelName := range labelNames {
		v := labels.Get(labelName)
		if v == "" {
			continue
		}
		*attrs = append(*attrs, key.String(v))
		return true
	}
	return false
}

type ProvideContainerMetadataError struct {
	Err error
}

func (e *ProvideContainerMetadataError) Error() string {
	return fmt.Sprintf("failed to provide metadata: %s", e.Err)
}
func (e *ProvideContainerMetadataError) Unwrap() error { return e.Err }
