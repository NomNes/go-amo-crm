package amo

const pipelineEntity = "pipelines"

type Pipeline struct {
	Id           int              `json:"id,omitempty"`
	Name         string           `json:"name"`
	Sort         int              `json:"sort"`
	IsMain       bool             `json:"is_main"`
	IsUnsortedOn bool             `json:"is_unsorted_on"`
	IsArchive    *bool            `json:"is_archive,omitempty"`
	AccountId    *int             `json:"account_id,omitempty"`
	Embedded     PipelineEmbedded `json:"_embedded,omitempty"`
	client       *AmoCrm
}

type PipelineEmbedded struct {
	Statuses []PipelineStatus `json:"statuses,omitempty"`
}

type PipelineStatus struct {
	Id         int    `json:"id,omitempty"`
	Name       string `json:"name"`
	Sort       int    `json:"sort,omitempty"`
	IsEditable bool   `json:"is_editable,omitempty"`
	PipelineId int    `json:"pipeline_id,omitempty"`
	Color      string `json:"color,omitempty"`
	Type       int    `json:"type,omitempty"`
	AccountId  int    `json:"account_id,omitempty"`
}

// GetPipelines return slice of Pipeline
func (a *AmoCrm) GetPipelines() ([]Pipeline, *Pages, error) {
	var pipelines []Pipeline
	pages, err := a.getItems([]string{leadsEntity, pipelineEntity}, nil, &pipelines)
	for i := range pipelines {
		pipelines[i].client = a
	}
	return pipelines, pages, err
}

// GetPipeline return Pipeline by id
func (a *AmoCrm) GetPipeline(id int) (*Pipeline, error) {
	var pipeline *Pipeline
	err := a.getItem([]string{leadsEntity, pipelineEntity}, &id, nil, &pipeline)
	if err != nil {
		return nil, err
	}
	pipeline.client = a
	return pipeline, nil
}

// NewPipeline return Pipeline for current AmoCrm client
func (a *AmoCrm) NewPipeline(contact *Pipeline) *Pipeline {
	if contact == nil {
		contact = &Pipeline{}
	}
	contact.client = a
	return contact
}

// Save add or update Pipeline
func (p *Pipeline) Save() error {
	return p.client.addItem([]string{leadsEntity, pipelineEntity}, p, p.Id > 0, &p.Id)
}

// Delete Pipeline
func (p *Pipeline) Delete() error {
	return p.client.deleteItem([]string{leadsEntity, pipelineEntity}, p.Id)
}
