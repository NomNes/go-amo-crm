package amo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func (a *AmoCrm) getHost() string {
	return "https://" + a.subdomain + ".amocrm.ru"
}

func (a *AmoCrm) getApiUri() string {
	return fmt.Sprintf("/api/v%d", VERSION)
}

type ErrorRes struct {
	Hint             string `json:"hint"`
	Title            string `json:"title"`
	Type             string `json:"type"`
	Status           int    `json:"status"`
	Detail           string `json:"detail"`
	ValidationErrors []struct {
		Errors []struct {
			Code    string `json:"code"`
			Path    string `json:"path"`
			Details string `json:"details"`
		} `json:"errors"`
	} `json:"validation-errors"`
}

func (a *AmoCrm) request(method, path string, jsonBody interface{}, r interface{}, auth bool) error {
	if auth {
		err := a.Restore(false)
		if err != nil {
			return err
		}
	}
	var br io.Reader
	if jsonBody != nil {
		b, err := json.Marshal(jsonBody)
		if err != nil {
			return err
		}
		br = bytes.NewReader(b)
		// log.Println(">>>", method, path, string(b))
	} else {
		// log.Println(">>>", method, path)
	}
	req, err := http.NewRequest(method, a.getHost()+path, br)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "amoCRM-API-Library/1.0")
	if auth {
		if d := a.Storage.Get(); d != nil {
			req.Header.Set("Authorization", fmt.Sprintf("%s %s", d.TokenType, d.AccessToken))
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	// log.Println("<<<", method, path, string(body))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		var errJson ErrorRes
		log.Println(string(body))
		err := json.Unmarshal(body, &errJson)
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s\n%s", path, err, string(body)))
		}
		return errors.New(fmt.Sprintf("%d %s\n%s\n%s\n%s\n%s", errJson.Status, errJson.Title, errJson.Hint, errJson.Detail, errJson.Type, errJson.ValidationErrors))
	}
	return json.Unmarshal(body, &r)
}

func (a *AmoCrm) post(path string, jsonBody interface{}, r interface{}, auth bool) error {
	return a.request(http.MethodPost, path, jsonBody, &r, auth)
}

func (a *AmoCrm) patch(path string, jsonBody interface{}, r interface{}, auth bool) error {
	return a.request(http.MethodPatch, path, jsonBody, &r, auth)
}

func serializeQuery(query map[string]string) string {
	var queryItems []string
	for key, value := range query {
		queryItems = append(queryItems, fmt.Sprintf("%s=%s", key, value))
	}
	if len(queryItems) > 0 {
		return fmt.Sprintf("?%s", strings.Join(queryItems, "&"))
	}
	return ""
}

func (a *AmoCrm) get(path string, query map[string]string, r interface{}, auth bool) error {
	return a.request(http.MethodGet, path+serializeQuery(query), nil, &r, auth)
}

func (a *AmoCrm) delete(path string, query map[string]string, r interface{}, auth bool) error {
	return a.request(http.MethodDelete, path+serializeQuery(query), nil, &r, auth)
}

func (a *AmoCrm) apiPost(path string, jsonBody interface{}, r interface{}) error {
	return a.post(a.getApiUri()+path, jsonBody, &r, true)
}

func (a *AmoCrm) apiPatch(path string, jsonBody interface{}, r interface{}) error {
	return a.patch(a.getApiUri()+path, jsonBody, &r, true)
}

func (a *AmoCrm) apiGet(path string, query map[string]string, r interface{}) error {
	return a.get(a.getApiUri()+path, query, &r, true)
}

func (a *AmoCrm) apiDelete(path string, query map[string]string, r interface{}) error {
	return a.delete(a.getApiUri()+path, query, &r, true)
}

