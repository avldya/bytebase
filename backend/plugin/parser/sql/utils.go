package parser

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	tidbparser "github.com/pingcap/tidb/parser"
	tidbast "github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/plugin/parser/sql/ast"
)

// SingleSQL is a separate SQL split from multi-SQL.
type SingleSQL struct {
	Text     string
	LastLine int
	// The sql is empty, such as `/* comments */;` or just `;`.
	Empty bool
}

// GetSQLFingerprint returns the fingerprint of the SQL.
func GetSQLFingerprint(engineType EngineType, sql string) (string, error) {
	switch engineType {
	case MySQL, TiDB, MariaDB:
		return getMySQLFingerprint(sql)
	default:
		return "", errors.Errorf("engine type is not supported: %s", engineType)
	}
}

// From https://github.com/percona/percona-toolkit/blob/af686fe186d1fca4c4392c8fa75c31a00c8fb273/bin/pt-query-digest#L2930
func getMySQLFingerprint(query string) (string, error) {
	// Match SQL queries generated by mysqldump command.
	if matched, _ := regexp.MatchString(`\ASELECT /\*!40001 SQL_NO_CACHE \*/ \* FROM `, query); matched {
		return "mysqldump", nil
	}
	// Match SQL queries generated by Percona Toolkit.
	if matched, _ := regexp.MatchString(`/\*\w+\.\w+:[0-9]/[0-9]\*/`, query); matched {
		return "percona-toolkit", nil
	}
	// Match administrator commands.
	if matched, _ := regexp.MatchString(`\Aadministrator command: `, query); matched {
		return query, nil
	}
	// Match stored procedure call statements.
	if matched, _ := regexp.MatchString(`\A\s*(call\s+\S+)\(`, query); matched {
		return strings.ToLower(regexp.MustCompile(`\A\s*(call\s+\S+)\(`).FindStringSubmatch(query)[1]), nil
	}
	// Match INSERT INTO or REPLACE INTO statements.
	if beginning := regexp.MustCompile(`(?i)((?:INSERT|REPLACE)(?: IGNORE)?\s+INTO.+?VALUES\s*\(.*?\))\s*,\s*\(`).FindStringSubmatch(query); len(beginning) > 0 {
		query = beginning[1]
	}

	// Match multi-line comments and single-line comments, and remove them.
	mlcRe := regexp.MustCompile(`(?s)/\*.*?\*/`)
	olcRe := regexp.MustCompile(`--.*$`)
	query = mlcRe.ReplaceAllString(query, "")
	query = olcRe.ReplaceAllString(query, "")

	// Replace the database name in USE statements with a question mark (?).
	query = regexp.MustCompile(`(?i)\Ause \S+\z`).ReplaceAllString(query, "use ?")

	// Replace escape characters and special characters in SQL queries with a question mark (?).
	query = regexp.MustCompile(`([^\\])(\\')`).ReplaceAllString(query, "$1")
	query = regexp.MustCompile(`([^\\])(\\")`).ReplaceAllString(query, "$1")
	query = regexp.MustCompile(`\\\\`).ReplaceAllString(query, "")
	query = regexp.MustCompile(`\\'`).ReplaceAllString(query, "")
	query = regexp.MustCompile(`\\"`).ReplaceAllString(query, "")
	query = regexp.MustCompile(`([^\\])(".*?[^\\]?")`).ReplaceAllString(query, "$1?")
	query = regexp.MustCompile(`([^\\])('.*?[^\\]?')`).ReplaceAllString(query, "$1?")

	// Replace boolean values in SQL queries with a question mark (?).
	query = regexp.MustCompile(`\bfalse\b|\btrue\b`).ReplaceAllString(query, "?")

	// Replace MD5 values in SQL queries with a question mark (?).
	if matched, _ := regexp.MatchString(`([._-])[a-f0-9]{32}`, query); matched {
		query = regexp.MustCompile(`([._-])[a-f0-9]{32}`).ReplaceAllString(query, "$1?")
	}

	// Replace numbers in SQL queries with a question mark (?).
	if matched, _ := regexp.MatchString(`\b[0-9+-][0-9a-f.xb+-]*`, query); matched {
		query = regexp.MustCompile(`\b[0-9+-][0-9a-f.xb+-]*`).ReplaceAllString(query, "?")
	}

	// Replace special characters in SQL queries with a question mark (?).
	if matched, _ := regexp.MatchString(`[xb+-]\?`, query); matched {
		query = regexp.MustCompile(`[xb+-]\?`).ReplaceAllString(query, "?")
	} else {
		query = regexp.MustCompile(`[xb.+-]\?`).ReplaceAllString(query, "?")
	}

	// Remove spaces and line breaks in SQL queries.
	query = strings.TrimSpace(query)
	query = strings.TrimRight(query, "\n\r\f ")
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")
	query = strings.ToLower(query)

	// Replace NULL values in SQL queries with a question mark (?).
	query = regexp.MustCompile(`\bnull\b`).ReplaceAllString(query, "?")

	// Replace IN and VALUES clauses in SQL queries with a question mark (?).
	query = regexp.MustCompile(`\b(in|values?)(?:[\s,]*\([\s?,]*\))+`).ReplaceAllString(query, "$1(?+)")

	var err error
	query, err = collapseUnion(query)
	if err != nil {
		return "", err
	}

	// Replace numbers in the LIMIT clause of SQL queries with a question mark (?).
	query = regexp.MustCompile(`\blimit \?(?:, ?\?| offset \?)?`).ReplaceAllString(query, "limit ?")

	// Remove ASC sorting in SQL queries.
	if matched, _ := regexp.MatchString(`\border by `, query); matched {
		ascRegexp := regexp.MustCompile(`(.+?)\s+asc`)
		for {
			if matched := ascRegexp.MatchString(query); matched {
				query = ascRegexp.ReplaceAllString(query, "$1")
			} else {
				break
			}
		}
	}

	return query, nil
}

