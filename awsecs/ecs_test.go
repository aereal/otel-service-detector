package awsecs_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	otelservicedetector "github.com/aereal/otel-service-detector"
	"github.com/aereal/otel-service-detector/awsecs"
)

const (
	label1 = "name-1"
	label2 = "name-2"
)

func TestProvider_ProvideContainerMetadata(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, deadline)
	}
	defer cancel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"Image": "app:latest",
			"Labels": map[string]string{
				label1: "value-1",
				label2: "value-2",
			},
		})
	}))
	t.Cleanup(srv.Close)
	t.Setenv("ECS_CONTAINER_METADATA_URI_V4", srv.URL)
	p := awsecs.New()
	got, err := p.ProvideContainerMetadata(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := &otelservicedetector.ContainerMetadata{
		Image: "app:latest",
		Labels: otelservicedetector.LabelsMap{
			label1: "value-1",
			label2: "value-2",
		},
	}
	if got.Image != want.Image {
		t.Errorf("Image: want=%q got=%q", want.Image, got.Image)
	}
	cmpLabel(t, label1, want.Labels, got.Labels)
	cmpLabel(t, label2, want.Labels, got.Labels)
}

func cmpLabel(t *testing.T, name string, wantLabels, gotLabels otelservicedetector.Labels) {
	want := wantLabels.Get(name)
	got := gotLabels.Get(name)
	if want != got {
		t.Errorf("Labels[%s]: want=%q got=%q", name, want, got)
	}
}
