# go-db3

[![GoDoc](https://godoc.org/github.com/adnsv/go-db3?status.svg)](https://godoc.org/github.com/adnsv/go-db3)


Libraries for managing SQLite databases in GO

Features:

- retrieving schema in existing databases
- validating scanned schema for table presense, fields, and (to some extent) field types
- convenience routines for creating new schema

Warning: unstable API, WIP

## Reference

Automatically generated documentation for the package can be viewed online here:
http://pkg.go.dev/github.com/adnsv/go-db3

## ORM Tags and Struct Bindings

The package provides ORM capabilities through struct tags that define how struct fields map to database columns. Here's a comprehensive guide to the tag syntax:

### Basic Tag Syntax

- `orm:"column_name"` - Required binding to a column
- `orm:"?column_name"` - Optional binding (won't error if column missing)
- `orm:"$column_name"` - Conditional binding (converts to empty string if NULL)
- `orm:"$?column_name"` - Optional conditional binding (combines both behaviors)
- `orm:"name1|name2"` - Alternative column names (binds to first existing column)

### Nested Struct Tags

For nested structs, use these modifiers:
- `orm:"!"` - Required nested struct
- `orm:"?"` - Optional nested struct

### Examples

Basic struct with required and optional fields:
```go
type User struct {
    ID       int    `orm:"id"`              // required field
    Name     string `orm:"?display_name"`    // optional field
    Status   string `orm:"$status"`         // conditional field
    Bio      string `orm:"$?description"`   // optional conditional field
    Email    string `orm:"email|contact"`   // alternative names
}
```

Embedded struct example:
```go
type Customer struct {
    ID      int     `orm:"id"`
    Address struct {
        Street string `orm:"address_street"`
        City   string `orm:"address_city"`
    } `orm:"!"` // required nested struct
    
    OptionalInfo struct {
        Phone string `orm:"optional_phone"`
        Notes string `orm:"optional_notes"`
    } `orm:"?"` // optional nested struct
}
```

Complex nested structure:
```go
type Customer struct {
    ID      int    `orm:"id"`
    Name    string `orm:"name"`
    Address struct {
        Street string `orm:"address_street"`
        City   string `orm:"address_city"`
        Contact struct {
            Phone string `orm:"address_contact_phone"`
            Email string `orm:"address_contact_email"`
        } `orm:"!"`
    } `orm:"!"` // required nested struct
    
    OptionalInfo struct {
        Notes string `orm:"optional_notes"`
        Tags  string `orm:"?optional_tags"`
    } `orm:"?"` // optional nested struct
}
```

### Tag Behavior

1. Required fields (`orm:"column_name"`):
   - Must exist in the database table
   - Will return error if column is missing
   - Multiple column names can be specified with pipe separator

2. Optional fields (`orm:"?column_name"`):
   - Can be missing from the database
   - Skipped if column doesn't exist
   - Inherits optional status from parent struct if parent is marked optional
   - Can be combined with conditional modifier (`$?column_name`)

3. Conditional fields (`orm:"$column_name"`):
   - Converts NULL values to empty string/zero value
   - Generates SQL CASE expression: `(case when column notnull then column else "" end)`
   - Must specify column name after the `$` prefix
   - Can be combined with optional modifier (`$?column_name`) to handle missing columns

4. Alternative names (`orm:"name1|name2"`):
   - Tries each name in order
   - Binds to first matching column
   - Can be combined with optional and conditional modifiers (e.g., `orm:"?name1|name2"`, `orm:"$name1|name2"`)

5. Nested structs:
   - `!` - All fields must be present (default behavior)
   - `?` - Fields can be missing
   - Can be nested multiple levels deep
   - All fields in a nested struct inherit the optional status if parent is marked optional
   - Each field must specify its complete column name - there is no automatic prefixing
   - The struct tag (`!` or `?`) only controls whether missing fields are allowed
   - Empty tag (`orm:""`) means the field is ignored in binding

### Error Handling

The package will return specific errors in these cases:
- `ErrTableDoesNotExist` - When the specified table is not found
- `ErrEmptyTableSchema` - When table exists but has no columns
- `ErrNoBindingsProduced` - When no valid bindings were created (e.g., no orm tags)
- `ErrMissingColumns` - When required columns are not found in the table
