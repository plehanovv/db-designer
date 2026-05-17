package analyzer

import (
	"db-designer-vkr/internal/model"
	"strings"
)

var ignoredEntities = map[string]bool{
	"The":  true,
	"This": true,
	"That": true,
}

func ValidateEntityCandidates(
	candidates []model.EntityCandidate,
) []model.EntityCandidate {

	var valid []model.EntityCandidate

	for _, candidate := range candidates {

		if ignoredEntities[candidate.Name] {
			continue
		}

		if len(candidate.Name) < 2 {
			continue
		}

		candidate.Name = normalizeEntityName(candidate.Name)

		valid = append(valid, candidate)
	}

	return valid
}

func normalizeEntityName(name string) string {

	name = strings.TrimSpace(name)

	return strings.Title(strings.ToLower(name))
}
