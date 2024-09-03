package schema

import (
	"bytes"
	"fmt"
)

func ExampleTable_CreateStatements() {

	d := Table{
		Name: "MyTable",
		Columns: []*Column{
			{Name: "uid", Type: UUID, Comment: "key"},
			{Name: "name", Type: Text},
			{Name: "created_at", Type: Timestamp},
			{Name: "deleted_at", Type: Timestamp, Nullable: true},
			{Name: "untyped", Nullable: true},
			{Name: "age", Type: Int, Nullable: true},
			{Name: "gen", Generated: &Generated{Expression: `name + " " + age`, Storage: Virtual}},
		},
		PK: []string{"uid"},
		ForeignKeys: []*ForeignKey{
			{
				ChildKey:    []string{"uid"},
				ParentTable: "other_table", ParentKey: []string{"uid"},
				OnDelete: Cascade,
				OnUpdate: Cascade,
			},
		},
		WithoutRowID: true,
		Indices: []*Index{
			{Name: "idx_name", Columns: []string{"name"}, Unique: true},
		},
	}

	b := bytes.Buffer{}
	d.CreateStatements(&b)
	fmt.Print(b.String())
	// Output:
	// create table MyTable (
	//     uid         uuid       not null,
	//     name        text       not null,
	//     created_at  timestamp  not null,
	//     deleted_at  timestamp,
	//     untyped,
	//     age         int,
	//     gen                    not null generated always as (name + " " + age) virtual,
	//     primary key (uid),
	//     foreign key (uid) references other_table(uid) on delete cascade on update cascade
	// ) without rowid;
	// create unique index idx_name on MyTable(name);
}
