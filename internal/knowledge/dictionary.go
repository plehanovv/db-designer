package knowledge

type WordType string

const (
	EntityType    WordType = "ENTITY"
	AttributeType WordType = "ATTRIBUTE"
	VerbType      WordType = "VERB"
)

var Dictionary = map[string]WordType{

	// entities
	"student":    EntityType,
	"teacher":    EntityType,
	"course":     EntityType,
	"group":      EntityType,
	"university": EntityType,

	// attributes
	"name":  AttributeType,
	"email": AttributeType,
	"age":   AttributeType,
	"phone": AttributeType,
	"title": AttributeType,

	// verbs
	"has":      VerbType,
	"contains": VerbType,
	"includes": VerbType,
	"teaches":  VerbType,
	"enrolls":  VerbType,
	"belongs":  VerbType,
}
