# github.com/aereal/otel-service-detector/awsecs

Package awsecs provides the container metadata using [ECS task metadata endpoint][].

## Synopsis

```go
import (
	otelservicedetector "github.com/aereal/otel-service-detector"
	"github.com/aereal/otel-service-detector/awsecs"
	"go.opentelemetry.io/otel/sdk/resource"
)

func _() {
  detector := otelservicedetector.New(otelservicedetector.WithContainerMetadataProvider(awsecs.New()))
  _, _ = resource.New(context.TODO(), resource.WithDetectors(detector))
}
```

[ECS task metadata endpoint]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html
