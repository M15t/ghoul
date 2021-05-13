package dbutil

import (
	"reflect"
	"strings"

	"github.com/imdatngo/gowhere"
	"gorm.io/gorm"
)

// NewDB creates new DB instance
func NewDB(model interface{}) *DB {
	return &DB{model, nil}
}

// DB represents the client for common usages
type DB struct {
	// Model must be set to a specific model instance. e.g: model.User{}
	Model interface{}
	// GDB holds previous DB instance that just executed the query
	GDB *gorm.DB
}

// Intf represents the common db interface
type Intf interface {
	// Create creates a new record on database.
	// `input` must be a non-nil pointer of the model. e.g: `input := &model.User{}`
	Create(db *gorm.DB, input interface{}) error
	// View returns single record matching the given conditions.
	// `output` must be a non-nil pointer of the model. e.g: `output := new(model.User)`
	// Note: RecordNotFound error is returned when there is no record that matches the conditions
	View(db *gorm.DB, output interface{}, cond ...interface{}) error
	// List returns list of records retrievable after filter & pagination if given.
	// `output` must be a non-nil pointer of slice of the model. e.g: `data := []*model.User{}; db.List(dbconn, &data, nil, nil)`
	// `lq` can be nil, then no filter & pagination are applied
	// `count` can also be nil, then no extra query is executed to get the total count
	List(db *gorm.DB, output interface{}, lq *ListQueryCondition, count *int) error
	// Update updates data of the records matching the given conditions.
	// `updates` could be a model struct or map[string]interface{}
	// Note: DB.Model must be provided in order to get the correct model/table
	Update(db *gorm.DB, updates interface{}, cond ...interface{}) error
	// Delete deletes record matching given conditions.
	// `cond` can be an instance of the model, then primary key will be used as the condition
	Delete(db *gorm.DB, cond ...interface{}) error
	// Exist checks whether there is record matching the given conditions.
	Exist(db *gorm.DB, cond ...interface{}) (bool, error)
	// CreateInBatches creates batch of new record on database.
	// `input` must be a array non-nil pointer of the model. e.g: `input := []*model.User`
	CreateInBatches(db *gorm.DB, input interface{}, batchSize int) error
}

// ListQueryCondition holds data used for db queries
type ListQueryCondition struct {
	Filter  *gowhere.Plan
	Sort    []string
	Page    int
	PerPage int
}

// Create creates a new record on database.
func (cdb *DB) Create(db *gorm.DB, input interface{}) error {
	cdb.GDB = db.Create(input)
	return cdb.GDB.Error
}

// View returns single record matching the given conditions.
func (cdb *DB) View(db *gorm.DB, output interface{}, cond ...interface{}) error {
	where := ParseCond(cond...)
	cdb.GDB = db.First(output, where...)
	return cdb.GDB.Error
}

// List returns list of records retrievable after filter & pagination if given.
func (cdb *DB) List(db *gorm.DB, output interface{}, lq *ListQueryCondition, count *int64) error {
	if lq != nil {
		if lq.Filter != nil {
			db = db.Where(lq.Filter.SQL(), lq.Filter.Vars()...)
		}

		if lq.PerPage > 0 {
			db = db.Limit(lq.PerPage)
			if lq.Page > 1 {
				db = db.Offset(lq.Page*lq.PerPage - lq.PerPage)
			}
		}

		if lq.Sort != nil && len(lq.Sort) > 0 {
			// Note: It's up to who using this package to validate the sort fields!
			db = db.Order(strings.Join(lq.Sort, ", "))
		}
	}

	cdb.GDB = db.Find(output)
	if err := cdb.GDB.Error; err != nil {
		return err
	}

	// Only count total records if requested
	if count != nil {
		if err := cdb.GDB.Limit(-1).Offset(-1).Count(count).Error; err != nil {
			return err
		}
	}

	return nil
}

// Update updates data of the records matching the given conditions.
func (cdb *DB) Update(db *gorm.DB, updates interface{}, cond ...interface{}) error {
	db = db.Model(cdb.Model)
	if len(cond) > 0 {
		where := ParseCond(cond...)
		db = db.Where(where[0], where[1:]...)
	}
	cdb.GDB = db.Omit("id").Updates(updates)
	return cdb.GDB.Error
}

// Delete deletes record matching given conditions.
func (cdb *DB) Delete(db *gorm.DB, cond ...interface{}) error {
	if len(cond) == 1 {
		val := reflect.ValueOf(cond[0])
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct {
			return db.Delete(cond[0]).Error
		}
	}
	where := ParseCond(cond...)
	cdb.GDB = db.Delete(cdb.Model, where...)
	return cdb.GDB.Error
}

// Exist checks whether there is record matching the given conditions.
func (cdb *DB) Exist(db *gorm.DB, cond ...interface{}) (bool, error) {
	var count int64
	count = 0
	where := ParseCond(cond...)
	cdb.GDB = db.Model(cdb.Model).Where(where[0], where[1:]...).Count(&count)
	return count > 0, cdb.GDB.Error
}

// CreateInBatches creates batch of new record on database.
func (cdb *DB) CreateInBatches(db *gorm.DB, input interface{}) error {
	cdb.GDB = db.CreateInBatches(input, 1000)
	return cdb.GDB.Error
}
