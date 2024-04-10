[![status][ci-status-badge]][ci-status]
[![PkgGoDev][pkg-go-dev-badge]][pkg-go-dev]

# otel-service-detector

otel-service-detector provides OpenTelemetry resource detector that detects service attributes from the container metadata.

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

## Installation

```sh
go get github.com/aereal/otel-service-detector
```

## License

See LICENSE file.

[pkg-go-dev]: https://pkg.go.dev/github.com/aereal/otel-service-detector
[pkg-go-dev-badge]: https://pkg.go.dev/badge/aereal/otel-service-detector
[ci-status-badge]: https://github.com/aereal/otel-service-detector/workflows/CI/badge.svg?branch=main
[ci-status]: https://github.com/aereal/otel-service-detector/actions/workflows/CI
