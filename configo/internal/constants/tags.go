package constants

// ConfigTag is a type for configuration tags.
type ConfigTag string

const (
	// CONFIG_TAG_NAME is the tag for the name of the configuration. It is required tag.
	CONFIG_TAG_NAME ConfigTag = "name"
	// CONFIG_TAG_DESCRIPTION is the tag for the description of the configuration. It is optional tag.
	CONFIG_TAG_DESCRIPTION ConfigTag = "description"
	// CONFIG_TAG_TYPE is the tag for the type of the configuration. It is required
	CONFIG_TAG_TYPE ConfigTag = "type"
	// CONFIG_TAG_REQUIRED is the tag for the required value of the configuration. It should be true or false. It will be false by default.
	CONFIG_TAG_REQUIRED ConfigTag = "required"
	// CONFIG_TAG_CHOICES is the tag for the choices of the configuration. It is required if the type is choice.
	CONFIG_TAG_CHOICES ConfigTag = "choices"
)

func (ct ConfigTag) String() string {
	return string(ct)
}
