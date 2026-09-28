package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
)

type recoveryDriver struct {
	offline    bool
	failSchema bool
	executions int
}

func (d *recoveryDriver) Open(string) (driver.Conn, error)    { return d, nil }
func (d *recoveryDriver) Close() error                        { return nil }
func (d *recoveryDriver) Begin() (driver.Tx, error)           { return nil, errors.New("unused") }
func (d *recoveryDriver) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unused") }
func (d *recoveryDriver) Ping(context.Context) error {
	if d.offline {
		return errors.New("offline")
	}
	return nil
}
func (d *recoveryDriver) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	d.executions++
	if d.failSchema {
		return nil, errors.New("schema unavailable")
	}
	return driver.RowsAffected(0), nil
}

func TestDatabaseRecovery(t *testing.T) {
	d := &recoveryDriver{offline: true}
	sql.Register("recovery-test", d)
	var err error
	DB, err = sql.Open("recovery-test", "")
	if err != nil {
		t.Fatal(err)
	}
	schemaReady = false
	t.Cleanup(func() { DB.Close(); DB = nil; schemaReady = false })
	ctx := context.Background()
	if EnsureReady(ctx) == nil {
		t.Fatal("outage should fail")
	}
	d.offline = false
	d.failSchema = true
	if EnsureReady(ctx) == nil {
		t.Fatal("schema failure should fail")
	}
	d.failSchema = false
	if err := EnsureReady(ctx); err != nil {
		t.Fatalf("recovery: %v", err)
	}
	if d.executions != 3 {
		t.Fatalf("expected failed schema query and two successful queries, got %d", d.executions)
	}
	if err := EnsureReady(ctx); err != nil {
		t.Fatal(err)
	}
	if d.executions != 3 {
		t.Fatal("schema initialization repeated after success")
	}
	d.offline = true
	if EnsureReady(ctx) == nil {
		t.Fatal("later outage should fail")
	}
	d.offline = false
	if err := EnsureReady(ctx); err != nil {
		t.Fatalf("second recovery: %v", err)
	}
}
