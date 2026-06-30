package api

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	ocrsvc "ai.ro/syncra/internal/ocr"
)

type createDatasetRequest struct {
	Name           string                `json:"name"`
	UserID         string                `json:"user_id"`
	SchemaID       string                `json:"schema_id"`
	SelectedFields []ocrsvc.DatasetField `json:"selected_fields"`
}

type updateDatasetRequest struct {
	Name           string                `json:"name"`
	SchemaID       string                `json:"schema_id"`
	SelectedFields []ocrsvc.DatasetField `json:"selected_fields"`
}

type datasetListCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Sort      string    `json:"sort"`
}

type datasetDateBounds struct {
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

const maxDatasetNameCharacters = 160
const maxDatasetFieldCount = 100
const maxDatasetRequestBytes int64 = 1 << 20

var (
	errInvalidDatasetSchemaID       = errors.New("invalid schema_id")
	errInvalidDatasetSelectedFields = errors.New("invalid selected_fields")
	errDatasetNotFound              = errors.New("dataset not found")
	errDatasetSchemaUnavailable     = errors.New("dataset schema unavailable")
)

func validateDatasetName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errors.New("name is required")
	}
	if utf8.RuneCountInString(name) > maxDatasetNameCharacters {
		return "", errors.New("name must be at most 160 characters")
	}
	return name, nil
}

func parseDatasetID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("invalid dataset id")
	}
	return id, nil
}

func parseDatasetSchemaID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errInvalidDatasetSchemaID
	}
	return id, nil
}

func parseDatasetDateBounds(c *gin.Context) (datasetDateBounds, error) {
	createdFrom, err := parseOCRJobTimeQuery(c, "created_from")
	if err != nil {
		return datasetDateBounds{}, errors.New("invalid created_from")
	}
	createdTo, err := parseOCRJobTimeQuery(c, "created_to")
	if err != nil {
		return datasetDateBounds{}, errors.New("invalid created_to")
	}
	if createdFrom != nil && createdTo != nil && createdFrom.After(*createdTo) {
		return datasetDateBounds{}, errors.New("created_from must be before or equal to created_to")
	}

	return datasetDateBounds{CreatedFrom: createdFrom, CreatedTo: createdTo}, nil
}

func validateDatasetSchemaBelongsToUser(db *gorm.DB, userID string, schemaID uuid.UUID) (ocrsvc.ExtractionSchema, error) {
	var schema ocrsvc.ExtractionSchema
	err := db.Where("id = ? AND user_id = ?", schemaID, userID).First(&schema).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ocrsvc.ExtractionSchema{}, errInvalidDatasetSchemaID
		}
		return ocrsvc.ExtractionSchema{}, err
	}
	return schema, nil
}

func validateDatasetSelectedFields(schema ocrsvc.ExtractionSchema, fields []ocrsvc.DatasetField) error {
	if len(fields) > maxDatasetFieldCount {
		return errInvalidDatasetSelectedFields
	}
	if err := ocrsvc.ValidateDatasetFields(json.RawMessage(schema.SchemaJSON), fields); err != nil {
		return errInvalidDatasetSelectedFields
	}
	return nil
}

func datasetBindErrorMessage(err error) string {
	if strings.Contains(err.Error(), "request body too large") {
		return "request body too large"
	}
	return "invalid JSON body"
}

func marshalDatasetSelectedFields(fields []ocrsvc.DatasetField) (datatypes.JSON, error) {
	rawFields, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(rawFields), nil
}

func (h *Handler) CreateDataset(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxDatasetRequestBytes)

	var req createDatasetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, datasetBindErrorMessage(err))
		return
	}
	userID, err := parseRequiredUserID(req.UserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	name, err := validateDatasetName(req.Name)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	schemaID, err := parseDatasetSchemaID(req.SchemaID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	db := h.DB.WithContext(c.Request.Context())
	if err := validateCollectionUserExists(db, userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}
	schema, err := validateDatasetSchemaBelongsToUser(db, userID, schemaID)
	if err != nil {
		if errors.Is(err, errInvalidDatasetSchemaID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate schema")
		return
	}
	if err := validateDatasetSelectedFields(schema, req.SelectedFields); err != nil {
		writeError(c, http.StatusBadRequest, errInvalidDatasetSelectedFields.Error())
		return
	}
	rawFields, err := marshalDatasetSelectedFields(req.SelectedFields)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to encode selected fields")
		return
	}

	dataset := ocrsvc.Dataset{
		UserID:         userID,
		SchemaID:       schema.ID,
		Name:           name,
		SelectedFields: rawFields,
	}
	if err := db.Create(&dataset).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save dataset")
		return
	}

	c.JSON(http.StatusCreated, makeDatasetResponse(dataset, schema.Name, req.SelectedFields))
}

