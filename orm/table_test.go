package orm

import (
	"testing"
)

func TestBindReceivers(t *testing.T) {
	tests := []struct {
		name        string
		table       *Table
		dst         interface{}
		wantErr     bool
		errType     error
		numBindings int
	}{
		{
			name: "basic struct binding",
			table: &Table{
				Name:    "test",
				columns: nameset{"id": struct{}{}, "name": struct{}{}},
			},
			dst: &struct {
				ID   int    `orm:"id"`
				Name string `orm:"name"`
			}{},
			wantErr:     false,
			numBindings: 2,
		},
		{
			name: "missing required column",
			table: &Table{
				Name:    "test",
				columns: nameset{"id": struct{}{}},
			},
			dst: &struct {
				ID   int    `orm:"id"`
				Name string `orm:"name"`
			}{},
			wantErr: true,
			errType: ErrMissingColumns([]string{"name"}),
		},
		{
			name: "optional column not present",
			table: &Table{
				Name:    "test",
				columns: nameset{"id": struct{}{}},
			},
			dst: &struct {
				ID   int    `orm:"id"`
				Name string `orm:"?name"`
			}{},
			wantErr:     false,
			numBindings: 1,
		},
		{
			name: "conditional field",
			table: &Table{
				Name:    "test",
				columns: nameset{"id": struct{}{}, "status": struct{}{}},
			},
			dst: &struct {
				ID     int    `orm:"id"`
				Status string `orm:"$status"`
			}{},
			wantErr:     false,
			numBindings: 2,
		},
		{
			name: "no orm tags",
			table: &Table{
				Name:    "test",
				columns: nameset{"id": struct{}{}, "name": struct{}{}},
			},
			dst: &struct {
				ID   int
				Name string
			}{},
			wantErr: true,
			errType: ErrNoBindingsProduced,
		},
		{
			name: "nested struct with !",
			table: &Table{
				Name:    "test",
				columns: nameset{"id": struct{}{}, "inner_field": struct{}{}},
			},
			dst: &struct {
				ID    int `orm:"id"`
				Inner struct {
					Field string `orm:"inner_field"`
				} `orm:"!"`
			}{},
			wantErr:     false,
			numBindings: 2,
		},
		{
			name: "nested struct with ?",
			table: &Table{
				Name:    "test",
				columns: nameset{"id": struct{}{}},
			},
			dst: &struct {
				ID    int `orm:"id"`
				Inner struct {
					Field string `orm:"missing_field"`
				} `orm:"?"`
			}{},
			wantErr:     false,
			numBindings: 1,
		},
		{
			name: "complex nested structs with multiple levels",
			table: &Table{
				Name: "test",
				columns: nameset{
					"id":               struct{}{},
					"name":             struct{}{},
					"address_street":   struct{}{},
					"address_city":     struct{}{},
					"contact_email":    struct{}{},
					"contact_phone":    struct{}{},
					"metadata_created": struct{}{},
					"metadata_updated": struct{}{},
				},
			},
			dst: &struct {
				ID      int    `orm:"id"`
				Name    string `orm:"name"`
				Address struct {
					Street string `orm:"address_street"`
					City   string `orm:"address_city"`
				} `orm:"!"`
				Contact struct {
					Email string `orm:"contact_email"`
					Phone string `orm:"contact_phone"`
				} `orm:"!"`
				Metadata struct {
					Created string `orm:"metadata_created"`
					Updated string `orm:"metadata_updated"`
				} `orm:"!"`
			}{},
			wantErr:     false,
			numBindings: 8,
		},
		{
			name: "mixed optional and required nested structs",
			table: &Table{
				Name: "test",
				columns: nameset{
					"id":            struct{}{},
					"name":          struct{}{},
					"primary_email": struct{}{},
					"backup_email":  struct{}{},
				},
			},
			dst: &struct {
				ID             int    `orm:"id"`
				Name           string `orm:"name"`
				PrimaryContact struct {
					Email string `orm:"primary_email"`
				} `orm:"!"`
				BackupContact struct {
					Email string `orm:"backup_email"`
				} `orm:"?"`
				MissingContact struct {
					Phone string `orm:"missing_phone"`
				} `orm:"?"`
			}{},
			wantErr:     false,
			numBindings: 4,
		},
		{
			name: "deeply nested structs with conditional fields",
			table: &Table{
				Name: "test",
				columns: nameset{
					"id":             struct{}{},
					"status":         struct{}{},
					"user_name":      struct{}{},
					"user_email":     struct{}{},
					"settings_theme": struct{}{},
				},
			},
			dst: &struct {
				ID     int    `orm:"id"`
				Status string `orm:"$status"`
				User   struct {
					Name     string `orm:"user_name"`
					Email    string `orm:"$user_email"`
					Settings struct {
						Theme  string `orm:"$settings_theme"`
						Layout string `orm:"?settings_layout"`
					} `orm:"!"`
				} `orm:"!"`
			}{},
			wantErr:     false,
			numBindings: 5,
		},
		{
			name: "nested structs with missing required fields",
			table: &Table{
				Name: "test",
				columns: nameset{
					"id":   struct{}{},
					"name": struct{}{},
				},
			},
			dst: &struct {
				ID       int    `orm:"id"`
				Name     string `orm:"name"`
				Required struct {
					Field string `orm:"required_field"`
				} `orm:"!"`
			}{},
			wantErr: true,
			errType: ErrMissingColumns([]string{"required_field"}),
		},
		{
			name: "embedded struct",
			table: &Table{
				Name: "test",
				columns: nameset{
					"id":     struct{}{},
					"name":   struct{}{},
					"street": struct{}{},
					"city":   struct{}{},
				},
			},
			dst: &struct {
				ID      int    `orm:"id"`
				Name    string `orm:"name"`
				Address struct {
					Street string `orm:"street"`
					City   string `orm:"city"`
				} `orm:"!"`
			}{},
			wantErr:     false,
			numBindings: 4,
		},
		{
			name: "multiple embedded structs",
			table: &Table{
				Name: "test",
				columns: nameset{
					"id":          struct{}{},
					"home_street": struct{}{},
					"home_city":   struct{}{},
					"work_street": struct{}{},
					"work_city":   struct{}{},
				},
			},
			dst: &struct {
				ID          int `orm:"id"`
				HomeAddress struct {
					Street string `orm:"home_street"`
					City   string `orm:"home_city"`
				} `orm:"!"`
				WorkAddress struct {
					Street string `orm:"work_street"`
					City   string `orm:"work_city"`
				} `orm:"!"`
			}{},
			wantErr:     false,
			numBindings: 5,
		},
		{
			name: "embedded struct with optional fields",
			table: &Table{
				Name: "test",
				columns: nameset{
					"id":            struct{}{},
					"contact_email": struct{}{},
				},
			},
			dst: &struct {
				ID      int `orm:"id"`
				Contact struct {
					Email string `orm:"contact_email"`
					Phone string `orm:"?contact_phone"`
				} `orm:"!"`
			}{},
			wantErr:     false,
			numBindings: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bindings, err := tt.table.bind_receivers(tt.dst)

			if tt.wantErr {
				if err == nil {
					t.Errorf("bind_receivers() expected error but got none")
					return
				}
				if tt.errType != nil && err.Error() != tt.errType.Error() {
					t.Errorf("bind_receivers() error = %v, want %v", err, tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("bind_receivers() unexpected error: %v", err)
				return
			}

			if len(bindings.receivers) != tt.numBindings {
				t.Errorf("bind_receivers() got %d bindings, want %d", len(bindings.receivers), tt.numBindings)
			}
			if len(bindings.selectors) != tt.numBindings {
				t.Errorf("bind_receivers() got %d selectors, want %d", len(bindings.selectors), tt.numBindings)
			}
		})
	}
}

func TestBindReceiversInvalidInput(t *testing.T) {
	table := &Table{
		Name:    "test",
		columns: nameset{"id": struct{}{}},
	}

	t.Run("non-pointer input", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("bind_receivers() expected panic for non-pointer input")
			}
		}()

		var s struct{}
		table.bind_receivers(s)
	})

	t.Run("pointer to non-struct", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("bind_receivers() expected panic for pointer to non-struct")
			}
		}()

		var i int
		table.bind_receivers(&i)
	})
}
