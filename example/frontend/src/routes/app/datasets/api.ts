import {
	exportDataset as clientExportDataset,
	type DatasetExportOptions,
} from "$lib/client/datasets";

export {
	createDataset,
	deleteDataset,
	fetchDatasetRows,
	fetchDatasets,
	getDataset,
	parseContentDispositionFilename,
	updateDataset,
} from "$lib/client/datasets";
export { clientExportDataset as exportDataset };
export type {
	CreateDatasetInput,
	DatasetColumnResponse,
	DatasetExportOptions,
	DatasetExportResponse,
	DatasetField,
	DatasetListQuery,
	DatasetListResponse,
	DatasetResponse,
	DatasetRowResponse,
	DatasetRowsQuery,
	DatasetRowsResponse,
	JsonValue,
	UpdateDatasetInput,
} from "$lib/client/datasets";

export type DatasetExportFormat = "csv" | "xlsx";
export type DatasetDownloadExportOptions = Omit<DatasetExportOptions, "format"> & {
	format?: never;
};

export function downloadDatasetExport(
	fetchFn: typeof fetch,
	id: string,
	format: DatasetExportFormat,
	options: DatasetDownloadExportOptions = {}
) {
	return clientExportDataset(fetchFn, id, { ...options, format });
}