func (h *Handler) ListDatasets(c *gin.Context) {
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	sortDirection, err := parseOCRJobListSort(c.Query("sort"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	size, err := parseOCRJobListSize(c.Query("size"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	cursor, err := parseDatasetListCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	query := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID)
	if cursor != nil {
		operator := "<"
		if sortDirection == "asc" {
			operator = ">"
		}
		query = query.Where("(created_at, id) "+operator+" (?, ?)", cursor.CreatedAt, cursor.ID)
	}

	order := "created_at desc, id desc"
	if sortDirection == "asc" {
		order = "created_at asc, id asc"
	}
	var datasets []ocrsvc.Dataset
	if err := query.Order(order).Limit(size + 1).Find(&datasets).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list datasets")
		return
	}

	var nextCursor *string
	if len(datasets) > size {
		datasets = datasets[:size]
		if len(datasets) > 0 {
			encoded, err := encodeDatasetListCursor(datasets[len(datasets)-1], sortDirection)
			if err != nil {
				writeError(c, http.StatusInternalServerError, "failed to encode next cursor")
				return
			}
			nextCursor = &encoded
		}
	}

	out, err := h.datasetResponses(c.Request.Context(), datasets, userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load datasets")
		return
	}
	c.JSON(http.StatusOK, DatasetListResponse{Datasets: out, NextCursor: nextCursor})
}

func (h *Handler) GetDataset(c *gin.Context) {
	id, err := parseDatasetID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var dataset ocrsvc.Dataset
	if err := h.DB.WithContext(c.Request.Context()).
		Where("id = ? AND user_id = ?", id, userID).
		First(&dataset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "dataset not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load dataset")
		return
	}

	resp, err := h.datasetResponse(c.Request.Context(), dataset)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load dataset")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetDatasetRows(c *gin.Context) {
	id, err := parseDatasetID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	sortDirection, err := parseOCRJobListSort(c.Query("sort"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	size, err := parseOCRJobListSize(c.Query("size"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	cursor, err := parseOCRDocumentListCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	dateBounds, err := parseDatasetDateBounds(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	dataset, response, fields, err := h.loadDatasetForUser(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, errDatasetNotFound) {
			writeError(c, http.StatusNotFound, errDatasetNotFound.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load dataset")
		return
	}
	rows, nextCursor, err := h.datasetRows(c.Request.Context(), dataset, fields, size, sortDirection, cursor, dateBounds)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load dataset rows")
		return
	}

	c.JSON(http.StatusOK, DatasetRowsResponse{
		Dataset:    response,
		Columns:    datasetColumnResponses(fields),
		Rows:       rows,
		NextCursor: nextCursor,
	})
}

func (h *Handler) ExportDataset(c *gin.Context) {
	id, err := parseDatasetID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	format := strings.ToLower(strings.TrimSpace(c.DefaultQuery("format", "csv")))
	if format != "csv" && format != "xlsx" {
		writeError(c, http.StatusBadRequest, "format must be csv or xlsx")
		return
	}
	sortDirection, err := parseOCRJobListSort(c.Query("sort"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	dateBounds, err := parseDatasetDateBounds(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	dataset, _, fields, err := h.loadDatasetForUser(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, errDatasetNotFound) {
			writeError(c, http.StatusNotFound, errDatasetNotFound.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load dataset")
		return
	}
	columns := datasetColumnResponses(fields)
	switch format {
	case "csv":
		file, err := h.datasetCSVFile(c.Request.Context(), dataset, columns, fields, sortDirection, dateBounds)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to export dataset")
			return
		}
		defer func() {
			name := file.Name()
			_ = file.Close()
			_ = os.Remove(name)
		}()
		c.Header("Content-Disposition", `attachment; filename="`+datasetExportFilename(dataset, format)+`"`)
		c.Header("Content-Type", "text/csv; charset=utf-8")
		if _, err := io.Copy(c.Writer, file); err != nil {
			_ = c.Error(err)
		}
	case "xlsx":
		file, err := h.datasetXLSXFile(c.Request.Context(), dataset, columns, fields, sortDirection, dateBounds)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to export dataset")
			return
		}
		defer func() {
			_ = file.Close()
		}()
		c.Header("Content-Disposition", `attachment; filename="`+datasetExportFilename(dataset, format)+`"`)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		if err := file.Write(c.Writer); err != nil {
			_ = c.Error(err)
		}
	}
}

func (h *Handler) UpdateDataset(c *gin.Context) {
	id, err := parseDatasetID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxDatasetRequestBytes)
	var req updateDatasetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, datasetBindErrorMessage(err))
		return
	}
	name, err := validateDatasetName(req.Name)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	schemaID, err := parseDatasetSchemaID(req.SchemaID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var dataset ocrsvc.Dataset
	var schema ocrsvc.ExtractionSchema
	db := h.DB.WithContext(c.Request.Context())
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&dataset).Error; err != nil {
			return err
		}
		loadedSchema, err := validateDatasetSchemaBelongsToUser(tx, userID, schemaID)
		if err != nil {
			return err
		}
		if err := validateDatasetSelectedFields(loadedSchema, req.SelectedFields); err != nil {
			return err
		}
		rawFields, err := marshalDatasetSelectedFields(req.SelectedFields)
		if err != nil {
			return err
		}

		result := tx.Model(&ocrsvc.Dataset{}).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(map[string]any{
				"name":            name,
				"schema_id":       loadedSchema.ID,
				"selected_fields": rawFields,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&dataset).Error; err != nil {
			return err
		}
		schema = loadedSchema
		return nil
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "dataset not found")
			return
		}
		if errors.Is(err, errInvalidDatasetSchemaID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, errInvalidDatasetSelectedFields) {
			writeError(c, http.StatusBadRequest, errInvalidDatasetSelectedFields.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to update dataset")
		return
	}

	c.JSON(http.StatusOK, makeDatasetResponse(dataset, schema.Name, req.SelectedFields))
}

func (h *Handler) DeleteDataset(c *gin.Context) {
	id, err := parseDatasetID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result := h.DB.WithContext(c.Request.Context()).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&ocrsvc.Dataset{})
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete dataset")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "dataset not found")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) loadDatasetForUser(ctx context.Context, id uuid.UUID, userID string) (ocrsvc.Dataset, DatasetResponse, []ocrsvc.DatasetField, error) {
	var dataset ocrsvc.Dataset
	if err := h.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&dataset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ocrsvc.Dataset{}, DatasetResponse{}, nil, errDatasetNotFound
		}
		return ocrsvc.Dataset{}, DatasetResponse{}, nil, err
	}

	var schema ocrsvc.ExtractionSchema
	if err := h.DB.WithContext(ctx).
		Select("id", "name").
		Where("id = ? AND user_id = ?", dataset.SchemaID, userID).
		First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ocrsvc.Dataset{}, DatasetResponse{}, nil, errDatasetSchemaUnavailable
		}
		return ocrsvc.Dataset{}, DatasetResponse{}, nil, err
	}

	fields, err := datasetFieldsFromJSON(dataset.SelectedFields)
	if err != nil {
		return ocrsvc.Dataset{}, DatasetResponse{}, nil, err
	}
	return dataset, makeDatasetResponse(dataset, schema.Name, fields), fields, nil
}