func collapseUnion(query string) (string, error) {
	// The origin perl code is:
	//   $query =~ s{                          # Collapse UNION
	//     \b(select\s.*?)(?:(\sunion(?:\sall)?)\s\1)+
	//	  }
	//	  {$1 /*repeat$2*/}xg;
	// But Golang doesn't support \1(back-reference).
	// So we use the following code to replace it.
	unionRegexp := regexp.MustCompile(`\s(union all|union)\s`)
	parts := unionRegexp.Split(query, -1)
	if len(parts) == 1 {
		return query, nil
	}
	// Add a sentinel node to the end of the slice.
	// Because we remove all comments before, so all parts are different from sentinel node.
	parts = append(parts, "/*Sentinel Node*/")
	separators := unionRegexp.FindAllString(query, -1)
	if len(parts) != len(separators)+2 {
		return "", errors.Errorf("find %d parts, but %d separators", len(parts)-1, len(separators))
	}
	start := 0
	var buf bytes.Buffer
	if _, err := buf.WriteString(parts[start]); err != nil {
		return "", err
	}
	for i, part := range parts {
		if i == 0 {
			continue
		}
		if part == parts[start] {
			continue
		}
		if i == start+1 {
			// The i-th part is not equal to the front part.
			if _, err := buf.WriteString(separators[i-1]); err != nil {
				return "", err
			}
		} else {
			// deal with the same parts[start, i-1] and start < i-1.
			if _, err := buf.WriteString(" /*repeat"); err != nil {
				return "", err
			}
			// Write the last separator between the same parts[start, i-1].
			// In other words, the last separator is the separator between the i-th part and the (i-1)-th part.
			// So the index of the last separator is (i-1)-1.
			if _, err := buf.WriteString(separators[(i-1)-1]); err != nil {
				return "", err
			}
			if _, err := buf.WriteString("*/"); err != nil {
				return "", err
			}
		}
		start = i
		// Don't write the sentinel node.
		if i != len(parts)-1 {
			if _, err := buf.WriteString(parts[start]); err != nil {
				return "", err
			}
		}
	}
	return buf.String(), nil
}

