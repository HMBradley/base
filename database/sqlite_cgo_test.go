//go:build cgo
// +build cgo

// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package database

import (
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moov-io/base/log"
)

func TestSQLite__basic(t *testing.T) {
	db := CreateTestSQLiteDB(t)
	defer db.Close()

	err := db.DB.Ping()
	require.NoError(t, err)

	r, err := db.Query("select * from tests")
	require.NoError(t, err)
	require.NoError(t, r.Err())
	defer r.Close()

	if runtime.GOOS == "windows" {
		t.Skip("/dev/null doesn't exist on Windows")
	}

	// error case
	s, err := sqliteConnection(log.NewNopLogger(), "/tmp/path/doesnt/exist")
	require.NoError(t, err)

	ctx, cancelFunc := context.WithCancel(context.Background())

	conn, _ := s.Connect(ctx)
	err = conn.Ping()
	require.EqualError(t, err, "unable to open database file: no such file or directory")

	cancelFunc()

	conn.Close()
}

func TestSQLiteUniqueViolation(t *testing.T) {
	err := errors.New(`problem upserting depository="7d676c65eccd48090ff238a0d5e35eb6126c23f2", userId="80cfe1311d9eb7659d02cba9ee6cb04ed3739a85": UNIQUE constraint failed: depositories.depository_id`)
	if !UniqueViolation(err) {
		t.Error("should have matched unique violation")
	}
}