func (h *Handler) datasetRows(ctx context.Context, dataset ocrsvc.Dataset, fields []ocrsvc.DatasetField, size int, sortDirection string, cursor *ocrDocumentListCursor, dateBounds datasetDateBounds) ([]DatasetRowResponse, *string, error) {
	query := h.datasetDocumentQuery(ctx, dataset, sortDirection, cursor, dateBounds)
	if size > 0 {
		query = query.Limit(size + 1)
	}

	var documents []ocrsvc.OCRDocument
	if err := query.Find(&documents).Error; err != nil {
		return nil, nil, err
	}

	var nextCursor *string
	if size > 0 && len(documents) > size {
		documents = documents[:size]
		if len(documents) > 0 {
			encoded, err := encodeOCRDocumentListCursor(documents[len(documents)-1], sortDirection)
			if err != nil {
				return nil, nil, err
			}
			nextCursor = &encoded
		}
	}

	rows := make([]DatasetRowResponse, 0, len(documents))
	for _, doc := range documents {
		values, err := ocrsvc.ProjectDatasetValues(json.RawMessage(doc.AnnotationJSON), fields)
		if err != nil {
			return nil, nil, err
		}
		rows = append(rows, DatasetRowResponse{
			DocumentID: doc.ID,
			Filename:   doc.OriginalFilename,
			CreatedAt:  doc.CreatedAt,
			Values:     values,
		})
	}
	return rows, nextCursor, nil
}

