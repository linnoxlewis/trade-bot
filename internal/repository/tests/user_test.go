package tests

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/repository"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	user := &domain.User{
		ID:       1,
		Username: "testUser",
	}

	mock.ExpectExec("INSERT INTO users").WithArgs(user.ID, user.Username, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Errorf("error was not expected while creating user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestExistUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	testCases := []struct {
		id    int64
		exist bool
	}{
		{1, true},
		{2, false},
	}

	for _, tc := range testCases {
		mock.ExpectQuery("SELECT EXISTS").WithArgs(tc.id).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(tc.exist))

		exist := repo.ExistUser(context.Background(), tc.id)
		if exist != tc.exist {
			t.Errorf("expected existence to be %v, got %v for user id %d", tc.exist, exist, tc.id)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestIsAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	testCases := []struct {
		id      int64
		isAdmin bool
		err     error
	}{
		{1, true, nil},
		{2, false, nil},
	}

	for _, tc := range testCases {
		mock.ExpectQuery("SELECT EXISTS").WithArgs(tc.id, true).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(tc.isAdmin))

		isAdmin, err := repo.IsAdmin(context.Background(), tc.id)
		if err != tc.err {
			t.Errorf("expected error %v, got %v", tc.err, err)
		}
		if isAdmin != tc.isAdmin {
			t.Errorf("expected isAdmin to be %v, got %v", tc.isAdmin, isAdmin)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAdminIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2)
	mock.ExpectQuery("SELECT id FROM users WHERE is_admin").WithArgs(true).WillReturnRows(rows)

	ids, err := repo.GetAdminIds(context.Background())
	if err != nil {
		t.Errorf("error was not expected while getting admin ids: %s", err)
	}

	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Errorf("expected ids [1 2], got %v", ids)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
