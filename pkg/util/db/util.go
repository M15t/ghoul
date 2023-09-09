package dbutil

import (
	"fmt"

	"github.com/imdatngo/gowhere"
	"gorm.io/gorm"
)

// parseCondWithConfig returns standard [sqlString, vars] format for query, powered by gowhere package (configurable version)
func parseCondWithConfig(cfg gowhere.Config, cond ...interface{}) []interface{} {
	if len(cond) == 1 {
		switch c := cond[0].(type) {
		case map[string]interface{}, []interface{}:
			cond[0] = gowhere.WithConfig(cfg).Where(c)
		}

		if plan, ok := cond[0].(*gowhere.Plan); ok {
			return append([]interface{}{plan.SQL()}, plan.Vars()...)
		}
	}
	return cond
}

// parseCond returns standard [sqlString, vars] format for query, powered by gowhere package (with default config)
func parseCond(cond ...interface{}) []interface{} {
	return parseCondWithConfig(gowhere.DefaultConfig, cond...)
}

// InTransaction defines the transaction wrapper function
type InTransaction func(tx *gorm.DB) error

// Transaction execute the input func in a transaction
func Transaction(db *gorm.DB, fn InTransaction) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = fmt.Errorf(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic: %+v", x)
			}
		}
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()
	return fn(tx)
}
