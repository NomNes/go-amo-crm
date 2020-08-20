package amo

const (
	FieldTypeText          = "text"
	FieldTypeNumeric       = "numeric"
	FieldTypeCheckbox      = "checkbox"
	FieldTypeSelect        = "select"
	FieldTypeMultiselect   = "multiselect"
	FieldTypeMultiText     = "multitext"
	FieldTypeDate          = "date"
	FieldTypeUrl           = "url"
	FieldTypeTextarea      = "textarea"
	FieldTypeRadiobutton   = "radiobutton"
	FieldTypeStreetAddress = "streetaddress"
	FieldTypeSmartAddress  = "smart_address"
	FieldBirthday          = "birthday"
	FieldTypeLegalEntity   = "legal_entity"
	FieldTypeDateTime      = "date_time"
	FieldTypeItems         = "items"
	FieldTypeCategory      = "category"
	FieldTypePrice         = "price"

	customFieldsEntity = "custom_fields"
)

type CustomFieldValues struct {
	Id     int          `json:"field_id"`
	Name   string       `json:"field_name"`
	Code   *string      `json:"field_code"`
	Type   string       `json:"field_type"`
	Values []FieldValue `json:"values"`
}

type FieldValue struct {
	Value    interface{} `json:"value"`
	EnumId   int         `json:"enum_id,omitempty"`
	EnumCode string      `json:"enum_code,omitempty"`
}

type CustomField struct {
	Id               int                         `json:"id"`
	Name             string                      `json:"name"`
	Type             string                      `json:"type"`
	AccountId        int                         `json:"account_id"`
	Code             string                      `json:"code"`
	Sort             int                         `json:"sort"`
	IsApiOnly        bool                        `json:"is_api_only"`
	Enums            []CustomFieldEnum           `json:"enums"`
	GroupId          *int                        `json:"group_id"`
	RequiredStatuses []CustomFieldRequiredStatus `json:"required_statuses"`
	IsDeletable      bool                        `json:"is_deletable"`
	IsPredefined     bool                        `json:"is_predefined"`
	EntityType       string                      `json:"entity_type"`
	Remind           *string                     `json:"remind"`
}

type CustomFieldEnum struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
	Sort  int    `json:"sort"`
}

type CustomFieldRequiredStatus struct {
	PipelineId int `json:"pipeline_id"`
	StatusId   int `json:"status_id"`
}

func (a *AmoCrm) getCustomFields(entity string, limit, page int) ([]CustomField, *Pages, error) {
	var fields []CustomField
	pages, err := a.getItems([]string{entity, customFieldsEntity}, &entitiesQuery{Limit: limit, Page: page}, &fields)
	return fields, pages, err
}
