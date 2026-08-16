package dsqlloader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

type QueryKind string

const (
	KindOne   QueryKind = "one"
	KindMany  QueryKind = "many"
	KindCount QueryKind = "count"
	KindExec  QueryKind = "exec"
)

type Query struct {
	SQL       string    // SQL is the raw query text with the header line removed.
	Name      string    // Name is the PascalCase query name parsed from the header comment.
	Kind      QueryKind // Kind is the return kind parsed from the header comment (one, many, count, exec).
	File      string    // File is the path to the file of the query.
	StartLine int       // StartLine is the line number where the query starts.
}

type Domain struct {
	Name    string  // DomainName is the PascalCase name derived from the .dsql file name.
	Queries []Query // Queries holds every query found in the file.
}

// Load walks queriesDir, parses every .dsql file found, and returns
// one Domain per file with its queries extracted and validated.
// Files or query blocks that don't match the expected format are
// skipped with warning.
func Load(queriesDir string) ([]Domain, error) {
	var domains []Domain

	entries, err := os.ReadDir(queriesDir)
	if err != nil {
		return domains, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		fullpath := filepath.Join(".", queriesDir, filename)
		extname := filepath.Ext(filename)

		if extname != ".dsql" {
			continue
		}

		var domain Domain

		domainNameRaw := strings.TrimSuffix(filename, extname)
		domain.Name = strcase.ToCamel(domainNameRaw)

		fileContent, err := os.ReadFile(fullpath)
		if err != nil {
			return domains, fmt.Errorf("reading %s: %w", filename, err)
		}

		queries := parseQueries(fullpath, string(fileContent))
		domain.Queries = queries

		domains = append(domains, domain)
	}

	return domains, nil
}

type filectx struct {
	path string
	line int
}

func (c *filectx) nextline() {
	c.line++
}

func (c filectx) current() string {
	return fmt.Sprintf("%s:%d", c.path, c.line)
}

// parseQueries parses queries and their metadata from the provided file content.
func parseQueries(filepath string, content string) []Query {
	var queries []Query

	ctx := filectx{filepath, 0}
	// Mode is a flag for parsing modes.
	// 0 - waiting for the query to start.
	// 1 - parsing query.
	// 2 - skipping query without annotation.
	mode := 0

	queryNames := map[string]filectx{}

	var prevMetadataAt filectx
	hasMetadata := false

	var queryStartAt filectx
	queryName := ""
	queryKind := QueryKind("")

	queryBuilder := strings.Builder{}

	flushQuery := func() {
		queries = append(queries, Query{
			SQL:       strings.TrimSpace(queryBuilder.String()),
			Name:      queryName,
			Kind:      queryKind,
			File:      filepath,
			StartLine: queryStartAt.line,
		})

		queryBuilder.Reset()
		mode = 0
		hasMetadata = false
	}

	startMetadata := func(line string) {
		mode = 0
		hasMetadata = true
		prevMetadataAt = ctx

		queryName, queryKind = parseMetadata(ctx, line)

		prevQueryAt, ok := queryNames[queryName]
		if ok {
			log.Fatalf("duplicated query name %s at %s and %s\n", queryName, prevQueryAt.current(), ctx.current())
		}
		queryNames[queryName] = ctx
	}

	for rawLine := range strings.SplitSeq(content, "\n") {
		ctx.nextline()
		line := strings.TrimSpace(rawLine)

		if mode == 0 {
			// Skip empty lines while waiting for the query
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "-->") {
				if hasMetadata {
					// Warning if we've met another metadata line before.
					fmt.Printf("Warning: metadata without query at %s\n", prevMetadataAt.current())
				}

				startMetadata(line)

				continue
			}

			// Skip comments while waiting for the query
			if strings.HasPrefix(line, "--") {
				continue
			}

			// If there is no metadata, enter the abandoned query mode.
			if !hasMetadata {
				log.Printf("Warning: query without metadata at %s\n", ctx.current())
				mode = 2

				continue
			}

			// Join current line to the query and enter into the parse query mode.
			mode = 1
			queryStartAt = ctx

			queryBuilder.WriteString(line)
			queryBuilder.WriteString("\n")
		} else if mode == 1 {
			if strings.HasPrefix(line, "-->") {
				flushQuery()

				startMetadata(line)

				continue
			}

			queryBuilder.WriteString(line)
			queryBuilder.WriteString("\n")
		} else if mode == 2 {
			// Skip lines in abandoned query mode until metadata line.

			if strings.HasPrefix(line, "-->") {
				startMetadata(line)

				continue
			}
		}
	}

	if mode == 1 {
		flushQuery()
	} else if hasMetadata {
		fmt.Printf("Warning: metadata without query at %s\n", prevMetadataAt.current())
	}

	return queries
}

func parseMetadata(ctx filectx, line string) (string, QueryKind) {
	parts := strings.Fields(line)

	if len(parts) != 3 {
		log.Fatalf("invalid metadata line at %s\n", ctx.current())
		return "", QueryKind("")
	}

	return strcase.ToCamel(parts[1]), parseQueryKind(ctx, parts[2])
}

func parseQueryKind(ctx filectx, raw string) QueryKind {
	switch raw {
	case "one", "many", "count", "exec":
		return QueryKind(raw)
	default:
		log.Fatalf("invalid query kind \"%s\" at %s\n", raw, ctx.current())
		return QueryKind("")
	}
}
