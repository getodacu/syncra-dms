package orgunits

import (
	"sort"
	"time"
)

type TreeNode struct {
	ID          string     `json:"id"`
	ParentID    *string    `json:"parentId"`
	Name        string     `json:"name"`
	Code        *string    `json:"code"`
	Description *string    `json:"description"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
	Children    []TreeNode `json:"children"`
}

func BuildTree(units []Unit) []TreeNode {
	unitsByParent := make(map[string][]Unit)
	roots := make([]Unit, 0)

	for _, unit := range units {
		if unit.ParentID == nil {
			roots = append(roots, unit)
			continue
		}
		unitsByParent[*unit.ParentID] = append(unitsByParent[*unit.ParentID], unit)
	}

	sortUnits(roots)
	for parentID := range unitsByParent {
		sortUnits(unitsByParent[parentID])
	}

	visited := map[string]bool{}
	var buildChildren func(parentID string) []TreeNode
	buildChildren = func(parentID string) []TreeNode {
		children := unitsByParent[parentID]
		nodes := make([]TreeNode, 0, len(children))
		for _, child := range children {
			if visited[child.ID] {
				continue
			}
			visited[child.ID] = true
			nodes = append(nodes, unitToTreeNode(child, buildChildren(child.ID)))
		}
		return nodes
	}

	nodes := make([]TreeNode, 0, len(roots))
	for _, root := range roots {
		if visited[root.ID] {
			continue
		}
		visited[root.ID] = true
		nodes = append(nodes, unitToTreeNode(root, buildChildren(root.ID)))
	}
	remaining := make([]Unit, 0)
	for _, unit := range units {
		if !visited[unit.ID] {
			remaining = append(remaining, unit)
		}
	}
	sortUnits(remaining)
	for _, unit := range remaining {
		if visited[unit.ID] {
			continue
		}
		visited[unit.ID] = true
		nodes = append(nodes, unitToTreeNode(unit, buildChildren(unit.ID)))
	}
	return nodes
}

func DescendantIDs(rootID string, units []Unit) map[string]bool {
	unitsByParent := make(map[string][]Unit)
	for _, unit := range units {
		if unit.ParentID != nil {
			unitsByParent[*unit.ParentID] = append(unitsByParent[*unit.ParentID], unit)
		}
	}

	descendants := make(map[string]bool)
	visited := map[string]bool{rootID: true}

	var walk func(parentID string)
	walk = func(parentID string) {
		for _, child := range unitsByParent[parentID] {
			if visited[child.ID] {
				continue
			}
			visited[child.ID] = true
			descendants[child.ID] = true
			walk(child.ID)
		}
	}
	walk(rootID)

	return descendants
}

func sortUnits(units []Unit) {
	sort.Slice(units, func(i, j int) bool {
		if units[i].Name == units[j].Name {
			return units[i].ID < units[j].ID
		}
		return units[i].Name < units[j].Name
	})
}

func unitToTreeNode(unit Unit, children []TreeNode) TreeNode {
	return TreeNode{
		ID:          unit.ID,
		ParentID:    unit.ParentID,
		Name:        unit.Name,
		Code:        unit.Code,
		Description: unit.Description,
		CreatedAt:   formatTreeTime(unit.CreatedAt),
		UpdatedAt:   formatTreeTime(unit.UpdatedAt),
		Children:    children,
	}
}

func formatTreeTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}
