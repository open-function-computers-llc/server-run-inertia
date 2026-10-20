package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-function-computers-llc/server-run-inertia/account"
	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultRowLimit = 100
	maxRowLimit     = 5000
)

// columnMeta is the server-trusted metadata for one column, sourced from
// information_schema (never from client input). The Vue front-end uses it to
// render edit forms (type hints, nullable => empty means NULL, PKs readonly,
// auto-increment skipped on insert).
type columnMeta struct {
	Name          string `json:"name"`
	DataType      string `json:"type"`
	ColumnType    string `json:"columnType"`
	Nullable      bool   `json:"nullable"`
	IsPK          bool   `json:"isPk"`
	AutoIncrement bool   `json:"autoIncrement"`
}

// openDBForAccount opens a MySQL connection to the named database, but ONLY
// if that database actually belongs to the given account. This is the
// security boundary for the whole database-admin feature: an operator can
// never reach a database that isn't listed in the account's settings.json.
func (s *server) openDBForAccount(accountName, dbName string) (*sql.DB, error) {
	// The account must be registered with this server (boot snapshot) — the
	// page routes list accounts from disk, but only registered ones are
	// reachable here. We then re-read settings.json via account.Find so
	// freshly-added databases work without a server restart; the security
	// rule (db must be listed in the account's settings.json) is unchanged.
	registered := false
	for _, a := range s.accounts {
		if a.Name == accountName {
			registered = true
			break
		}
	}
	if !registered {
		return nil, fmt.Errorf("account %q is not registered with this server", accountName)
	}

	a, err := account.Find(accountName)
	if err != nil {
		return nil, err
	}
	for _, db := range a.Databases {
		if db.Name == dbName {
			// parseTime is deliberately OFF. With it on, the driver turns
			// DATE/DATETIME/TIME into time.Time, which JSON-marshals as RFC3339
			// with a bogus Z ("2024-03-15T00:00:00Z"); the client would then do
			// timezone math on a value that has no timezone, shifting it on save,
			// and MySQL TIME values outside 00:00:00-23:59:59 cannot round-trip
			// at all. Left off, the driver returns these as []byte holding MySQL's
			// own canonical text ("2024-03-15", "08:30:00", "2024-03-15 00:00:00"),
			// which normalizeRowValues turns into plain strings, so what the UI
			// reads is exactly what it writes back — no reinterpretation.
			dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=10s&readTimeout=10s",
				db.Username, db.Password, db.Host, db.Name)
			conn, err := sql.Open("mysql", dsn)
			if err != nil {
				return nil, err
			}
			// verify the connection actually works
			if err = conn.Ping(); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		}
	}
	return nil, fmt.Errorf("database %q is not configured for account %q", dbName, accountName)
}

// fetchColumnInfo returns the server-trusted column list for a table in the
// given schema. An empty slice means the table doesn't exist in this schema.
func fetchColumnInfo(conn *sql.DB, ctx context.Context, schema, table string) ([]columnMeta, error) {
	rows, err := conn.QueryContext(ctx,
		"SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA, COLUMN_TYPE "+
			"FROM information_schema.COLUMNS "+
			"WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? "+
			"ORDER BY ORDINAL_POSITION", schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := []columnMeta{}
	for rows.Next() {
		var name, dataType, nullable, key, extra, columnType string
		if err := rows.Scan(&name, &dataType, &nullable, &key, &extra, &columnType); err != nil {
			return nil, err
		}
		cols = append(cols, columnMeta{
			Name:          name,
			DataType:      dataType,
			ColumnType:    columnType,
			Nullable:      strings.EqualFold(nullable, "YES"),
			IsPK:          strings.EqualFold(key, "PRI"),
			AutoIncrement: strings.Contains(strings.ToLower(extra), "auto_increment"),
		})
	}
	return cols, rows.Err()
}

// normalizeRowValues converts driver-returned []byte to string. The mysql
// driver yields []byte for BINARY/VARBINARY/BLOB and some text columns and
// encoding/json would marshal those as base64, showing garbage in the UI.
func normalizeRowValues(values []any) {
	for i, v := range values {
		if b, ok := v.([]byte); ok {
			values[i] = string(b)
		}
	}
}

// handleDbTables lists the tables (with row counts and engine) for a database
// belonging to the given account.
// GET /api/accounts/{account}/databases/{db}/tables
func (s *server) handleDbTables() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("account")
		dbName := r.PathValue("db")

		if accountName == "" || dbName == "" {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "account and db are required",
			})
			return
		}

		conn, err := s.openDBForAccount(accountName, dbName)
		if err != nil {
			sendJSONError(w, http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}
		defer conn.Close()

		ctx := r.Context()

		// information_schema is the standard, read-only way to introspect.
		rows, err := conn.QueryContext(ctx,
			"SELECT TABLE_NAME, ENGINE, TABLE_ROWS, TABLE_COMMENT "+
				"FROM information_schema.TABLES "+
				"WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = 'BASE TABLE' "+
				"ORDER BY TABLE_NAME", dbName)
		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": "Error querying tables: " + err.Error(),
			})
			return
		}
		defer rows.Close()

		type tableInfo struct {
			Name    string `json:"name"`
			Engine  string `json:"engine"`
			Rows    int64  `json:"rows"`
			Comment string `json:"comment"`
		}
		tables := []tableInfo{}
		for rows.Next() {
			var t tableInfo
			if err := rows.Scan(&t.Name, &t.Engine, &t.Rows, &t.Comment); err != nil {
				// some engines leave rows/comment NULL; don't abort the list
				continue
			}
			tables = append(tables, t)
		}
		if err = rows.Err(); err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		sendJSON(w, map[string]any{"tables": tables})
	}
}

