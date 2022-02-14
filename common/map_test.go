package common

import (
	"testing"
)

type MigrateState string

const (
	Migrating    MigrateState = "migrating"
	Migrated     MigrateState = "migrated"
	MigrateError MigrateState = "migrate_error"
)

func TestMap(t *testing.T) {
	m := map[string]interface{}{"state": "migrating", "type": "vm-migrate"}
	t.Log(len(m))
	t.Log(m["state"] == Migrating)
	t.Log(m["type"] == "vm-migrate")
}
