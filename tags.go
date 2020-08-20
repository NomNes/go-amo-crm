package amo

const tagsEntity = "tags"

type Tag IdNameItem

func (a *AmoCrm) getTags(entity string, limit, page int) ([]Tag, *Pages, error) {
	var tags []Tag
	pages, err := a.getItems([]string{entity, tagsEntity}, &entitiesQuery{Limit: limit, Page: page}, &tags)
	return tags, pages, err
}