// SplitMultiSQLAndNormalize split multiple SQLs and normalize them.
// For MySQL, filter DELIMITER statements and replace all non-semicolon delimiters with semicolons.
func SplitMultiSQLAndNormalize(engineType EngineType, statement string) ([]SingleSQL, error) {
	switch engineType {
	case MySQL:
		list, err := SplitMultiSQL(MySQL, statement)
		if err != nil {
			return nil, err
		}

		var result []SingleSQL
		delimiter := `;`
		for _, sql := range list {
			if IsDelimiter(sql.Text) {
				delimiter, err = ExtractDelimiter(sql.Text)
				if err != nil {
					return nil, err
				}
				continue
			}
			if delimiter != ";" {
				result = append(result, SingleSQL{
					Text:     fmt.Sprintf("%s;", strings.TrimSuffix(sql.Text, delimiter)),
					LastLine: sql.LastLine,
					Empty:    sql.Empty,
				})
			} else {
				result = append(result, sql)
			}
		}
		return result, nil
	default:
		return SplitMultiSQL(engineType, statement)
	}
}

// SplitMultiSQL splits statement into a slice of the single SQL.
func SplitMultiSQL(engineType EngineType, statement string) ([]SingleSQL, error) {
	var list []SingleSQL
	var err error
	switch engineType {
	case Oracle, MSSQL:
		t := newTokenizer(statement)
		list, err = t.splitStandardMultiSQL()
	case Postgres, Redshift:
		t := newTokenizer(statement)
		list, err = t.splitPostgreSQLMultiSQL()
	case MySQL, TiDB, MariaDB, OceanBase:
		t := newTokenizer(statement)
		list, err = t.splitMySQLMultiSQL()
	default:
		return nil, errors.Errorf("engine type is not supported: %s", engineType)
	}

	if err != nil {
		return nil, err
	}

	var result []SingleSQL
	for _, sql := range list {
		if sql.Empty {
			continue
		}
		if engineType == Oracle {
			sql.Text = strings.TrimRight(sql.Text, " \n\t;")
		}
		result = append(result, sql)
	}
	return result, nil
}

// SplitMultiSQLStream splits statement stream into a slice of the single SQL.
func SplitMultiSQLStream(engineType EngineType, src io.Reader, f func(string) error) ([]SingleSQL, error) {
	var list []SingleSQL
	var err error
	switch engineType {
	case Oracle, MSSQL:
		t := newStreamTokenizer(src, f)
		list, err = t.splitStandardMultiSQL()
	case Postgres, Redshift:
		t := newStreamTokenizer(src, f)
		list, err = t.splitPostgreSQLMultiSQL()
	case MySQL, TiDB, MariaDB, OceanBase:
		t := newStreamTokenizer(src, f)
		list, err = t.splitMySQLMultiSQL()
	default:
		return nil, errors.Errorf("engine type is not supported: %s", engineType)
	}

	if err != nil {
		return nil, err
	}

	var result []SingleSQL
	for _, sql := range list {
		if sql.Empty {
			continue
		}
		if engineType == Oracle {
			sql.Text = strings.TrimRight(sql.Text, " \n\t;")
		}
		result = append(result, sql)
	}

	return result, nil
}

// SetLineForCreateTableStmt sets the line for columns and table constraints in CREATE TABLE statements.
func SetLineForCreateTableStmt(engineType EngineType, node *ast.CreateTableStmt) error {
	switch engineType {
	case Postgres:
		t := newTokenizer(node.Text())
		firstLine := node.LastLine() - strings.Count(node.Text(), "\n")
		return t.setLineForPGCreateTableStmt(node, firstLine)
	default:
		return errors.Errorf("engine type is not supported: %s", engineType)
	}
}

