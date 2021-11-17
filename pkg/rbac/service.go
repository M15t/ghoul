package rbac

import (
	"net/http"

	"ghoul/pkg/rbac/casbinadapter"
	"ghoul/pkg/server"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"gorm.io/gorm"
)

// Custom errors
var (
	ErrForbiddenAccess = server.NewHTTPError(http.StatusForbidden, "FORBIDDEN", "You don't have permission to access the requested resource")
	ErrForbiddenAction = server.NewHTTPError(http.StatusForbidden, "FORBIDDEN", "You don't have permission to perform this action")
)

// Config represents the config for RBAC service
type Config struct {
	Model     model.Model
	Adapter   persist.Adapter
	GormDB    *gorm.DB
	EnableLog bool
}

// RBAC is RBAC application service
type RBAC struct {
	*casbin.Enforcer
}

// Intf represents common interface for the RBAC service
type Intf interface {
	Enforce(rvals ...interface{}) bool
}

// DefaultConfig represents the default configuration
var DefaultConfig = Config{
	Model:     NewRBACModel(),
	Adapter:   nil,
	GormDB:    nil,
	EnableLog: true,
}

// New creates new RBAC service with default configuration
func New() *RBAC {
	return NewWithConfig(DefaultConfig)
}

// NewWithConfig creates new RBAC service with custom configuration
func NewWithConfig(cfg Config) *RBAC {
	if cfg.Model == nil {
		cfg.Model = DefaultConfig.Model
	}
	if cfg.GormDB == nil {
		cfg.GormDB = DefaultConfig.GormDB
	}
	if cfg.GormDB != nil {
		cfg.Adapter = casbinadapter.NewAdapter(cfg.GormDB)
	} else if cfg.Adapter == nil {
		cfg.Adapter = DefaultConfig.Adapter
	}

	var ce *casbin.Enforcer
	if cfg.Adapter != nil {
		ce = casbin.NewEnforcer(cfg.Model, cfg.Adapter, cfg.EnableLog)
	} else {
		ce = casbin.NewEnforcer(cfg.Model, cfg.EnableLog)
	}

	return &RBAC{ce}
}

// NewRBACModel initializes the RBAC casbin model
func NewRBACModel() model.Model {
	m := casbin.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `g(r.sub, p.sub) && (r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*")`)
	return m
}

// NewRBACWithLevelInheritanceModel initializes the RBAC with level inheritance model
func NewRBACWithLevelInheritanceModel() model.Model {
	m := casbin.NewModel()
	m.AddDef("r", "r", "sub, lvl, obj, act")
	m.AddDef("p", "p", "sub, lvl, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("g", "g2", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `((r.sub != p.sub && g(r.sub, p.sub)) || (r.sub == p.sub && g2(p.lvl, r.lvl))) && (r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*")`)
	return m
}

// NewRBACWithDomainModel initializes the RBAC with domain model
func NewRBACWithDomainModel() model.Model {
	m := casbin.NewModel()
	m.AddDef("r", "r", "sub, dom, obj, act")
	m.AddDef("p", "p", "sub, dom, obj, act")
	m.AddDef("g", "g", "_, _, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act`)
	return m
}
