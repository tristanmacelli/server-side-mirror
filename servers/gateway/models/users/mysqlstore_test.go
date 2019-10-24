package users

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// a failing test case
func TestGetByID(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	columns := []string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}

	mock.ExpectPrepare("SELECT \\* FROM users")
	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("1").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByID(1); err != nil {
		t.Errorf("was not expecting an error, but there was none")
		fmt.Print(user)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByEmail(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	columns := []string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}

	mock.ExpectPrepare("SELECT \\* FROM users")
	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("user@domain.com").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByEmail("user@domain.com"); err != nil {
		t.Errorf("was not expecting an error, but there was none")
		fmt.Print(user)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByUsername(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	columns := []string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}

	mock.ExpectPrepare("SELECT \\* FROM users")
	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("Sam").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByUserName("Sam"); err != nil {
		t.Errorf("was not expecting an error, but there was none")
		fmt.Print(user)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsert(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO users(email, passHash, username, firstname, lastname, photoURL) VALUES (?,?,?,?,?,?)")
	mock.ExpectExec("INSERT INTO users").
		WithArgs()

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByUserName("Sam"); err != nil {
		t.Errorf("was not expecting an error, but there was none")
		fmt.Print(user)
	}
}
