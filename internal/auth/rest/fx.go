package rest

import (
	"context"
	"fmt"
	"prodigo/pkg/config"

	"go.uber.org/fx"
)

var Module = fx.Module("rest",
	fx.Provide(New),
	fx.Invoke(func(lc fx.Lifecycle, s *Server, conf *config.Config) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := s.Start(conf.Host, conf.Port); err != nil {
					return fmt.Errorf("failed to start server: %w", err)
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return s.Stop(ctx)
			},
		})
	}),
)
