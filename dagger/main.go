package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Depado/daggo/dagger/internal/dagger"
)

const GOVERSION = "1.22"
const BINARY = "daggo"

var gomodCache = dag.CacheVolume(fmt.Sprintf("go-mod-%s", GOVERSION))
var buildCache = dag.CacheVolume(fmt.Sprintf("go-build-%s", GOVERSION))

type Daggo struct{}

func (m *Daggo) BaseGoEnv(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("golang:%s-alpine", GOVERSION)).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", gomodCache).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", buildCache).
		WithEnvVariable("GOCACHE", "/go/build-cache").WithExec([]string{"go", "get"})
}

func (m *Daggo) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BaseGoEnv(source).
		WithExec([]string{"go", "test", "./..."}).
		Stdout(ctx)
}

func (m *Daggo) Build(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BaseGoEnv(source).
		WithExec([]string{"go", "build"}).
		Stdout(ctx)
}

func (m *Daggo) Docker(ctx context.Context, source *dagger.Directory) *dagger.Container {
	builder := m.BaseGoEnv(source).WithExec([]string{"go", "build", "-o", BINARY})
	binPath := fmt.Sprintf("/bin/%s", BINARY)
	return dag.Container().
		From("gcr.io/distroless/static").
		WithLabel("org.opencontainers.image.title", BINARY).
		WithLabel("org.opencontainers.image.version", "1.0").
		WithLabel("org.opencontainers.image.created", time.Now().String()).
		WithLabel("org.opencontainers.image.source", "https://github.com/Depado/daggo").
		WithLabel("org.opencontainers.image.licenses", "MIT").
		WithFile(binPath, builder.File(fmt.Sprintf("/src/%s", BINARY))).
		WithEntrypoint([]string{binPath, "serve", "--http=0.0.0.0:8080"})
}

func (m *Daggo) Publish(ctx context.Context, registry, imageName, username string, password *dagger.Secret, source *dagger.Directory) error {
	_, err := m.Docker(ctx, source).
		WithRegistryAuth(registry, username, password).
		Publish(ctx, registry+"/"+imageName)
	return err
}

func (m *Daggo) Run(ctx context.Context, source *dagger.Directory) (string, error) {
	if _, err := m.Build(ctx, source); err != nil {
		return "", err
	}
	if _, err := m.Test(ctx, source); err != nil {
		return "", err
	}

	return "", nil
}
