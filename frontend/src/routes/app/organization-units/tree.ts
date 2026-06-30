export type OrganizationUnitNode = {
	id: string;
	parentId?: string | null;
	name: string;
	code?: string | null;
	description?: string | null;
	isArchived?: boolean;
	createdAt: string;
	updatedAt: string;
	children: OrganizationUnitNode[];
};

export type FlatOrganizationUnitNode = Omit<OrganizationUnitNode, 'children'> & {
	depth: number;
	isArchived?: boolean;
};

export function flattenUnitTree(
	units: OrganizationUnitNode[],
	depth = 0
): FlatOrganizationUnitNode[] {
	return units.flatMap((unit) => {
		const { children, ...row } = unit;

		return [{ ...row, depth }, ...flattenUnitTree(children, depth + 1)];
	});
}

export function findUnit(
	units: OrganizationUnitNode[],
	unitId: string | null | undefined
): OrganizationUnitNode | null {
	if (!unitId) {
		return null;
	}

	for (const unit of units) {
		if (unit.id === unitId) {
			return unit;
		}

		const child = findUnit(unit.children, unitId);
		if (child) {
			return child;
		}
	}

	return null;
}

export function countUnits(units: OrganizationUnitNode[]): number {
	return units.reduce((total, unit) => total + 1 + countUnits(unit.children), 0);
}

export function selectInitialUnit(
	units: OrganizationUnitNode[],
	requestedId?: string | null
): OrganizationUnitNode | null {
	return findUnit(units, requestedId) ?? units[0] ?? null;
}

export function collectMoveTargets(
	units: OrganizationUnitNode[],
	selectedId: string
): FlatOrganizationUnitNode[] {
	const rows: FlatOrganizationUnitNode[] = [];

	for (const unit of units) {
		collectMoveTargetRows(unit, selectedId, 0, rows);
	}

	return rows;
}

function collectMoveTargetRows(
	unit: OrganizationUnitNode,
	selectedId: string,
	depth: number,
	rows: FlatOrganizationUnitNode[]
) {
	if (unit.id === selectedId || isArchived(unit)) {
		return;
	}

	const { children, ...row } = unit;
	rows.push({ ...row, depth });

	for (const child of children) {
		collectMoveTargetRows(child, selectedId, depth + 1, rows);
	}
}

function isArchived(unit: OrganizationUnitNode) {
	return unit.isArchived === true;
}
