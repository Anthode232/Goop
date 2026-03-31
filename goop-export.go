package goop

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportData represents structured data for export
type ExportData struct {
	Elements []ElementData `json:"elements"`
}

// ElementData represents a single element's data
type ElementData struct {
	Tag      string            `json:"tag"`
	Text     string            `json:"text"`
	HTML     string            `json:"html"`
	Attrs    map[string]string `json:"attributes"`
	Children int               `json:"children_count"`
}

// ToJSON converts a Root element to JSON string
func (r Root) ToJSON() (string, error) {
	if r.Error != nil {
		return "", r.Error
	}

	data := r.ToExportData()
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %v", err)
	}

	return string(jsonBytes), nil
}

// ToCSV converts a Root element and its children to CSV format
func (r Root) ToCSV() (string, error) {
	if r.Error != nil {
		return "", r.Error
	}

	var records [][]string

	// Add header
	records = append(records, []string{"Tag", "Text", "HTML", "Attributes", "Children"})

	// Convert element and its children to records
	elementRecords := r.toCSVRecords()
	records = append(records, elementRecords...)

	var csvBuilder strings.Builder
	csvWriter := csv.NewWriter(&csvBuilder)

	for _, record := range records {
		if err := csvWriter.Write(record); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV: %v", err)
	}

	return csvBuilder.String(), nil
}

// ToExportData converts a Root element to ExportData structure
func (r Root) ToExportData() ExportData {
	data := ExportData{}

	if r.Error == nil {
		elementData := r.toElementData()
		data.Elements = append(data.Elements, elementData)

		// Add children
		children := r.Children()
		for _, child := range children {
			childData := child.ToExportData()
			data.Elements = append(data.Elements, childData.Elements...)
		}
	}

	return data
}

// SaveJSON saves the element data to a JSON file
func (r Root) SaveJSON(filename string) error {
	jsonData, err := r.ToJSON()
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(jsonData), 0644)
}

// SaveCSV saves the element data to a CSV file
func (r Root) SaveCSV(filename string) error {
	csvData, err := r.ToCSV()
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(csvData), 0644)
}

// ExportAllToJSON converts multiple Root elements to JSON
func ExportAllToJSON(elements []Root) (string, error) {
	data := ExportData{}

	for _, element := range elements {
		if element.Error == nil {
			elementData := element.toElementData()
			data.Elements = append(data.Elements, elementData)
		}
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %v", err)
	}

	return string(jsonBytes), nil
}

// ExportAllToCSV converts multiple Root elements to CSV
func ExportAllToCSV(elements []Root) (string, error) {
	var records [][]string

	// Add header
	records = append(records, []string{"Tag", "Text", "HTML", "Attributes", "Children"})

	for _, element := range elements {
		if element.Error == nil {
			elementRecords := element.toCSVRecords()
			records = append(records, elementRecords...)
		}
	}

	var csvBuilder strings.Builder
	csvWriter := csv.NewWriter(&csvBuilder)

	for _, record := range records {
		if err := csvWriter.Write(record); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV: %v", err)
	}

	return csvBuilder.String(), nil
}

// SaveAllJSON saves multiple elements to a JSON file
func SaveAllJSON(elements []Root, filename string) error {
	jsonData, err := ExportAllToJSON(elements)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(jsonData), 0644)
}

// SaveAllCSV saves multiple elements to a CSV file
func SaveAllCSV(elements []Root, filename string) error {
	csvData, err := ExportAllToCSV(elements)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(csvData), 0644)
}

// toElementData converts Root to ElementData
func (r Root) toElementData() ElementData {
	data := ElementData{
		Tag:   r.NodeValue,
		Text:  safeText(r),
		HTML:  r.HTML(),
		Attrs: safeAttrs(r),
	}

	if r.Pointer != nil {
		data.Children = len(r.Children())
	}

	return data
}

// safeText safely extracts text without panicking
func safeText(r Root) string {
	defer func() {
		if r := recover(); r != nil {
			// Handle panic from Text() method
		}
	}()
	return r.Text()
}

// safeAttrs safely extracts attributes without panicking
func safeAttrs(r Root) map[string]string {
	defer func() {
		if r := recover(); r != nil {
			// Handle panic from Attrs() method
		}
	}()
	return r.Attrs()
}

// toCSVRecords converts Root to CSV records
func (r Root) toCSVRecords() [][]string {
	var records [][]string

	if r.Error == nil {
		attrsStr := formatAttributes(safeAttrs(r))
		text := sanitizeCSVField(safeText(r))
		html := sanitizeCSVField(r.HTML())

		record := []string{
			r.NodeValue,
			text,
			html,
			attrsStr,
			fmt.Sprintf("%d", len(r.Children())),
		}
		records = append(records, record)

		// Add children records
		children := r.Children()
		for _, child := range children {
			childRecords := child.toCSVRecords()
			records = append(records, childRecords...)
		}
	}

	return records
}

// formatAttributes converts attributes map to string
func formatAttributes(attrs map[string]string) string {
	if len(attrs) == 0 {
		return ""
	}

	var pairs []string
	for key, value := range attrs {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(pairs, "; ")
}

// sanitizeCSVField cleans text for CSV output
func sanitizeCSVField(text string) string {
	// Remove newlines and tabs, replace with spaces
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Truncate very long text
	if len(text) > 1000 {
		text = text[:997] + "..."
	}

	return text
}

// LoadJSON loads JSON data from file and converts to Root elements
func LoadJSON(filename string) ([]Root, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var exportData ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	var elements []Root
	for _, elementData := range exportData.Elements {
		// Note: This creates simplified Root elements from exported data
		// Full HTML reconstruction would require more complex parsing
		element := Root{
			NodeValue: elementData.Tag,
			Error:     nil,
		}
		elements = append(elements, element)
	}

	return elements, nil
}

// LoadCSV loads CSV data from file
func LoadCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %v", err)
	}

	return records, nil
}
