package orgunits

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUnitBeforeCreateSetsIDWhenEmpty(t *testing.T) {
	unit := Unit{}
	if err := unit.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate() error = %v", err)
	}
	if _, err := uuid.Parse(unit.ID); err != nil {
		t.Fatalf("BeforeCreate() ID = %q, want UUID: %v", unit.ID, err)
	}

	existingID := uuid.NewString()
	unit = Unit{ID: existingID}
	if err := unit.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate() error = %v", err)
	}
	if unit.ID != existingID {
		t.Fatalf("BeforeCreate() ID = %q, want existing ID %q", unit.ID, existingID)
	}
}

func TestNormalizeName(t *testing.T) {
	got, err := NormalizeName(" Finance ")
	if err != nil {
		t.Fatalf("NormalizeName() error = %v", err)
	}
	if got != "Finance" {
		t.Fatalf("NormalizeName() = %q, want Finance", got)
	}

	if got, err := NormalizeName(" \t\n "); err == nil {
		t.Fatalf("NormalizeName() error = nil, name = %q", got)
	}

	if got, err := NormalizeName(strings.Repeat("ă", MaxNameCharacters+1)); err == nil {
		t.Fatalf("NormalizeName() error = nil, name = %q", got)
	}
}

func TestNormalizeCode(t *testing.T) {
	t.Run("trims and uppercases non-empty code", func(t *testing.T) {
		got, err := NormalizeCode(" fin-ap ")
		if err != nil {
			t.Fatalf("NormalizeCode() error = %v", err)
		}
		if got == nil || *got != "FIN-AP" {
			t.Fatalf("NormalizeCode() = %v, want FIN-AP", got)
		}
	})

	t.Run("returns nil for whitespace", func(t *testing.T) {
		got, err := NormalizeCode(" \t\n ")
		if err != nil {
			t.Fatalf("NormalizeCode() error = %v", err)
		}
		if got != nil {
			t.Fatalf("NormalizeCode() = %v, want nil", *got)
		}
	})

	t.Run("rejects codes over forty characters", func(t *testing.T) {
		got, err := NormalizeCode(strings.Repeat("A", MaxCodeCharacters+1))
		if err == nil {
			t.Fatalf("NormalizeCode() error = nil, code = %v", got)
		}
	})

	t.Run("accepts exact boundary ascii code", func(t *testing.T) {
		got, err := NormalizeCode(strings.Repeat("A", MaxCodeCharacters))
		if err != nil {
			t.Fatalf("NormalizeCode() error = %v", err)
		}
		if got == nil || *got != strings.Repeat("A", MaxCodeCharacters) {
			t.Fatalf("NormalizeCode() = %v, want boundary code", got)
		}
	})

	t.Run("rejects non-ascii codes", func(t *testing.T) {
		got, err := NormalizeCode("FÎN")
		if err == nil {
			t.Fatalf("NormalizeCode() error = nil, code = %v", got)
		}
	})
}

func TestNormalizeDescription(t *testing.T) {
	got := NormalizeDescription(" Finance department ")
	if got == nil || *got != "Finance department" {
		t.Fatalf("NormalizeDescription() = %v, want Finance department", got)
	}

	if got := NormalizeDescription(" \t\n "); got != nil {
		t.Fatalf("NormalizeDescription() = %v, want nil", *got)
	}
}

