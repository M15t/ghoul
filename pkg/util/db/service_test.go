package dbutil

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	type args struct {
		dialect string
		dbPsn   string
		cfg     *gorm.Config
	}
	tests := []struct {
		name    string
		args    args
		wantDb  *gorm.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDb, err := New(tt.args.dialect, tt.args.dbPsn, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDb, tt.wantDb) {
				t.Errorf("New() = %v, want %v", gotDb, tt.wantDb)
			}
		})
	}
}
