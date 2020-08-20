package amo

import "fmt"

type linkData struct {
	Id       int           `json:"to_entity_id"`
	Type     string        `json:"to_entity_type"`
	Metadata *linkMetadata `json:"metadata,omitempty"`
}

type linkMetadata struct {
	CatalogId int  `json:"catalog_id,omitempty"`
	Quantity  int  `json:"quantity,omitempty"`
	IsMain    bool `json:"is_main,omitempty"`
	UpdatedBy bool `json:"updated_by,omitempty"`
}

func (a *AmoCrm) link(entity string, id int, data []linkData) error {
	var r interface{}
	return a.apiPost(fmt.Sprintf("/%s/%d/link", entity, id), data, &r)
}
