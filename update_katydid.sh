set -xe
(cd implementations/go-katydid-auto-reflect && GOPROXY=direct go get github.com/katydid/validator-go-jsonschema@main && go mod tidy)
(cd implementations/go-katydid-auto-json && GOPROXY=direct go get github.com/katydid/validator-go-jsonschema@main && go mod tidy)
(cd implementations/go-katydid-mem-reflect && GOPROXY=direct go get github.com/katydid/validator-go-jsonschema@main && go mod tidy)
(cd implementations/go-katydid-mem-json && GOPROXY=direct go get github.com/katydid/validator-go-jsonschema@main && go mod tidy)