// SetLineForMySQLCreateTableStmt sets the line for columns and table constraints in MySQL CREATE TABLE statments.
// This is a temporary function. Because we do not convert tidb AST to our AST. So we have to implement this.
// TODO(rebelice): remove it.
func SetLineForMySQLCreateTableStmt(node *tidbast.CreateTableStmt) error {
	// exclude CREATE TABLE ... AS and CREATE TABLE ... LIKE statement.
	if len(node.Cols) == 0 {
		return nil
	}
	firstLine := node.OriginTextPosition() - strings.Count(node.Text(), "\n")
	return newTokenizer(node.Text()).setLineForMySQLCreateTableStmt(node, firstLine)
}

// ExtractTiDBUnsupportStmts returns a list of unsupported statements in TiDB extracted from the `stmts`,
// and returns the remaining statements supported by TiDB from `stmts`.
func ExtractTiDBUnsupportStmts(stmts string) ([]string, string, error) {
	var unsupportStmts []string
	var supportedStmts bytes.Buffer
	// We use our bb tokenizer to help us split the multi-statements into statement list.
	singleSQLs, err := SplitMultiSQL(MySQL, stmts)
	if err != nil {
		return nil, "", errors.Wrapf(err, "cannot split multi sql %q via bytebase parser", stmts)
	}
	for _, sql := range singleSQLs {
		content := sql.Text
		if isTiDBUnsupportStmt(content) {
			unsupportStmts = append(unsupportStmts, content)
			continue
		}
		_, _ = supportedStmts.Write([]byte(content))
		_, _ = supportedStmts.Write([]byte("\n"))
	}
	return unsupportStmts, supportedStmts.String(), nil
}

// isTiDBUnsupportStmt returns true if this statement is unsupported in TiDB.
func isTiDBUnsupportStmt(stmt string) bool {
	if IsTiDBUnsupportDDLStmt(stmt) {
		return true
	}
	// Match DELIMITER statement
	// Now, we assume that all input comes from our mysqldump, and the tokenizer can split the mysqldump DELIMITER statement
	// in one singleSQL correctly, so we can handle it easily here by checking the prefix.
	return IsDelimiter(stmt)
}

// IsTiDBUnsupportDDLStmt checks whether the `stmt` is unsupported DDL statement in TiDB, the following statements are unsupported:
// 1. `CREATE TRIGGER`
// 2. `CREATE EVENT`
// 3. `CREATE FUNCTION`
// 4. `CREATE PROCEDURE`.
func IsTiDBUnsupportDDLStmt(stmt string) bool {
	objects := []string{
		"TRIGGER",
		"EVENT",
		"FUNCTION",
		"PROCEDURE",
	}
	createRegexFmt := "(?i)^\\s*CREATE\\s+(DEFINER=(`(.)+`|(.)+)@(`(.)+`|(.)+)(\\s)+)?%s\\s+"
	dropRegexFmt := "(?i)^\\s*DROP\\s+%s\\s+"
	for _, obj := range objects {
		createRegexp := fmt.Sprintf(createRegexFmt, obj)
		re := regexp.MustCompile(createRegexp)
		if re.MatchString(stmt) {
			return true
		}
		dropRegexp := fmt.Sprintf(dropRegexFmt, obj)
		re = regexp.MustCompile(dropRegexp)
		if re.MatchString(stmt) {
			return true
		}
	}
	return false
}

// IsDelimiter returns true if the statement is a delimiter statement.
func IsDelimiter(stmt string) bool {
	delimiterRegex := `(?i)^\s*DELIMITER\s+`
	re := regexp.MustCompile(delimiterRegex)
	return re.MatchString(stmt)
}

// ExtractDelimiter extracts the delimiter from the delimiter statement.
func ExtractDelimiter(stmt string) (string, error) {
	delimiterRegex := `(?i)^\s*DELIMITER\s+(?P<DELIMITER>[^\s\\]+)\s*`
	re := regexp.MustCompile(delimiterRegex)
	matchList := re.FindStringSubmatch(stmt)
	index := re.SubexpIndex("DELIMITER")
	if index >= 0 && index < len(matchList) {
		return matchList[index], nil
	}
	return "", errors.Errorf("cannot extract delimiter from %q", stmt)
}

