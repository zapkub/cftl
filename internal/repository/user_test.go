package repository

import (
	"context"
	"testing"

	"github.com/zapkub/cftl/internal"
)

func TestInsertSession(t *testing.T) {
	ctx := context.Background()
	defer ResetTestDB(testDB, t)
	_, err := insertUser(ctx, testDB.db, &internal.User{
		Email:    "rungsikorn@me.com",
		Username: "rungsikorn",
		Name:     "zdcdos",
	})
	if err != nil {
		t.Fatalf("insertUser error: %v", err)
	}

	token, err := insertSession(ctx, testDB.db, &internal.Session{
		AccessToken:  "test_token",
		Email:        "rungsikorn@me.com",
		RefreshToken: "1234123",
		Origin:       internal.SessionOriginGithub,
	})
	if err != nil {
		t.Fatalf("insertSession error: %v", err)
	}

	if expected, actual := "test_token", token; expected != actual {
		t.Fatalf("insertSession return %s but expected %s", actual, expected)
	}

}

func TestGetSession(t *testing.T) {
	ctx := context.Background()
	defer ResetTestDB(testDB, t)
	_, err := insertUser(ctx, testDB.db, &internal.User{
		Email:    "rungsikorn@me.com",
		Username: "rungsikorn",
		Name:     "zdcdos",
	})
	if err != nil {
		t.Fatalf("insertUser error: %v", err)
	}

	token, err := insertSession(ctx, testDB.db, &internal.Session{
		AccessToken:  "test_token",
		Email:        "rungsikorn@me.com",
		RefreshToken: "1234123",
		Origin:       internal.SessionOriginGithub,
	})
	if err != nil {
		t.Fatalf("insertSession error: %v", err)
	}

	ss, err := getSession(ctx, testDB.db, token)
	if err != nil {
		t.Fatalf("getSession error: %v", err)
	}

	if actual, expected := ss.Email, "rungsikorn@me.com"; actual != expected {
		t.Fatalf("getSession expecte %s but retrieve %s", expected, actual)
	}
}

func TestInsertUser(t *testing.T) {

	ctx := context.Background()
	defer ResetTestDB(testDB, t)

	ID, err := insertUser(ctx, testDB.db, &internal.User{
		Email:    "rungsikorn@me.com",
		Username: "rungsikorn",
		Name:     "zdcdos",
	})
	if err != nil {
		t.Fatalf("insert user error: %+v", err)
	}

	if result, expected := ID, "rungsikorn@me.com"; result != expected {
		t.Logf("InsertUser() expect %s but receive %s", expected, result)
	}

	userinfo, err := getUser(ctx, testDB.db, "rungsikorn@me.com")
	if err != nil {
		t.Fatalf("get user error: %+v", err)
	}
	if actual, expected := userinfo.Name, "zdcdos"; actual != expected {
		t.Fatalf("InsertUser() expect %v but retrieve %v", expected, actual)
	}

}