func TestBuildTreeSortsAndNestsUnits(t *testing.T) {
	parentID := "parent"
	otherRootID := "other-root"
	now := time.Date(2026, 6, 30, 9, 15, 0, 123, time.FixedZone("UTC+2", 2*60*60))
	later := now.Add(time.Hour)

	code := "FIN"
	description := "Finance department"
	units := []Unit{
		{ID: "child-b", ParentID: &parentID, Name: "Beta", CreatedAt: now, UpdatedAt: later},
		{ID: otherRootID, Name: "Accounting", CreatedAt: now, UpdatedAt: later},
		{ID: "child-a", ParentID: &parentID, Name: "Alpha", Code: &code, Description: &description, CreatedAt: later, UpdatedAt: later},
		{ID: parentID, Name: "Finance", CreatedAt: now, UpdatedAt: now},
	}

	tree := BuildTree(units)
	if len(tree) != 2 {
		t.Fatalf("BuildTree() returned %d roots, want 2", len(tree))
	}
	if tree[0].ID != otherRootID || tree[1].ID != parentID {
		t.Fatalf("root order = [%s %s], want [%s %s]", tree[0].ID, tree[1].ID, otherRootID, parentID)
	}

	finance := tree[1]
	if len(finance.Children) != 2 {
		t.Fatalf("Finance children = %d, want 2", len(finance.Children))
	}
	if finance.Children[0].ID != "child-a" || finance.Children[1].ID != "child-b" {
		t.Fatalf("Finance child order = [%s %s], want [child-a child-b]", finance.Children[0].ID, finance.Children[1].ID)
	}
	if finance.Children[0].Code == nil || *finance.Children[0].Code != code {
		t.Fatalf("child-a code = %v, want %s", finance.Children[0].Code, code)
	}
	if finance.Children[0].Description == nil || *finance.Children[0].Description != description {
		t.Fatalf("child-a description = %v, want %s", finance.Children[0].Description, description)
	}

	wantCreatedAt := now.UTC().Format(time.RFC3339Nano)
	if tree[0].CreatedAt != wantCreatedAt {
		t.Fatalf("CreatedAt = %q, want %q", tree[0].CreatedAt, wantCreatedAt)
	}
}

func TestBuildTreeOrdersEqualNamesByID(t *testing.T) {
	parentID := "parent"
	tree := BuildTree([]Unit{
		{ID: parentID, Name: "Root"},
		{ID: "child-b", ParentID: &parentID, Name: "Same"},
		{ID: "child-a", ParentID: &parentID, Name: "Same"},
	})

	if len(tree) != 1 || len(tree[0].Children) != 2 {
		t.Fatalf("BuildTree() = %#v, want one root with two children", tree)
	}
	if tree[0].Children[0].ID != "child-a" || tree[0].Children[1].ID != "child-b" {
		t.Fatalf("same-name child order = [%s %s], want [child-a child-b]", tree[0].Children[0].ID, tree[0].Children[1].ID)
	}
}

func TestBuildTreeSurfacesUnreachableUnitsWithoutRecursingForever(t *testing.T) {
	missingParentID := "missing-parent"
	cycleAID := "cycle-a"
	cycleBID := "cycle-b"
	units := []Unit{
		{ID: "root", Name: "Root"},
		{ID: "orphan", ParentID: &missingParentID, Name: "Orphan"},
		{ID: cycleAID, ParentID: &cycleBID, Name: "Cycle A"},
		{ID: cycleBID, ParentID: &cycleAID, Name: "Cycle B"},
	}

	tree := BuildTree(units)
	if len(tree) != 3 {
		t.Fatalf("BuildTree() root count = %d, want root, cycle component, and orphan: %#v", len(tree), tree)
	}
	if tree[0].ID != "root" {
		t.Fatalf("first root = %s, want root", tree[0].ID)
	}
	if tree[1].ID != "cycle-a" || len(tree[1].Children) != 1 || tree[1].Children[0].ID != "cycle-b" {
		t.Fatalf("cycle component = %#v, want cycle-a with cycle-b child", tree[1])
	}
	if tree[2].ID != "orphan" {
		t.Fatalf("third surfaced unit = %s, want orphan", tree[2].ID)
	}
}

func TestDescendantIDsReturnsDescendantsOnly(t *testing.T) {
	rootID := "root"
	childID := "child"
	grandchildID := "grandchild"
	unrelatedParentID := "unrelated-parent"

	units := []Unit{
		{ID: rootID, Name: "Root"},
		{ID: childID, ParentID: &rootID, Name: "Child"},
		{ID: grandchildID, ParentID: &childID, Name: "Grandchild"},
		{ID: "unrelated-child", ParentID: &unrelatedParentID, Name: "Unrelated child"},
		{ID: unrelatedParentID, Name: "Unrelated parent"},
	}

	got := DescendantIDs(rootID, units)
	if !got[childID] || !got[grandchildID] {
		t.Fatalf("DescendantIDs() = %v, want child and grandchild", got)
	}
	if got[rootID] {
		t.Fatalf("DescendantIDs() includes root %q", rootID)
	}
	if got["unrelated-child"] || got[unrelatedParentID] {
		t.Fatalf("DescendantIDs() includes unrelated units: %v", got)
	}
}