// TypeString returns the string representation of the type for MySQL.
func TypeString(tp byte) string {
	switch tp {
	case mysql.TypeTiny:
		return "tinyint"
	case mysql.TypeShort:
		return "smallint"
	case mysql.TypeInt24:
		return "mediumint"
	case mysql.TypeLong:
		return "int"
	case mysql.TypeLonglong:
		return "bigint"
	case mysql.TypeFloat:
		return "float"
	case mysql.TypeDouble:
		return "double"
	case mysql.TypeNewDecimal:
		return "decimal"
	case mysql.TypeVarchar:
		return "varchar"
	case mysql.TypeBit:
		return "bit"
	case mysql.TypeTimestamp:
		return "timestamp"
	case mysql.TypeDatetime:
		return "datetime"
	case mysql.TypeDate:
		return "date"
	case mysql.TypeDuration:
		return "time"
	case mysql.TypeJSON:
		return "json"
	case mysql.TypeEnum:
		return "enum"
	case mysql.TypeSet:
		return "set"
	case mysql.TypeTinyBlob:
		return "tinyblob"
	case mysql.TypeMediumBlob:
		return "mediumblob"
	case mysql.TypeLongBlob:
		return "longblob"
	case mysql.TypeBlob:
		return "blob"
	case mysql.TypeVarString:
		return "varbinary"
	case mysql.TypeString:
		return "binary"
	case mysql.TypeGeometry:
		return "geometry"
	}
	return "unknown"
}

// ExtractDatabaseList extracts all databases from statement.
func ExtractDatabaseList(engineType EngineType, statement string) ([]string, error) {
	switch engineType {
	case MySQL, TiDB, MariaDB, OceanBase:
		return extractMySQLDatabaseList(statement)
	default:
		return nil, errors.Errorf("engine type is not supported: %s", engineType)
	}
}

func newMySQLParser() *tidbparser.Parser {
	p := tidbparser.New()

	// To support MySQL8 window function syntax.
	// See https://github.com/bytebase/bytebase/issues/175.
	p.EnableWindowFunc(true)

	return p
}

func extractMySQLDatabaseList(statement string) ([]string, error) {
	databaseMap := make(map[string]bool)

	p := newMySQLParser()
	nodeList, _, err := p.Parse(statement, "", "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parser statement %q", statement)
	}

	for _, node := range nodeList {
		databaseList := extractDatabaseListFromNode(node)
		for _, database := range databaseList {
			databaseMap[database] = true
		}
	}

	var databaseList []string
	for database := range databaseMap {
		databaseList = append(databaseList, database)
	}
	sort.Slice(databaseList, func(i, j int) bool {
		return databaseList[i] < databaseList[j]
	})
	return databaseList, nil
}

// extractDatabaseListFromNode extracts all the database from node.
func extractDatabaseListFromNode(in tidbast.Node) []string {
	tableNameList := ExtractMySQLTableList(in, false /* asName */)

	databaseMap := make(map[string]bool)
	for _, tableName := range tableNameList {
		databaseMap[tableName.Schema.O] = true
	}

	var databaseList []string
	for databaseName := range databaseMap {
		databaseList = append(databaseList, databaseName)
	}

	sort.Strings(databaseList)
	return databaseList
}

// ExtractMySQLTableList extracts all the TableNames from node.
// If asName is true, extract AsName prior to OrigName.
func ExtractMySQLTableList(in tidbast.Node, asName bool) []*tidbast.TableName {
	input := []*tidbast.TableName{}
	return extractTableList(in, input, asName)
}

// -------------------------------------------- DO NOT TOUCH --------------------------------------------

