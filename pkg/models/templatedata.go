package models

// Template Data type to pass data to templates
// 		In HTML files, the required data will be defined with a name. But the actual data type may
//		be something else. This is why we have couple of map types. The last Data field is a mapping
//		from strings to ANY object Type
// CSRFToken is for Cross Site Request Forgery if a template has form in it.
// Flash, Warning and Error are just messages
type TemplateData struct {
	StringMap map[string]string
	IntMap map[string]int
	FloatMap map[string]float64
	Data map[string]interface{}
	CSRFToken string
	Flash string
	Warning string
	Error string
}