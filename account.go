package amo

const (
	accountEntity = "account"
)

type Account struct {
	Id                      int               `json:"id"`
	Name                    string            `json:"name"`
	Subdomain               string            `json:"subdomain"`
	CreatedAt               int               `json:"created_at"`
	CreatedBy               int               `json:"created_by"`
	UpdatedAt               int               `json:"updated_at"`
	UpdatedBy               int               `json:"updated_by"`
	CurrentUser             int               `json:"current_user"`
	Country                 string            `json:"country"`
	CustomersMode           string            `json:"customers_mode"`
	IsUnsortedOn            bool              `json:"is_unsorted_on"`
	MobileFeatureVersion    int               `json:"mobile_feature_version"`
	IsLossReasonEnabled     bool              `json:"is_loss_reason_enabled"`
	IsHelpbotEnabled        bool              `json:"is_helpbot_enabled"`
	IsTechnicalAccount      bool              `json:"is_technical_account"`
	ContactNameDisplayOrder int               `json:"contact_name_display_order"`
	AmojoId                 *string           `json:"amojo_id,omitempty"`
	Uuid                    *string           `json:"uuid,omitempty"`
	Version                 *int              `json:"version,omitempty"`
	DatetimeSettings        *DatetimeSettings `json:"datetime_settings"`
	Embedded                AccountEmbedded   `json:"_embedded"`
}

type AccountEmbedded struct {
	AmojoRights map[string]bool `json:"amojo_rights,omitempty"`
	UsersGroups []IdNameItem    `json:"users_groups,omitempty"`
	TaskTypes   []TaskType      `json:"task_types,omitempty"`
}

type TaskType struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Color  *string `json:"color"`
	IconId *string `json:"icon_id"`
	Code   string  `json:"code"`
}

type DatetimeSettings struct {
	DatePattern      string `json:"date_pattern"`
	ShortDatePattern string `json:"short_date_pattern"`
	ShortTimePattern string `json:"short_time_pattern"`
	DateFormat       string `json:"date_format"`
	TimeFormat       string `json:"time_format"`
	Timezone         string `json:"timezone"`
	TimezoneOffset   string `json:"timezone_offset"`
}

// GetAccount return current Account
// Available with WithAmojoId, WithUuid, WithAmojoRights, WithUsersGroups, WithVersion, WithDatetimeSettings
func (a *AmoCrm) GetAccount(with []string) (*Account, error) {
	var account *Account
	return account, a.getItem([]string{accountEntity}, nil, &entitiesQuery{With: with}, &account)
}
