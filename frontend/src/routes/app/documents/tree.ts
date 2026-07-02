export type DocumentFolderNode = {
	id: string;
	parentId?: string | null;
	organizationUnitId: string;
	name: string;
	description?: string | null;
	deletedAt?: string | null;
	createdAt: string;
	updatedAt: string;
	children: DocumentFolderNode[];
};

export type FlatDocumentFolderNode = Omit<DocumentFolderNode, 'children'> & {
	depth: number;
};

export type RepositoryDocument = {
	id: string;
	folderId: string;
	organizationUnitId: string;
	originalFileName: string;
	displayName: string;
	mimeType: string;
	extension?: string | null;
	sizeBytes: number;
	sha256Hash: string;
	deletedAt?: string | null;
	createdAt: string;
	updatedAt: string;
};

export type RepositoryRow =
	| {
			type: 'folder';
			id: string;
			name: string;
			folder: DocumentFolderNode;
	  }
	| {
			type: 'document';
			id: string;
			name: string;
			document: RepositoryDocument;
	  };

export function flattenFolderTree(
	folders: DocumentFolderNode[],
	depth = 0
): FlatDocumentFolderNode[] {
	return folders.flatMap((folder) => {
		const { children, ...row } = folder;

		return [{ ...row, depth }, ...flattenFolderTree(children, depth + 1)];
	});
}

export function findFolder(
	folders: DocumentFolderNode[],
	folderId: string | null | undefined
): DocumentFolderNode | null {
	if (!folderId) {
		return null;
	}

	for (const folder of folders) {
		if (folder.id === folderId) {
			return folder;
		}

		const child = findFolder(folder.children, folderId);
		if (child) {
			return child;
		}
	}

	return null;
}

export function selectInitialFolder(
	folders: DocumentFolderNode[],
	requestedId?: string | null
): DocumentFolderNode | null {
	return findFolder(folders, requestedId) ?? folders[0] ?? null;
}

export function collectFolderMoveTargets(
	folders: DocumentFolderNode[],
	selectedId: string
): FlatDocumentFolderNode[] {
	const rows: FlatDocumentFolderNode[] = [];

	for (const folder of folders) {
		collectFolderMoveTargetRows(folder, selectedId, 0, rows);
	}

	return rows;
}

export function repositoryRows(
	folders: DocumentFolderNode[],
	documents: RepositoryDocument[]
): RepositoryRow[] {
	const folderRows: RepositoryRow[] = folders
		.map((folder) => ({
			type: 'folder' as const,
			id: folder.id,
			name: folder.name,
			folder
		}))
		.sort(compareRepositoryRows);
	const documentRows: RepositoryRow[] = documents
		.map((document) => ({
			type: 'document' as const,
			id: document.id,
			name: document.displayName,
			document
		}))
		.sort(compareRepositoryRows);

	return [...folderRows, ...documentRows];
}

function collectFolderMoveTargetRows(
	folder: DocumentFolderNode,
	selectedId: string,
	depth: number,
	rows: FlatDocumentFolderNode[]
) {
	if (folder.id === selectedId || isDeleted(folder)) {
		return;
	}

	const { children, ...row } = folder;
	rows.push({ ...row, depth });

	for (const child of children) {
		collectFolderMoveTargetRows(child, selectedId, depth + 1, rows);
	}
}

function compareRepositoryRows(first: RepositoryRow, second: RepositoryRow) {
	return compareText(first.name, second.name) || compareText(first.id, second.id);
}

function compareText(first: string, second: string) {
	const normalizedFirst = first.toLowerCase();
	const normalizedSecond = second.toLowerCase();

	if (normalizedFirst < normalizedSecond) return -1;
	if (normalizedFirst > normalizedSecond) return 1;
	if (first < second) return -1;
	if (first > second) return 1;
	return 0;
}

function isDeleted(folder: DocumentFolderNode) {
	return typeof folder.deletedAt === 'string' && folder.deletedAt.trim() !== '';
}