// handleDbTableData returns the columns, per-column metadata, primary key
// columns, one LIMIT/OFFSET page of rows, and the exact row total for a
// single table. The server builds the SELECT itself: identifiers come only
// from its own information_schema queries (the `sort` column is validated
// against them, `dir` against a two-value whitelist) and every value binds as
// a parameter, so this handler is structurally read-only and immune to
// SQL injection.
// GET /api/accounts/{account}/databases/{db}/tables/{table}?limit=N&offset=M&sort=col&dir=asc|desc
func (s *server) handleDbTableData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("account")
		dbName := r.PathValue("db")
		table := r.PathValue("table")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		sortBy := r.URL.Query().Get("sort")
		sortDir := strings.ToLower(r.URL.Query().Get("dir"))

		if accountName == "" || dbName == "" || table == "" {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "account, db and table are required",
			})
			return
		}

		limit := defaultRowLimit
		if limitStr != "" {
			l, err := strconv.Atoi(limitStr)
			if err != nil || l < 1 {
				sendJSONError(w, http.StatusBadRequest, map[string]string{
					"error": "`limit` must be a positive integer",
				})
				return
			}
			limit = l
		}
		if limit > maxRowLimit {
			limit = maxRowLimit
		}

		offset := 0
		if offsetStr != "" {
			o, err := strconv.Atoi(offsetStr)
			if err != nil || o < 0 {
				sendJSONError(w, http.StatusBadRequest, map[string]string{
					"error": "`offset` must be a non-negative integer",
				})
				return
			}
			offset = o
		}

		conn, err := s.openDBForAccount(accountName, dbName)
		if err != nil {
			sendJSONError(w, http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}
		defer conn.Close()

		ctx := r.Context()

		// Server-trusted column metadata; also doubles as the "table exists
		// in this schema" check (guards typos and cross-database reach).
		cols, err := fetchColumnInfo(conn, ctx, dbName, table)
		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if len(cols) == 0 {
			sendJSONError(w, http.StatusNotFound, map[string]string{
				"error": "Table not found in this database",
			})
			return
		}

		var columnNames []string
		var primaryKey []string
		for _, c := range cols {
			columnNames = append(columnNames, c.Name)
			if c.IsPK {
				primaryKey = append(primaryKey, c.Name)
			}
		}

		// Sort whitelist: the client may only sort by a real column of THIS
		// table (checked against the information_schema metadata above, not
		// trusted from the request) and only in asc/desc. An unknown sort
		// column is an error rather than a silent ignore, so a stale UI can't
		// quietly show unsorted data after a rename.
		if sortBy != "" {
			found := false
			for _, n := range columnNames {
				if n == sortBy {
					found = true
					break
				}
			}
			if !found {
				sendJSONError(w, http.StatusBadRequest, map[string]string{
					"error": "`sort` must name a column of this table",
				})
				return
			}
			if sortDir == "" {
				sortDir = "asc"
			}
			if sortDir != "asc" && sortDir != "desc" {
				sendJSONError(w, http.StatusBadRequest, map[string]string{
					"error": "`dir` must be asc or desc",
				})
				return
			}
		} else {
			sortDir = "asc" // ignored; only PK tie-breakers are used
		}

		// ORDER BY is server-built from a whitelist: `sort` was validated
		// against information_schema above and `dir` against asc/desc, so the
		// only identifiers reaching the SQL text are backtick-quoted,
		// server-sourced names.
		colList := strings.Join(quoteIdentifiers(columnNames), ", ")
		query := fmt.Sprintf("SELECT %s FROM `%s`", colList, table)

		orderBy := []string{}
		if sortBy != "" {
			orderBy = append(orderBy, quoteIdentifier(sortBy)+" "+strings.ToUpper(sortDir))
			// A secondary key on the PK keeps rows with equal sort values
			// (including NULLs, which MySQL groups non-deterministically) in
			// a stable order, so paging never skips or duplicates a row.
			for _, pk := range primaryKey {
				if pk == sortBy {
					continue
				}
				orderBy = append(orderBy, quoteIdentifier(pk)+" ASC")
				break // one tie-breaker is enough
			}
		} else {
			// Default order: the PK when the table has one (InnoDB returns PK
			// order anyway; being explicit makes OFFSET paging coherent).
			for _, pk := range primaryKey {
				orderBy = append(orderBy, quoteIdentifier(pk)+" ASC")
				break
			}
		}
		if len(orderBy) > 0 {
			query += " ORDER BY " + strings.Join(orderBy, ", ")
		}
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

		// Exact total so the UI can page (row counts from information_schema
		// are estimates on InnoDB and unsuitable for computing page counts).
		var totalRows int64
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`", table)
		if err := conn.QueryRowContext(ctx, countQuery).Scan(&totalRows); err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": "Error counting rows: " + err.Error(),
			})
			return
		}

		dataRows, err := conn.QueryContext(ctx, query)
		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": "Error reading table data: " + err.Error(),
			})
			return
		}
		defer dataRows.Close()

		rowCols, err := dataRows.Columns()
		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		// Use []any per row so NULL columns decode to nil (JSON null)
		// instead of erroring on a scan into a fixed type.
		rowsOut := [][]any{}
		for dataRows.Next() {
			values := make([]any, len(rowCols))
			valuePtrs := make([]any, len(rowCols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			if err := dataRows.Scan(valuePtrs...); err != nil {
				sendJSONError(w, http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
				return
			}
			normalizeRowValues(values)
			rowsOut = append(rowsOut, values)
		}
		if err = dataRows.Err(); err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		sendJSON(w, map[string]any{
			"columns":    rowCols,
			"rows":       rowsOut,
			"limit":      limit,
			"offset":     offset,
			"totalRows":  totalRows,
			"columnInfo": cols,
			"primaryKey": primaryKey,
		})
	}
}

// collectFormColumnValues reads `col_<name>` / `pk_<name>` form fields,
// keeping only fields whose suffix matches a real column of the table
// (information_schema is the whitelist). Empty string on a nullable column
// becomes NULL; everything else is passed through verbatim as a query
// parameter value — values NEVER touch the SQL text itself.
func collectFormColumnValues(r *http.Request, byName map[string]columnMeta, prefix string) ([]string, []any) {
	var names []string
	var args []any
	for key, values := range r.PostForm {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		name := strings.TrimPrefix(key, prefix)
		meta, ok := byName[name]
		if !ok {
			continue // silently skip unknown columns
		}
		raw := ""
		if len(values) > 0 {
			raw = values[0]
		}
		if raw == "" && meta.Nullable {
			args = append(args, nil)
		} else {
			args = append(args, raw)
		}
		names = append(names, name)
	}
	return names, args
}

// handleDbInsertRow inserts one row into a table belonging to the account's
// database. Only real columns (per information_schema) are accepted;
// auto-increment columns are ignored. All values are bound as ? parameters,
// so user input never appears in the SQL text.
// POST /api/accounts/{account}/databases/{db}/tables/{table}/rows
//
//	multipart form fields: col_<name>=value...
func (s *server) handleDbInsertRow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("account")
		dbName := r.PathValue("db")
		table := r.PathValue("table")

		if accountName == "" || dbName == "" || table == "" {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "account, db and table are required",
			})
			return
		}

		if err := r.ParseMultipartForm(1 << 20); err != nil {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "invalid form body: " + err.Error(),
			})
			return
		}

		conn, err := s.openDBForAccount(accountName, dbName)
		if err != nil {
			sendJSONError(w, http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}
		defer conn.Close()

		ctx := r.Context()

		cols, err := fetchColumnInfo(conn, ctx, dbName, table)
		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if len(cols) == 0 {
			sendJSONError(w, http.StatusNotFound, map[string]string{
				"error": "Table not found in this database",
			})
			return
		}

		byName := map[string]columnMeta{}
		for _, c := range cols {
			byName[c.Name] = c
		}

		names, args := collectFormColumnValues(r, byName, "col_")
		if len(names) == 0 {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "no column values provided",
			})
			return
		}

		placeholders := make([]string, len(names))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
			table,
			strings.Join(quoteIdentifiers(names), ", "),
			strings.Join(placeholders, ", "))

		res, err := conn.ExecContext(ctx, query, args...)
		if err != nil {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "Error inserting row: " + err.Error(),
			})
			return
		}
		affected, _ := res.RowsAffected()
		insertID, _ := res.LastInsertId()

		s.logger.Info(fmt.Sprintf("db-admin: INSERT %d row(s) into %s.%s (account %s)",
			affected, dbName, table, accountName))

		sendJSON(w, map[string]any{
			"message":  fmt.Sprintf("Inserted %d row(s)", affected),
			"insertId": insertID,
		})
	}
}

// handleDbUpdateRow updates one row, identified by its full primary key.
// Requires the table to have a primary key (no PK => no safe way to target a
// single row). SET columns come only from `col_<name>` form fields and WHERE
// only from `pk_<name>` fields, both validated against information_schema;
// all values bind as ? parameters. PK columns are not editable (they are the
// row identity this handler targets).
// PUT /api/accounts/{account}/databases/{db}/tables/{table}/rows
//
//	multipart form fields: pk_<name>=value..., col_<name>=value...
func (s *server) handleDbUpdateRow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("account")
		dbName := r.PathValue("db")
		table := r.PathValue("table")

		if accountName == "" || dbName == "" || table == "" {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "account, db and table are required",
			})
			return
		}

		if err := r.ParseMultipartForm(1 << 20); err != nil {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "invalid form body: " + err.Error(),
			})
			return
		}

		conn, err := s.openDBForAccount(accountName, dbName)
		if err != nil {
			sendJSONError(w, http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}
		defer conn.Close()

		ctx := r.Context()

		cols, err := fetchColumnInfo(conn, ctx, dbName, table)
		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if len(cols) == 0 {
			sendJSONError(w, http.StatusNotFound, map[string]string{
				"error": "Table not found in this database",
			})
			return
		}

		byName := map[string]columnMeta{}
		var pks []string
		for _, c := range cols {
			byName[c.Name] = c
			if c.IsPK {
				pks = append(pks, c.Name)
			}
		}
		if len(pks) == 0 {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "This table has no primary key, so individual rows can't be safely edited. Add a primary key to enable editing.",
			})
			return
		}

		// SET clause from col_* fields; PK columns are excluded.
		setNames, args := collectFormColumnValues(r, byName, "col_")
		setNames = removeStrings(setNames, pks) // strip any PK the client sent as col_
		args = args[:len(setNames)]             // keep args aligned with names
		if len(setNames) == 0 {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "no column values to update",
			})
			return
		}

		// WHERE clause from pk_* fields; every PK column must be provided.
		whereParts := []string{}
		for _, pk := range pks {
			if _, ok := r.PostForm["pk_"+pk]; !ok {
				sendJSONError(w, http.StatusBadRequest, map[string]string{
					"error": "missing primary key value for column " + pk,
				})
				return
			}
			whereParts = append(whereParts, quoteIdentifier(pk)+" = ?")
			args = append(args, r.PostFormValue("pk_"+pk))
		}

		setParts := make([]string, len(setNames))
		for i, n := range setNames {
			setParts[i] = quoteIdentifier(n) + " = ?"
		}
		query := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s",
			table,
			strings.Join(setParts, ", "),
			strings.Join(whereParts, " AND "))

		res, err := conn.ExecContext(ctx, query, args...)
		if err != nil {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "Error updating row: " + err.Error(),
			})
			return
		}
		affected, _ := res.RowsAffected()
		msg := fmt.Sprintf("Updated %d row(s)", affected)
		if affected == 0 {
			msg = "No rows matched the primary key — the row may have been changed or deleted since it was loaded"
		}

		s.logger.Info(fmt.Sprintf("db-admin: UPDATE %s.%s (account %s): %s",
			dbName, table, accountName, query))

		sendJSON(w, map[string]any{"message": msg})
	}
}

// handleDbDeleteRow deletes exactly one row, identified by its full primary
// key (same identity rule as the update handler). WHERE is built only from
// pk_<name> form fields validated against information_schema; all values bind
// as ? parameters. Tables without a PK are refused — there is no safe way to
// target a single row.
// DELETE /api/accounts/{account}/databases/{db}/tables/{table}/rows
//
//	multipart form fields: pk_<name>=value...
func (s *server) handleDbDeleteRow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountName := r.PathValue("account")
		dbName := r.PathValue("db")
		table := r.PathValue("table")

		if accountName == "" || dbName == "" || table == "" {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "account, db and table are required",
			})
			return
		}

		if err := r.ParseMultipartForm(1 << 20); err != nil {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "invalid form body: " + err.Error(),
			})
			return
		}

		conn, err := s.openDBForAccount(accountName, dbName)
		if err != nil {
			sendJSONError(w, http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}
		defer conn.Close()

		ctx := r.Context()

		cols, err := fetchColumnInfo(conn, ctx, dbName, table)
		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if len(cols) == 0 {
			sendJSONError(w, http.StatusNotFound, map[string]string{
				"error": "Table not found in this database",
			})
			return
		}

		pks := []string{}
		for _, c := range cols {
			if c.IsPK {
				pks = append(pks, c.Name)
			}
		}
		if len(pks) == 0 {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "This table has no primary key, so individual rows can't be safely deleted. Add a primary key to enable deleting.",
			})
			return
		}

		whereParts := []string{}
		args := []any{}
		for _, pk := range pks {
			if _, ok := r.PostForm["pk_"+pk]; !ok {
				sendJSONError(w, http.StatusBadRequest, map[string]string{
					"error": "missing primary key value for column " + pk,
				})
				return
			}
			whereParts = append(whereParts, quoteIdentifier(pk)+" = ?")
			args = append(args, r.PostFormValue("pk_"+pk))
		}

		query := fmt.Sprintf("DELETE FROM `%s` WHERE %s",
			table,
			strings.Join(whereParts, " AND "))

		res, err := conn.ExecContext(ctx, query, args...)
		if err != nil {
			sendJSONError(w, http.StatusBadRequest, map[string]string{
				"error": "Error deleting row: " + err.Error(),
			})
			return
		}
		affected, _ := res.RowsAffected()
		msg := fmt.Sprintf("Deleted %d row(s)", affected)
		if affected == 0 {
			msg = "No rows matched the primary key — the row may have been deleted already"
		}

		s.logger.Info(fmt.Sprintf("db-admin: DELETE %s.%s (account %s): %s",
			dbName, table, accountName, query))

		sendJSON(w, map[string]any{"message": msg})
	}
}

// quoteIdentifier wraps a single server-sourced identifier in backticks.
func quoteIdentifier(id string) string {
	return "`" + strings.ReplaceAll(id, "`", "``") + "`"
}

// quoteIdentifiers wraps each identifier in backticks so SELECT column lists
// are safe against reserved words (e.g. a column named `order` or `key`).
// The identifiers are server-sourced (information_schema), never client input.
func quoteIdentifiers(idents []string) []string {
	out := make([]string, len(idents))
	for i, id := range idents {
		out[i] = quoteIdentifier(id)
	}
	return out
}

// removeStrings returns a copy of base with all elements found in drop removed.
func removeStrings(base, drop []string) []string {
	dropSet := map[string]bool{}
	for _, d := range drop {
		dropSet[d] = true
	}
	out := []string{}
	for _, b := range base {
		if !dropSet[b] {
			out = append(out, b)
		}
	}
	return out
}
