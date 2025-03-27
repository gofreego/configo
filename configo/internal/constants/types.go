package constants

type ConfigType string

const (
	// CONFIG_TYPE_STRING is the type for string configuration, it will show a textbox on ui.
	CONFIG_TYPE_STRING ConfigType = "string"
	// CONFIG_TYPE_NUMBER is the type for integer configuration, it will show a number input on ui.
	CONFIG_TYPE_NUMBER ConfigType = "number"
	// CONFIG_TYPE_BOOLEAN is the type for boolean configuration, it will show a checkbox on ui.
	CONFIG_TYPE_BOOLEAN ConfigType = "boolean"
	// CONFIG_TYPE_JSON is the type for json configuration, it will show a textarea on ui which will have json formatting.
	CONFIG_TYPE_JSON ConfigType = "json"
	// CONFIG_TYPE_FLOAT is the type for float configuration, it will show a number input on ui.
	CONFIG_TYPE_BIG_TEXT ConfigType = "bigText"
	// CONFIG_TYPE_CHOICE is the type for choice configuration, it will show a dropdown on ui and it should have type string.
	CONFIG_TYPE_CHOICE ConfigType = "choice"
	//CONFIG_TYPE_PARENT
	CONFIG_TYPE_PARENT ConfigType = "parent"
	//CONFIG_TYPE_LIST
	CONFIG_TYPE_LIST ConfigType = "list"
)