// extractTableList extracts all the TableNames from node.
// If asName is true, extract AsName prior to OrigName.
// Privilege check should use OrigName, while expression may use AsName.
// WARNING: copy from TiDB core code, do NOT touch!
func extractTableList(node tidbast.Node, input []*tidbast.TableName, asName bool) []*tidbast.TableName {
	switch x := node.(type) {
	case *tidbast.SelectStmt:
		if x.From != nil {
			input = extractTableList(x.From.TableRefs, input, asName)
		}
		if x.Where != nil {
			input = extractTableList(x.Where, input, asName)
		}
		if x.With != nil {
			for _, cte := range x.With.CTEs {
				input = extractTableList(cte.Query, input, asName)
			}
		}
		for _, f := range x.Fields.Fields {
			if s, ok := f.Expr.(*tidbast.SubqueryExpr); ok {
				input = extractTableList(s, input, asName)
			}
		}
	case *tidbast.DeleteStmt:
		input = extractTableList(x.TableRefs.TableRefs, input, asName)
		if x.IsMultiTable {
			for _, t := range x.Tables.Tables {
				input = extractTableList(t, input, asName)
			}
		}
		if x.Where != nil {
			input = extractTableList(x.Where, input, asName)
		}
		if x.With != nil {
			for _, cte := range x.With.CTEs {
				input = extractTableList(cte.Query, input, asName)
			}
		}
	case *tidbast.UpdateStmt:
		input = extractTableList(x.TableRefs.TableRefs, input, asName)
		for _, e := range x.List {
			input = extractTableList(e.Expr, input, asName)
		}
		if x.Where != nil {
			input = extractTableList(x.Where, input, asName)
		}
		if x.With != nil {
			for _, cte := range x.With.CTEs {
				input = extractTableList(cte.Query, input, asName)
			}
		}
	case *tidbast.InsertStmt:
		input = extractTableList(x.Table.TableRefs, input, asName)
		input = extractTableList(x.Select, input, asName)
	case *tidbast.SetOprStmt:
		l := &tidbast.SetOprSelectList{}
		unfoldSelectList(x.SelectList, l)
		for _, s := range l.Selects {
			input = extractTableList(s.(tidbast.ResultSetNode), input, asName)
		}
	case *tidbast.PatternInExpr:
		if s, ok := x.Sel.(*tidbast.SubqueryExpr); ok {
			input = extractTableList(s, input, asName)
		}
	case *tidbast.ExistsSubqueryExpr:
		if s, ok := x.Sel.(*tidbast.SubqueryExpr); ok {
			input = extractTableList(s, input, asName)
		}
	case *tidbast.BinaryOperationExpr:
		if s, ok := x.R.(*tidbast.SubqueryExpr); ok {
			input = extractTableList(s, input, asName)
		}
	case *tidbast.SubqueryExpr:
		input = extractTableList(x.Query, input, asName)
	case *tidbast.Join:
		input = extractTableList(x.Left, input, asName)
		input = extractTableList(x.Right, input, asName)
	case *tidbast.TableSource:
		if s, ok := x.Source.(*tidbast.TableName); ok {
			if x.AsName.L != "" && asName {
				newTableName := *s
				newTableName.Name = x.AsName
				newTableName.Schema = model.NewCIStr("")
				input = append(input, &newTableName)
			} else {
				input = append(input, s)
			}
		} else if s, ok := x.Source.(*tidbast.SelectStmt); ok {
			if s.From != nil {
				var innerList []*tidbast.TableName
				innerList = extractTableList(s.From.TableRefs, innerList, asName)
				if len(innerList) > 0 {
					innerTableName := innerList[0]
					if x.AsName.L != "" && asName {
						newTableName := *innerList[0]
						newTableName.Name = x.AsName
						newTableName.Schema = model.NewCIStr("")
						innerTableName = &newTableName
					}
					input = append(input, innerTableName)
				}
			}
		}
	}
	return input
}

// WARNING: copy from TiDB core code, do NOT touch!
func unfoldSelectList(list *tidbast.SetOprSelectList, unfoldList *tidbast.SetOprSelectList) {
	for _, sel := range list.Selects {
		switch s := sel.(type) {
		case *tidbast.SelectStmt:
			unfoldList.Selects = append(unfoldList.Selects, s)
		case *tidbast.SetOprSelectList:
			unfoldSelectList(s, unfoldList)
		}
	}
}