const (
	WithLeads                  = "leads"
	WithCustomers              = "customers"
	WithCatalogElements        = "catalog_elements"
	WithAmojoId                = "amojo_id"
	WithUuid                   = "uuid"
	WithAmojoRights            = "amojo_rights"
	WithUsersGroups            = "users_groups"
	WithTaskTypes              = "task_types"
	WithVersion                = "version"
	WithDatetimeSettings       = "datetime_settings"
	WithIsPriceModifiedByRobot = "is_price_modified_by_robot"
	WithLossReason             = "loss_reason"
	WithContacts               = "contacts"
	WithOnlyDeleted            = "only_deleted"

	OrderByCreatedAt = "created_at"
	OrderByUpdatedAt = "updated_at"
	OrderById        = "id"
	OrderAsc         = "asc"
	OrderDesc        = "desc"
)

type entitiesQuery struct {
	Id    int
	Limit int
	Page  int
	With  []string
	Query string
	Order map[string]string
}

type Pages struct {
	Total int `json:"_total_items"`
	Page  int `json:"_page"`
	Pages int `json:"_page_count"`
}

type entitiesResponse struct {
	*Pages
	Embedded map[string]interface{} `json:"_embedded"`
}

func buildQuery(query *entitiesQuery) map[string]string {
	q := map[string]string{}
	if query != nil {
		if query.Id > 0 {
			q["id"] = fmt.Sprintf("%d", query.Id)
		}
		if query.Limit > 0 {
			q["limit_rows"] = fmt.Sprintf("%d", query.Limit)
		}
		if query.Page > 0 {
			q["limit_offset"] = fmt.Sprintf("%d", query.Page)
		}
		if len(query.With) > 0 {
			q["with"] = strings.Join(query.With, ",")
		}
		if query.Query != "" {
			q["query"] = query.Query
		}
		if query.Order != nil && len(query.Order) > 0 {
			for key, dem := range query.Order {
				q[fmt.Sprintf("order[%s]", key)] = dem
			}
		}
	}
	return q
}

func (a *AmoCrm) getItems(entity []string, query *entitiesQuery, v interface{}) (*Pages, error) {
	var r entitiesResponse
	var err error
	err = a.apiGet("/"+strings.Join(entity, "/"), buildQuery(query), &r)
	if err != nil {
		return r.Pages, err
	}
	if e, ok := r.Embedded[entity[len(entity)-1]]; ok {
		err = reMarshal(e, &v)
	}
	return r.Pages, err
}

func (a *AmoCrm) getItem(entity []string, id *int, query *entitiesQuery, v interface{}) error {
	path := "/" + strings.Join(entity, "/")
	if id != nil {
		path += fmt.Sprintf("/%d", *id)
	}
	return a.apiGet(path, buildQuery(query), &v)
}

type EntityWithTime interface {
	GetCreatedAtTime() time.Time
	GetUpdatedAtTime() time.Time
	SetCreatedAtTime(time.Time)
	SetUpdatedAtTime(time.Time)
}

func (a *AmoCrm) addItem(entity []string, v interface{}, update bool, id *int) error {
	var r map[string]interface{}
	var err error
	path := "/" + strings.Join(entity, "/")
	if update {
		if id != nil {
			path += fmt.Sprintf("/%d", *id)
			err = a.apiPatch(path, v, &r)
		} else {
			err = a.apiPatch(path, []interface{}{v}, &r)
		}
	} else {
		err = a.apiPost(path, []interface{}{v}, &r)
	}
	if err != nil {
		return err
	}
	if _, ok := r["id"]; ok {
		return reMarshal(r, &v)
	} else {
		if embedded, ok := r["_embedded"]; ok {
			if me, ok := embedded.(map[string]interface{}); ok {
				if e, ok := me[entity[len(entity)-1]]; ok {
					if s, ok := e.([]interface{}); ok && len(s) > 0 {
						if m, ok := s[0].(map[string]interface{}); ok {
							if id, ok := m["id"]; ok {
								if floatId, ok := id.(float64); ok {
									intId := int(floatId)
									return a.getItem(entity, &intId, nil, &v)
								}
							}
						}
					}
				}
			}
		}
	}
	return errors.New("unknown error")
}

func (a *AmoCrm) deleteItem(entity []string, id int) error {
	var r entitiesResponse
	return a.apiDelete(fmt.Sprintf("/%s/%d", strings.Join(entity, "/"), id), nil, &r)
}
