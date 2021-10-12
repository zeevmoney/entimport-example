// Code generated by entc, DO NOT EDIT.

package user

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldAge holds the string denoting the age field in the database.
	FieldAge = "age"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldLastName holds the string denoting the last_name field in the database.
	FieldLastName = "last_name"
	// FieldPhone holds the string denoting the phone field in the database.
	FieldPhone = "phone"
	// EdgeCars holds the string denoting the cars edge name in mutations.
	EdgeCars = "cars"
	// Table holds the table name of the user in the database.
	Table = "users"
	// CarsTable is the table that holds the cars relation/edge.
	CarsTable = "cars"
	// CarsInverseTable is the table name for the Car entity.
	// It exists in this package in order to avoid circular dependency with the "car" package.
	CarsInverseTable = "cars"
	// CarsColumn is the table column denoting the cars relation/edge.
	CarsColumn = "user_id"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldAge,
	FieldName,
	FieldLastName,
	FieldPhone,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}
