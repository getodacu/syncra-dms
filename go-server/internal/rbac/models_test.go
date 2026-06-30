package rbac

import "testing"

func TestNormalizeCode(t *testing.T) {
	code, err := NormalizeCode(" System Administrator ")
	if err != nil {
		t.Fatalf("NormalizeCode() error = %v", err)
	}
	if code != "system_administrator" {
		t.Fatalf("NormalizeCode() = %q, want system_administrator", code)
	}
}

func TestScopeTypeValidation(t *testing.T) {
	valid := []ScopeType{ScopeGlobal, ScopeOrganizationUnit, ScopeOrganizationUnitAndChildren}
	for _, scope := range valid {
		if !scope.Valid() {
			t.Fatalf("scope %q should be valid", scope)
		}
	}
	if ScopeType("bad").Valid() {
		t.Fatal("bad scope should be invalid")
	}
}

func TestUserStatusValidation(t *testing.T) {
	if !UserStatusActive.Valid() || !UserStatusSuspended.Valid() {
		t.Fatal("expected active and suspended to be valid")
	}
	if UserStatus("bad").Valid() {
		t.Fatal("bad status should be invalid")
	}
}
