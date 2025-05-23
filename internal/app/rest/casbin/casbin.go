package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

type Enforcer interface {
	Enforce(...any) (bool, error)
}

func New(conf, policy string) (Enforcer, error) {
	m, err := model.NewModelFromFile(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	a := fileadapter.NewAdapter(policy)
	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	return enforcer, nil
}
