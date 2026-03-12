package container

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"

	"github.com/base-go/backend/internal/auth"
	"github.com/base-go/backend/internal/rbac"
	"github.com/base-go/backend/pkg/cache"
	"github.com/base-go/backend/pkg/database"
	"github.com/base-go/backend/pkg/router"
	"github.com/base-go/backend/pkg/server"
)

func New() (*dig.Container, error) {
	// we use go.uber.org/dig for autowire dependencies
	container := dig.New()

	// provide dependencies injection

	if err := container.Provide(cache.NewCache); err != nil {
		return nil, err
	}

	// database
	if err := container.Provide(database.NewDatabase); err != nil {
		return nil, err
	}

	// auth module
	if err := container.Provide(auth.NewRepository); err != nil {
		return nil, err
	}

	if err := container.Provide(auth.NewService); err != nil {
		return nil, err
	}

	if err := container.Provide(auth.NewHandler); err != nil {
		return nil, err
	}

	// rbac module
	if err := container.Provide(rbac.NewRepository); err != nil {
		return nil, err
	}

	if err := container.Provide(rbac.NewService); err != nil {
		return nil, err
	}

	if err := container.Provide(rbac.NewHandler); err != nil {
		return nil, err
	}

	// other domain

	// end

	if err := container.Provide(router.SetupRoutes); err != nil {
		return nil, err
	}

	if err := container.Provide(ProvideHttpServer); err != nil {
		return nil, err
	}

	return container, nil
}

func ProvideHttpServer(mux *chi.Mux) (server.Server, error) {
	svr := server.New()
	svr.WithRoute(mux)
	return svr, nil
}
