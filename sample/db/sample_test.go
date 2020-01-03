package db

import (
	"go-db-testing-sample/dbtest"
	"go-db-testing-sample/sample/db/models"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestDBHandler_FetchSample(t *testing.T) {

	db, tearDown := dbtest.SetUpTestContainerPostgres(
		[]string{
			"sample",
		})
	defer tearDown()

	type fields struct {
		db *gorm.DB
	}
	type fixture struct {
		File  string
		Model interface{}
	}
	tests := []struct {
		name     string
		fields   fields
		fixtures []fixture
		want     *models.Sample
		wantErr  bool
	}{
		{
			name:   "sample",
			fields: fields{db: db},
			fixtures: []fixture{
				{
					File:  "test_sample",
					Model: &[]*models.Sample{},
				},
			},
			want: &models.Sample{
				ID:    1234,
				Name:  "kohey",
				Email: "hogehoge@sample.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := tt.fields.db.Begin()
			s := &DBHandler{
				db: tx,
			}

			for _, fixture := range tt.fixtures {
				if err := dbtest.SetTestData(s.db, fixture.File, fixture.Model); err != nil {
					t.Errorf("Set test data error: %s", err)
				}
			}

			got, err := s.FetchSample()
			if (err != nil) != tt.wantErr {
				t.Errorf("DBHandler.FetchSample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBHandler.FetchSample() = %v, want %v", got, tt.want)
			}
		})
	}
}
