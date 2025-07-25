package rest

import (
	"context"
	"fmt"
	"log"
	"prodigo/pkg/config"

	"go.uber.org/fx"
)

var Module = fx.Module("rest",
	fx.Provide(New),
	fx.Invoke(func(lc fx.Lifecycle, s *Server, conf *config.Config) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := s.Start(conf.AppHost, conf.AppPort); err != nil {
						log.Fatalf("failed to start server: %v", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				if err := s.Stop(ctx); err != nil {
					return fmt.Errorf("failed to stop server: %w", err)
				}
				return nil
			},
		})
	}),
)
