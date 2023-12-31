package am

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/adriancrafter/todoapp/internal/am/errors"
)

var (
	cfgKey = Key
)

const (
	GetAllKey     = "GetAll"
	GetOneKey     = "GetOne"
	CreateKey     = "CreateErr"
	UpdateKey     = "UpdateErr"
	SoftDeleteKey = "SoftDelete"
	DeleteKey     = "DeleteErr"
	PurgeKey      = "Purge"
)

const (
	queriesDir = "assets/queries/%s"
)

var (
	validSQLPrefixes = []string{
		"SELECT",
		"INSERT",
		"UPDATE",
		"DELETE",
		"CREATE",
		"ALTER",
		"DROP",
		"TRUNCATE",
		"USE",
		"SET",
		"BEGIN",
		"COMMIT",
		"ROLLBACK",
		"DECLARE",
		"FETCH",
		"OPEN",
		"CLOSE",
		"RETURN",
		"GRANT",
		"REVOKE",
	}

	validSQLPrefixesMap = make(map[string]struct{})
)

type (
	QueryManager struct {
		*SimpleCore
		dir     string
		queryFS embed.FS
		queries map[string]Queries
	}

	Queries struct {
		Model string
		List  map[string]string
	}
)

func NewQueryManager(queryFS embed.FS, dir string, opts ...Option) *QueryManager {
	return &QueryManager{
		SimpleCore: NewCore("query-manager", opts...),
		dir:        dir,
		queryFS:    queryFS,
		queries:    make(map[string]Queries),
	}
}

func (qm *QueryManager) Setup(ctx context.Context) {
	loadValidPrefixes()

	err := qm.LoadQueries()
	if err != nil {
		qm.Log().Errorf("error loading queries", err)
	}

	qm.dumpQueries()
}

func (qm *QueryManager) LoadQueries() error {
	modelNameRegex := regexp.MustCompile(`.*\/(\w+)\.sql`)

	path := fmt.Sprintf(queriesDir, qm.dir)
	files, err := fs.ReadDir(qm.queryFS, path)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		fileName := fileInfo.Name()
		modelNameMatches := modelNameRegex.FindStringSubmatch(fileName)
		modelName := ""

		if len(modelNameMatches) == 2 {
			modelName = modelNameMatches[1]
			modelName = toCamelCase(modelName)
		}

		sqlContent, err := qm.queryFS.ReadFile(fmt.Sprintf("assets/queries/pg/%s", fileName))
		if err != nil {
			return errors.Wrap(err, "error reading %s", fileName)
		}

		queries := make(map[string]string)
		queryParts := strings.Split(string(sqlContent), "--")
		var querySection string
		var queryName string
		for _, part := range queryParts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "Model:") {
				modelName = strings.TrimSpace(strings.TrimPrefix(part, "Model:"))

			} else {
				querySection = strings.TrimSpace(strings.TrimPrefix(part, "--"))
				parts := strings.SplitN(querySection, "\n", 2)
				if len(parts) == 2 && isValidPrefix(parts[1]) {
					queryName = strings.TrimSpace(parts[0])
					queries[queryName] = strings.TrimSpace(parts[1])
				}
			}
		}

		modelQueries := Queries{
			Model: modelName,
			List:  queries,
		}
		qm.queries[modelName] = modelQueries
	}

	return nil
}

func (qm *QueryManager) Get(model, queryName string) (query string, err error) {
	modelQueries, ok := qm.queries[string(model)]
	if !ok {
		return query, errors.NewErrorf("query '%s' for model '%s' not found", queryName, model)
	}

	query, ok = modelQueries.List[string(queryName)]
	if !ok {
		return query, errors.NewErrorf("query '%s' for model '%s' not found", queryName, model)
	}

	return query, nil
}

func (qm *QueryManager) dumpQueries() {
	compact := qm.Cfg().GetBool(cfgKey.LogOutputCompact)
	multiSpaces := regexp.MustCompile(`\s+`)
	qm.Log().Debug("----------------")
	for model, queries := range qm.queries {
		for name, query := range queries.List {
			if compact {
				query = strings.ReplaceAll(query, "\n", " ")
				query = strings.ReplaceAll(query, "\t", " ")
				query = multiSpaces.ReplaceAllString(query, " ")
			}
			qm.Log().Debugf("[%s] %s: %s", model, name, query)
		}
		qm.Log().Debug("================")
	}
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func loadValidPrefixes() {
	for _, prefix := range validSQLPrefixes {
		validSQLPrefixesMap[prefix] = struct{}{}
	}
}

func isValidPrefix(query string) bool {
	words := strings.Fields(query)
	if len(words) == 0 {
		return false
	}

	firstWord := strings.ToUpper(words[0])
	_, ok := validSQLPrefixesMap[firstWord]
	return ok
}
