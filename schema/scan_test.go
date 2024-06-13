package schema

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

func ExampleDatabase() {

	db, _ := sql.Open("sqlite3", ":memory:")
	_, err := db.Exec(`
		create table t1 (
			id       int64,
			uid      uuid  not null unique,
			n1       int   default 42,
			n2       int   default null,
			n3       int,
			nn1      int   not null default 42,
			nn2      int   not null default null,
			nn3      int   not null,
			s1       text  default "42",
			s2       text  default "null",
			s3       text  default null,
			s4       text,
			s5       text  not null default "42",
			s6       text  not null default "null",
			s7       text  not null default null,
			s8       text  not null,
			f1		 timestamp not null default current_timestamp,
			primary key (id)
		);
	`)
	if err != nil {
		fmt.Print(err)
	}
	_, err = db.Exec(`
		create table t2 (
			id      int64 not null unique,
			t_id    int64  not null,
			primary key (id)
			foreign key (t_id) references t(id) on delete cascade on update cascade
		);
	`)
	if err != nil {
		fmt.Print(err)
	}

	dbsch, err := Scan(db)
	if err != nil {
		fmt.Print(err)
	} else {
		y, _ := yaml.Marshal(dbsch.Tables)
		fmt.Print(string(y))
	}

	// Output:
	// - table: t2
	//   columns:
	//     - {name: id, type: int64}
	//     - {name: t_id, type: int64}
	//   indices:
	//     - {name: sqlite_autoindex_t2_1, unique: true}
	//   pk: [id]
	//   foreign_keys:
	//     - child_key: [t_id]
	//       parent_table: id
	//       parent_key: [id]
	//       on_delete: cascade
	//       on_update: cascade
	// - table: t1
	//   columns:
	//     - {name: id, type: int64, nullable: true}
	//     - {name: uid, type: uuid}
	//     - {name: n1, type: INT, nullable: true, default: 42}
	//     - {name: n2, type: INT, nullable: true, default: null}
	//     - {name: n3, type: INT, nullable: true}
	//     - {name: nn1, type: INT, default: 42}
	//     - {name: nn2, type: INT, default: null}
	//     - {name: nn3, type: INT}
	//     - {name: s1, type: TEXT, nullable: true, default: '"42"'}
	//     - {name: s2, type: TEXT, nullable: true, default: '"null"'}
	//     - {name: s3, type: TEXT, nullable: true, default: null}
	//     - {name: s4, type: TEXT, nullable: true}
	//     - {name: s5, type: TEXT, default: '"42"'}
	//     - {name: s6, type: TEXT, default: '"null"'}
	//     - {name: s7, type: TEXT, default: null}
	//     - {name: s8, type: TEXT}
	//     - {name: f1, type: timestamp, default: CURRENT_TIMESTAMP}
	//   indices:
	//     - {name: sqlite_autoindex_t1_2, unique: true}
	//     - {name: sqlite_autoindex_t1_1, unique: true}
	//   pk: [id]

}
