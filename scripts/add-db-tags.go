package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <go-file>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	// Validate the filename to prevent directory traversal attacks
	if !strings.HasSuffix(filename, ".go") {
		fmt.Fprintf(os.Stderr, "Error: file must have .go extension\n")
		os.Exit(1)
	}

	// Clean the path to prevent directory traversal
	filename = filepath.Clean(filename)

	// Check if the file exists and is a regular file
	if info, err := os.Stat(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error: file does not exist: %v\n", err)
		os.Exit(1)
	} else if !info.Mode().IsRegular() {
		fmt.Fprintf(os.Stderr, "Error: not a regular file\n")
		os.Exit(1)
	}

	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// Add db tags to all structs
	modified := false
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				if addDBTagsToStruct(typeSpec.Name.Name, structType) {
					modified = true
				}
			}
		}
		return true
	})

	if !modified {
		fmt.Println("No modifications needed")
		return
	}

	// Write the modified file back
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %v\n", closeErr)
		}
	}()

	if err := format.Node(f, fset, node); err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting file: %v\n", err)
		_ = f.Close()
		return
	}

	fmt.Printf("Successfully added db tags to %s\n", filename)
}

// addDBTagsToStruct adds db tags to all fields in a struct
func addDBTagsToStruct(structName string, structType *ast.StructType) bool {
	modified := false

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		fieldName := field.Names[0].Name

		// Skip if field already has db tag
		if hasDBTag(field) {
			continue
		}

		// Generate db tag
		dbTagValue := generateDBTag(fieldName)

		// Add the db tag
		addTagToField(field, "db", dbTagValue)

		// Special handling for ID field - set json tag to "-"
		if fieldName == "Id" {
			updateJSONTag(field, "-")
		}

		modified = true
	}

	return modified
}

// hasDBTag checks if a field already has a db tag
func hasDBTag(field *ast.Field) bool {
	if field.Tag == nil {
		return false
	}

	tagValue := strings.Trim(field.Tag.Value, "`")
	return strings.Contains(tagValue, "db:")
}

// generateDBTag converts a field name to snake_case for the db tag
func generateDBTag(fieldName string) string {
	// Special cases
	switch fieldName {
	case "Id", "ID":
		return "id"
	case "URL":
		return "url"
	case "IP":
		return "ip"
	}

	return toSnakeCase(fieldName)
}

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(s string) string {
	// Handle acronyms and special cases
	s = handleAcronyms(s)

	// Insert underscores before uppercase letters (except at the beginning)
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			// Check if the previous character is lowercase or if this is the start of a new word
			prevR := rune(s[i-1])
			if unicode.IsLower(prevR) || (i < len(s)-1 && unicode.IsLower(rune(s[i+1]))) {
				result.WriteRune('_')
			}
		}
		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}

// handleAcronyms handles common acronyms in field names
func handleAcronyms(s string) string {
	// Common acronyms and their replacements
	acronyms := map[string]string{
		"URL":   "Url",
		"API":   "Api",
		"HTTP":  "Http",
		"HTTPS": "Https",
		"IP":    "Ip",
		"JSON":  "Json",
		"XML":   "Xml",
		"SQL":   "Sql",
		"UUID":  "Uuid",
		"ID":    "Id",
		"EZBEQ": "Ezbeq",
		"AVR":   "Avr",
		"HDMI":  "Hdmi",
		"TV":    "Tv",
		"TMDB":  "Tmdb",
	}

	result := s
	for acronym, replacement := range acronyms {
		// Replace acronym at the end of string
		if strings.HasSuffix(result, acronym) {
			result = strings.TrimSuffix(result, acronym) + replacement
		}
		// Replace acronym followed by uppercase letter
		re := regexp.MustCompile(acronym + "([A-Z])")
		result = re.ReplaceAllString(result, replacement+"$1")
	}

	return result
}

// addTagToField adds a tag to a field's tag list
func addTagToField(field *ast.Field, key, value string) {
	newTag := fmt.Sprintf(`%s:%q`, key, value)

	if field.Tag == nil {
		// Create new tag
		field.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: "`" + newTag + "`",
		}
	} else {
		// Add to existing tags
		existingTag := strings.Trim(field.Tag.Value, "`")
		if existingTag == "" {
			field.Tag.Value = "`" + newTag + "`"
		} else {
			field.Tag.Value = "`" + existingTag + " " + newTag + "`"
		}
	}
}

// updateJSONTag updates the json tag value for a field
func updateJSONTag(field *ast.Field, value string) {
	if field.Tag == nil {
		return
	}

	tagValue := strings.Trim(field.Tag.Value, "`")

	// Parse existing tags
	tags := parseStructTags(tagValue)
	tags["json"] = value

	// Rebuild tag string
	tagParts := make([]string, 0, len(tags))
	for k, v := range tags {
		tagParts = append(tagParts, fmt.Sprintf(`%s:%q`, k, v))
	}

	field.Tag.Value = "`" + strings.Join(tagParts, " ") + "`"
}

// parseStructTags parses a struct tag string into a map
func parseStructTags(tag string) map[string]string {
	tags := make(map[string]string)

	// Simple regex to match tag:value pairs
	re := regexp.MustCompile(`(\w+):"([^"]*)"`)
	matches := re.FindAllStringSubmatch(tag, -1)

	for _, match := range matches {
		if len(match) == 3 {
			tags[match[1]] = match[2]
		}
	}

	return tags
}
