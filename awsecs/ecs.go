package awsecs // import "github.com/aereal/otel-service-detector/awsecs"

import (
	"context"
	"net/http"

	otelservicedetector "github.com/aereal/otel-service-detector"
	ecsmetadata "github.com/brunoscheufler/aws-ecs-metadata-go"
)

type Option func(*Provider)

func WithHTTPClient(client *http.Client) Option { return func(p *Provider) { p.client = client } }

func New(opts ...Option) *Provider {
	p := new(Provider)
	for _, o := range opts {
		o(p)
	}
	if p.client == nil {
		p.client = http.DefaultClient
	}
	return p
}

type Provider struct {
	client *http.Client
}

func (p *Provider) ProvideContainerMetadata(ctx context.Context) (*otelservicedetector.ContainerMetadata, error) {
	md, err := ecsmetadata.GetContainerV4(ctx, p.client)
	if err != nil {
		return nil, err
	}
	return &otelservicedetector.ContainerMetadata{
		Image:  md.Image,
		Labels: md.Labels,
	}, nil
}
