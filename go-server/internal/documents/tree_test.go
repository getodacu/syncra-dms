package documents

import (
	"testing"
	"time"
)

func TestBuildFolderTreeOrdersAndNestsFolders(t *testing.T) {
	rootID := "root"
	childID := "child"
	nodes := BuildFolderTree([]Folder{
		{ID: childID, ParentID: &rootID, Name: "Beta", CreatedAt: time.Unix(2, 0), UpdatedAt: time.Unix(2, 0)},
		{ID: rootID, Name: "Root", CreatedAt: time.Unix(1, 0), UpdatedAt: time.Unix(1, 0)},
		{ID: "alpha", ParentID: &rootID, Name: "Alpha", CreatedAt: time.Unix(3, 0), UpdatedAt: time.Unix(3, 0)},
	})
	if len(nodes) != 1 || len(nodes[0].Children) != 2 || nodes[0].Children[0].ID != "alpha" || nodes[0].Children[1].ID != childID {
		t.Fatalf("tree = %#v", nodes)
	}
}

func TestBuildFolderTreeOrdersEqualNamesByID(t *testing.T) {
	rootID := "root"
	nodes := BuildFolderTree([]Folder{
		{ID: rootID, Name: "Root"},
		{ID: "child-b", ParentID: &rootID, Name: "Same"},
		{ID: "child-a", ParentID: &rootID, Name: "Same"},
	})

	if len(nodes) != 1 || len(nodes[0].Children) != 2 {
		t.Fatalf("BuildFolderTree() = %#v, want one root with two children", nodes)
	}
	if nodes[0].Children[0].ID != "child-a" || nodes[0].Children[1].ID != "child-b" {
		t.Fatalf("same-name child order = [%s %s], want [child-a child-b]", nodes[0].Children[0].ID, nodes[0].Children[1].ID)
	}
}

func TestBuildFolderTreeSurfacesUnreachableFoldersWithoutRecursingForever(t *testing.T) {
	missingParentID := "missing-parent"
	cycleAID := "cycle-a"
	cycleBID := "cycle-b"
	nodes := BuildFolderTree([]Folder{
		{ID: "root", Name: "Root"},
		{ID: "orphan", ParentID: &missingParentID, Name: "Orphan"},
		{ID: cycleAID, ParentID: &cycleBID, Name: "Cycle A"},
		{ID: cycleBID, ParentID: &cycleAID, Name: "Cycle B"},
	})

	if len(nodes) != 3 {
		t.Fatalf("BuildFolderTree() root count = %d, want root, cycle component, and orphan: %#v", len(nodes), nodes)
	}
	if nodes[0].ID != "root" {
		t.Fatalf("first root = %s, want root", nodes[0].ID)
	}
	if nodes[1].ID != cycleAID || len(nodes[1].Children) != 1 || nodes[1].Children[0].ID != cycleBID {
		t.Fatalf("cycle component = %#v, want cycle-a with cycle-b child", nodes[1])
	}
	if nodes[2].ID != "orphan" {
		t.Fatalf("third surfaced folder = %s, want orphan", nodes[2].ID)
	}
}

func TestDescendantFolderIDsReturnsDescendantsOnly(t *testing.T) {
	rootID := "root"
	childID := "child"
	grandchildID := "grandchild"
	unrelatedParentID := "unrelated-parent"

	got := DescendantFolderIDs(rootID, []Folder{
		{ID: rootID, Name: "Root"},
		{ID: childID, ParentID: &rootID, Name: "Child"},
		{ID: grandchildID, ParentID: &childID, Name: "Grandchild"},
		{ID: "unrelated-child", ParentID: &unrelatedParentID, Name: "Unrelated child"},
		{ID: unrelatedParentID, Name: "Unrelated parent"},
	})

	if !got[childID] || !got[grandchildID] {
		t.Fatalf("DescendantFolderIDs() = %v, want child and grandchild", got)
	}
	if got[rootID] {
		t.Fatalf("DescendantFolderIDs() includes root %q", rootID)
	}
	if got["unrelated-child"] || got[unrelatedParentID] {
		t.Fatalf("DescendantFolderIDs() includes unrelated folders: %v", got)
	}
}