func (h *Handler) datasetDocumentQuery(ctx context.Context, dataset ocrsvc.Dataset, sortDirection string, cursor *ocrDocumentListCursor, dateBounds datasetDateBounds) *gorm.DB {
	query := h.DB.WithContext(ctx).
		Model(&ocrsvc.OCRDocument{}).
		Select("id", "original_filename", "created_at", "annotation_json").
		Where("user_id = ? AND schema_id = ?", dataset.UserID, dataset.SchemaID)
	if dateBounds.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *dateBounds.CreatedFrom)
	}
	if dateBounds.CreatedTo != nil {
		query = query.Where("created_at <= ?", *dateBounds.CreatedTo)
	}
	if cursor != nil {
		operator := "<"
		if sortDirection == "asc" {
			operator = ">"
		}
		query = query.Where("(created_at, id) "+operator+" (?, ?)", cursor.CreatedAt, cursor.ID)
	}

	order := "created_at desc, id desc"
	if sortDirection == "asc" {
		order = "created_at asc, id asc"
	}
	return query.Order(order)
}

func (h *Handler) forEachDatasetRow(ctx context.Context, dataset ocrsvc.Dataset, fields []ocrsvc.DatasetField, sortDirection string, dateBounds datasetDateBounds, fn func(DatasetRowResponse) error) error {
	query := h.datasetDocumentQuery(ctx, dataset, sortDirection, nil, dateBounds)
	rows, err := query.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var doc ocrsvc.OCRDocument
		if err := h.DB.ScanRows(rows, &doc); err != nil {
			return err
		}
		values, err := ocrsvc.ProjectDatasetValues(json.RawMessage(doc.AnnotationJSON), fields)
		if err != nil {
			return err
		}
		if err := fn(DatasetRowResponse{
			DocumentID: doc.ID,
			Filename:   doc.OriginalFilename,
			CreatedAt:  doc.CreatedAt,
			Values:     values,
		}); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (h *Handler) datasetCSVFile(ctx context.Context, dataset ocrsvc.Dataset, columns []DatasetColumnResponse, fields []ocrsvc.DatasetField, sortDirection string, dateBounds datasetDateBounds) (*os.File, error) {
	file, err := os.CreateTemp("", "syncra-dataset-export-*.csv")
	if err != nil {
		return nil, err
	}
	success := false
	defer func() {
		if !success {
			name := file.Name()
			_ = file.Close()
			_ = os.Remove(name)
		}
	}()

	if err := h.writeDatasetCSV(ctx, file, dataset, columns, fields, sortDirection, dateBounds); err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	success = true
	return file, nil
}

func (h *Handler) writeDatasetCSV(ctx context.Context, out io.Writer, dataset ocrsvc.Dataset, columns []DatasetColumnResponse, fields []ocrsvc.DatasetField, sortDirection string, dateBounds datasetDateBounds) error {
	writer := csv.NewWriter(out)
	header := datasetCSVHeader(columns)
	if err := writer.Write(header); err != nil {
		return err
	}

	if err := h.forEachDatasetRow(ctx, dataset, fields, sortDirection, dateBounds, func(row DatasetRowResponse) error {
		if err := writer.Write(datasetCSVRecord(columns, row)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

func (h *Handler) datasetXLSXFile(ctx context.Context, dataset ocrsvc.Dataset, columns []DatasetColumnResponse, fields []ocrsvc.DatasetField, sortDirection string, dateBounds datasetDateBounds) (*excelize.File, error) {
	file := excelize.NewFile()
	success := false
	defer func() {
		if !success {
			_ = file.Close()
		}
	}()

	sheet := sanitizeExcelSheetName(dataset.Name)
	if err := file.SetSheetName("Sheet1", sheet); err != nil {
		return nil, err
	}
	stream, err := file.NewStreamWriter(sheet)
	if err != nil {
		return nil, err
	}
	if err := stream.SetRow("A1", datasetXLSXHeader(columns)); err != nil {
		return nil, err
	}

	rowIndex := 2
	if err := h.forEachDatasetRow(ctx, dataset, fields, sortDirection, dateBounds, func(row DatasetRowResponse) error {
		cell, err := excelize.CoordinatesToCellName(1, rowIndex)
		if err != nil {
			return err
		}
		rowIndex++
		return stream.SetRow(cell, datasetXLSXRecord(columns, row))
	}); err != nil {
		return nil, err
	}
	if err := stream.Flush(); err != nil {
		return nil, err
	}
	success = true
	return file, nil
}

func datasetExportFilename(dataset ocrsvc.Dataset, format string) string {
	base := strings.ToLower(strings.TrimSpace(dataset.Name))
	var builder strings.Builder
	lastDash := false
	for _, r := range base {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			builder.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_':
			builder.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}
	name := strings.Trim(builder.String(), "-_")
	if name == "" {
		name = "dataset-" + dataset.ID.String()
	}
	return name + "." + format
}

func sanitizeExcelSheetName(raw string) string {
	name := strings.Map(func(r rune) rune {
		switch r {
		case ':', '\\', '/', '?', '*', '[', ']':
			return ' '
		}
		if r < 32 {
			return ' '
		}
		return r
	}, strings.TrimSpace(raw))
	name = normalizeExcelSheetNameBoundary(name)
	if name == "" {
		return "Dataset"
	}
	if utf8.RuneCountInString(name) <= 31 {
		return name
	}

	var builder strings.Builder
	count := 0
	for _, r := range name {
		if count >= 31 {
			break
		}
		builder.WriteRune(r)
		count++
	}
	name = normalizeExcelSheetNameBoundary(builder.String())
	if name == "" {
		return "Dataset"
	}
	return name
}

func (h *Handler) datasetResponse(ctx context.Context, dataset ocrsvc.Dataset) (DatasetResponse, error) {
	var schema ocrsvc.ExtractionSchema
	if err := h.DB.WithContext(ctx).
		Select("id", "name").
		Where("id = ? AND user_id = ?", dataset.SchemaID, dataset.UserID).
		First(&schema).Error; err != nil {
		return DatasetResponse{}, err
	}
	fields, err := datasetFieldsFromJSON(dataset.SelectedFields)
	if err != nil {
		return DatasetResponse{}, err
	}
	return makeDatasetResponse(dataset, schema.Name, fields), nil
}

func (h *Handler) datasetResponses(ctx context.Context, datasets []ocrsvc.Dataset, userID string) ([]DatasetResponse, error) {
	if len(datasets) == 0 {
		return []DatasetResponse{}, nil
	}

	schemaIDs := make([]uuid.UUID, 0, len(datasets))
	seen := make(map[uuid.UUID]struct{}, len(datasets))
	for _, dataset := range datasets {
		if _, ok := seen[dataset.SchemaID]; ok {
			continue
		}
		seen[dataset.SchemaID] = struct{}{}
		schemaIDs = append(schemaIDs, dataset.SchemaID)
	}

	var schemas []ocrsvc.ExtractionSchema
	if err := h.DB.WithContext(ctx).
		Select("id", "name").
		Where("id IN ? AND user_id = ?", schemaIDs, userID).
		Find(&schemas).Error; err != nil {
		return nil, err
	}
	schemaNames := make(map[uuid.UUID]string, len(schemas))
	for _, schema := range schemas {
		schemaNames[schema.ID] = schema.Name
	}

	responses := make([]DatasetResponse, 0, len(datasets))
	for _, dataset := range datasets {
		schemaName, ok := schemaNames[dataset.SchemaID]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		fields, err := datasetFieldsFromJSON(dataset.SelectedFields)
		if err != nil {
			return nil, err
		}
		responses = append(responses, makeDatasetResponse(dataset, schemaName, fields))
	}
	return responses, nil
}

func datasetFieldsFromJSON(raw datatypes.JSON) ([]ocrsvc.DatasetField, error) {
	if len(raw) == 0 {
		return []ocrsvc.DatasetField{}, nil
	}
	var fields []ocrsvc.DatasetField
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, err
	}
	if fields == nil {
		return []ocrsvc.DatasetField{}, nil
	}
	return fields, nil
}

func datasetColumnResponses(fields []ocrsvc.DatasetField) []DatasetColumnResponse {
	columns := make([]DatasetColumnResponse, 0, len(fields))
	for _, field := range fields {
		columns = append(columns, DatasetColumnResponse{
			Key:   field.Key,
			Label: field.Label,
			Path:  field.Path,
		})
	}
	return columns
}

func datasetExportHeader(columns []DatasetColumnResponse) []string {
	header := []string{"document_id", "filename", "created_at"}
	for _, column := range columns {
		header = append(header, column.Label)
	}
	return header
}

func datasetCSVHeader(columns []DatasetColumnResponse) []string {
	header := datasetExportHeader(columns)
	for i, value := range header {
		header[i] = sanitizeDatasetCSVCell(value)
	}
	return header
}

func datasetCSVRecord(columns []DatasetColumnResponse, row DatasetRowResponse) []string {
	record := []string{
		sanitizeDatasetCSVCell(row.DocumentID.String()),
		sanitizeDatasetCSVCell(row.Filename),
		sanitizeDatasetCSVCell(row.CreatedAt.UTC().Format(time.RFC3339Nano)),
	}
	for _, column := range columns {
		record = append(record, sanitizeDatasetCSVCell(datasetExportCellString(row.Values[column.Key])))
	}
	return record
}

func datasetXLSXHeader(columns []DatasetColumnResponse) []interface{} {
	header := datasetExportHeader(columns)
	values := make([]interface{}, 0, len(header))
	for _, value := range header {
		values = append(values, value)
	}
	return values
}

func datasetXLSXRecord(columns []DatasetColumnResponse, row DatasetRowResponse) []interface{} {
	values := []interface{}{
		row.DocumentID.String(),
		row.Filename,
		row.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
	for _, column := range columns {
		value := row.Values[column.Key]
		if value == nil {
			value = ""
		}
		values = append(values, value)
	}
	return values
}

func datasetExportCellString(value any) string {
	if value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}

func sanitizeDatasetCSVCell(value string) string {
	if value == "" {
		return value
	}
	switch value[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + value
	default:
		return value
	}
}

func normalizeExcelSheetNameBoundary(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "'")
	return strings.TrimSpace(name)
}

func makeDatasetResponse(dataset ocrsvc.Dataset, schemaName string, fields []ocrsvc.DatasetField) DatasetResponse {
	fieldResponses := make([]DatasetFieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResponses = append(fieldResponses, DatasetFieldResponse{
			Path:  field.Path,
			Key:   field.Key,
			Label: field.Label,
		})
	}
	return DatasetResponse{
		ID:             dataset.ID,
		CreatedAt:      dataset.CreatedAt,
		UpdatedAt:      dataset.UpdatedAt,
		UserID:         collectionUserIDString(dataset.UserID),
		Name:           dataset.Name,
		SchemaID:       dataset.SchemaID,
		SchemaName:     schemaName,
		SelectedFields: fieldResponses,
		FieldCount:     len(fieldResponses),
	}
}

func parseDatasetListCursor(raw string, sortDirection string) (*datasetListCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	var cursor datasetListCursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, errors.New("invalid cursor")
	}
	if cursor.ID == uuid.Nil || cursor.CreatedAt.IsZero() || (cursor.Sort != "asc" && cursor.Sort != "desc") {
		return nil, errors.New("invalid cursor")
	}
	if cursor.Sort != sortDirection {
		return nil, errors.New("cursor sort does not match sort")
	}
	return &cursor, nil
}

func encodeDatasetListCursor(dataset ocrsvc.Dataset, sortDirection string) (string, error) {
	raw, err := json.Marshal(datasetListCursor{
		CreatedAt: dataset.CreatedAt.UTC(),
		ID:        dataset.ID,
		Sort:      sortDirection,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
