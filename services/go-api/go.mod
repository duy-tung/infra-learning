module go-api

go 1.23.8

require (
	github.com/gin-gonic/gin v1.10.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/clickhouseexporter v0.128.0
	github.com/prometheus/client_golang v1.20.3
	github.com/rs/zerolog v1.34.0
	go.opentelemetry.io/collector/component v1.34.0
	go.opentelemetry.io/collector/config/configopaque v1.34.0
	go.opentelemetry.io/collector/exporter v0.128.0
	go.opentelemetry.io/collector/pdata v1.34.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.61.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0
	go.opentelemetry.io/otel v1.36.0
	go.opentelemetry.io/otel/sdk v1.36.0
	go.opentelemetry.io/otel/trace v1.36.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/ClickHouse/ch-go v0.66.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.36.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/sonic v1.13.2 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.26.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/knadh/koanf/providers/confmap v1.0.0 // indirect
	github.com/knadh/koanf/v2 v2.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.128.0 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/collector/config/configretry v1.34.0 // indirect
	go.opentelemetry.io/collector/confmap v1.34.0 // indirect
	go.opentelemetry.io/collector/consumer v1.34.0 // indirect
	go.opentelemetry.io/collector/consumer/consumererror v0.128.0 // indirect
	go.opentelemetry.io/collector/extension v1.34.0 // indirect
	go.opentelemetry.io/collector/extension/xextension v0.128.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.34.0 // indirect
	go.opentelemetry.io/collector/internal/telemetry v0.128.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.128.0 // indirect
	go.opentelemetry.io/collector/pipeline v0.128.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.11.0 // indirect
	go.opentelemetry.io/otel/log v0.12.2 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.17.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/grpc v1.72.2 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
