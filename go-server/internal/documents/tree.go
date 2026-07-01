package documents

import (
	"sort"
	"time"
)

type FolderTreeNode struct {
	ID          string           `json:"id"`
	ParentID    *string          `json:"parentId"`
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
	Children    []FolderTreeNode `json:"children"`
}

func BuildFolderTree(folders []Folder) []FolderTreeNode {
	foldersByParent := make(map[string][]Folder)
	roots := make([]Folder, 0)

	for _, folder := range folders {
		if folder.ParentID == nil {
			roots = append(roots, folder)
			continue
		}
		foldersByParent[*folder.ParentID] = append(foldersByParent[*folder.ParentID], folder)
	}

	sortFolders(roots)
	for parentID := range foldersByParent {
		sortFolders(foldersByParent[parentID])
	}

	visited := map[string]bool{}
	var buildChildren func(parentID string) []FolderTreeNode
	buildChildren = func(parentID string) []FolderTreeNode {
		children := foldersByParent[parentID]
		nodes := make([]FolderTreeNode, 0, len(children))
		for _, child := range children {
			if visited[child.ID] {
				continue
			}
			visited[child.ID] = true
			nodes = append(nodes, folderToTreeNode(child, buildChildren(child.ID)))
		}
		return nodes
	}

	nodes := make([]FolderTreeNode, 0, len(roots))
	for _, root := range roots {
		if visited[root.ID] {
			continue
		}
		visited[root.ID] = true
		nodes = append(nodes, folderToTreeNode(root, buildChildren(root.ID)))
	}
	remaining := make([]Folder, 0)
	for _, folder := range folders {
		if !visited[folder.ID] {
			remaining = append(remaining, folder)
		}
	}
	sortFolders(remaining)
	for _, folder := range remaining {
		if visited[folder.ID] {
			continue
		}
		visited[folder.ID] = true
		nodes = append(nodes, folderToTreeNode(folder, buildChildren(folder.ID)))
	}
	return nodes
}

func DescendantFolderIDs(rootID string, folders []Folder) map[string]bool {
	foldersByParent := make(map[string][]Folder)
	for _, folder := range folders {
		if folder.ParentID != nil {
			foldersByParent[*folder.ParentID] = append(foldersByParent[*folder.ParentID], folder)
		}
	}

	descendants := make(map[string]bool)
	visited := map[string]bool{rootID: true}

	var walk func(parentID string)
	walk = func(parentID string) {
		for _, child := range foldersByParent[parentID] {
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

func sortFolders(folders []Folder) {
	sort.Slice(folders, func(i, j int) bool {
		if folders[i].Name == folders[j].Name {
			return folders[i].ID < folders[j].ID
		}
		return folders[i].Name < folders[j].Name
	})
}

func folderToTreeNode(folder Folder, children []FolderTreeNode) FolderTreeNode {
	return FolderTreeNode{
		ID:          folder.ID,
		ParentID:    folder.ParentID,
		Name:        folder.Name,
		Description: folder.Description,
		CreatedAt:   formatFolderTreeTime(folder.CreatedAt),
		UpdatedAt:   formatFolderTreeTime(folder.UpdatedAt),
		Children:    children,
	}
}

func formatFolderTreeTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}
