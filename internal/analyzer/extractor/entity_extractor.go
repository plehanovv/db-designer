package extractor

import (
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func ExtractEntityCandidates(
	document nlp.Document,
) []model.EntityCandidate {

	var candidates []model.EntityCandidate

	unique := map[string]bool{}

	for _, sentence := range document.Sentences {

		for _, token := range sentence.Tokens {

			if token.Type != "ENTITY" {
				continue
			}

			if unique[token.Value] {
				continue
			}

			unique[token.Value] = true

			candidates = append(candidates, model.EntityCandidate{
				Name:       token.Value,
				Confidence: 0.8,
				Source:     sentence.Text,
			})
		}
	}

	return candidates
}
