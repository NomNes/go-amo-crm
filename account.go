package amo

import (
	"encoding/json"
	"strings"

	"github.com/NomNes/go-errors-sentry"
)

type Account struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Subdomain      string `json:"subdomain"`
	Currency       string `json:"currency"`
	Timezone       string `json:"timezone"`
	TimezoneOffset string `json:"timezone_offset"`
	Language       string `json:"language"`
	CurrentUser    int    `json:"current_user"`
	Embedded       struct {
		*CustomFields
		CF *uglyCustomFields `json:"custom_fields,omitempty"`
		// users	array	Список пользователей аккаунта с их правами
		// users/id	int	Уникальный идентификатор пользователя
		// users/name	string	Имя пользователя
		// users/login	string	Логин пользователя
		// users/language	string	Настройки языка пользователя
		// users/phone_number	string	Номер телефона пользователя
		// users/group_id	int	id группы, в которой состоит пользователь
		// users/is_active	bool	Активна учётная запись пользователя или нет, если нет, то доступ будет закрыт
		// users/is_free	bool	Является ли учётная запись пользователя бесплатной
		// users/is_admin	bool	Наличие прав администратора
		// users/rights	array	Права пользователя (описание формата см. здесь)
		// users/rights/mail	string	Доступ к корпоративной почте
		// users/rights/incoming_leads	string	Доступ к “неразобранному”
		// users/rights/catalogs	string	Права пользователя на создание/редактирование каталогов и их элементов
		// users/rights/lead_add	string	Права пользователя на добавление новых сделок
		// users/rights/lead_view	string	Права пользователя на просмотра существующих сделок
		// users/rights/lead_edit	string	Права пользователя на редактирование существующих сделок
		// users/rights/lead_delete	string	Права пользователя на удаление существующих сделок
		// users/rights/lead_export	string	Права пользователя на экспорт сделок
		// users/rights/contact_add	string	Права пользователя на добавление новых контактов
		// users/rights/contact_view	string	Права пользователя на просмотр существующих контактов
		// users/rights/contact_edit	string	Права пользователя на редактирование существующих контактов
		// users/rights/contact_delete	string	Права пользователя на удаление существующих контактов
		// users/rights/contact_export	string	Права пользователя на экспорт контактов
		// users/rights/company_add	string	Права пользователя на добавление новых компаний
		// users/rights/company_view	string	Права пользователя на просмотр существующих компаний
		// users/rights/company_edit	string	Права пользователя на редактирование существующих компаний
		// users/rights/company_delete	string	Права пользователя на удаление существующих компаний
		// users/rights/company_export	string	Права пользователя на экспорт существующих компаний
		// note_types	array	Список используемых в системе типов примечаний (подробное описание типов см. здесь)
		// note_types/id	int	Уникальный идентификатор примечания
		// note_types/name	string	Название примечания
		// note_types/code	string	Код примечания
		// note_types/editable	bool	Существует ли возможность редактирования примечания
		// task_types	array	Типы задач доступных для данного аккаунта
		// task_types/id	int	Уникальный идентификатор задачи
		// task_types/name	string	Название задачи
		// pipelines	array	Цифровые воронки имеющиеся на аккаунте
		// pipelines/id	int	Уникальный идентификатор воронки
		// pipelines/name	string	Название воронки
		// pipelines/sort	int	Порядковый номер воронки при отображении
		// pipelines/is_main	bool	Является ли воронка “главной”
		// pipelines/statuses	array	Этапы цифровой воронки
		// pipelines/statuses/id	int	Уникальный идентификатор этапа
		// pipelines/statuses/name	string	Название этапа
		// pipelines/statuses/sort	int	Порядковый номер этапа при отображении
		// pipelines/statuses/color	string	Цвет этапа (подробнее см. здесь)
		// pipelines/statuses/editable	bool	Есть ли возможность изменить или удалить этот этап
	} `json:"_embedded"`
	// date_pattern	array	Формат даты (описание формата см. здесь)
	// date_pattern/date	string	Дата, формат зависит от выбранного формата в аккаунте
	// date_pattern/time	string	Время, формат зависит от выбранного формата в аккаунте
	// date_pattern/date_time	string	Дата и время, формат зависит от выбранного формата в аккаунте
	// date_pattern/time_full	string	Время с точностью до секунды, формат зависит от выбранного формата в аккаунте
}

type uglyCustomFields struct {
	Contacts  interface{} `json:"contacts"`
	Leads     interface{} `json:"leads"`
	Companies interface{} `json:"companies"`
	Customers interface{} `json:"customers"`
}

func reUnmarshal(s interface{}, d interface{}) error {
	j, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err)
	}
	return errors.Wrap(json.Unmarshal(j, &d))
}

func (a *AmoCrm) GetAccount(with []string) (*Account, error) {
	path := "/api/v2/account"
	if with != nil {
		path += "?with=" + strings.Join(with, ",")
	}
	var r *Account
	err := a.get(path, &r)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if r.Embedded.CF != nil {
		r.Embedded.CustomFields = &CustomFields{}
		if f, ok := r.Embedded.CF.Contacts.(map[string]interface{}); ok {
			err = reUnmarshal(f, &r.Embedded.CustomFields.Contacts)
			if err != nil {
				return nil, errors.Wrap(err)
			}
		}
		if f, ok := r.Embedded.CF.Leads.(map[string]Field); ok {
			r.Embedded.CustomFields.Leads = f
		}
		if f, ok := r.Embedded.CF.Companies.(map[string]Field); ok {
			r.Embedded.CustomFields.Companies = f
		}
		if f, ok := r.Embedded.CF.Customers.(map[string]Field); ok {
			r.Embedded.CustomFields.Customers = f
		}
		r.Embedded.CF = nil
	}
	return r, nil
}
