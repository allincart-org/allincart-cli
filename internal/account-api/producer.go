package account_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ProducerEndpoint struct {
	c          *Client
	producerId int
}

func (e ProducerEndpoint) GetId() int {
	return e.producerId
}

func (c *Client) Producer(ctx context.Context) (*ProducerEndpoint, error) {
	r, err := c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/tenant/allocations", ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(r)
	if err != nil {
		return nil, err
	}

	var allocation companyAllocation
	if err := json.Unmarshal(body, &allocation); err != nil {
		return nil, fmt.Errorf("producer.profile: %v", err)
	}

	if !allocation.IsProducer {
		return nil, fmt.Errorf("this company is not unlocked as producer")
	}

	return &ProducerEndpoint{producerId: allocation.ProducerID, c: c}, nil
}

type companyAllocation struct {
	IsProducer bool `json:"isProducer"`
	ProducerID int  `json:"producerId"`
}

func (e ProducerEndpoint) Profile(ctx context.Context) (*Producer, error) {
	r, err := e.c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/producers", ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, err
	}

	var producers []Producer
	if err := json.Unmarshal(body, &producers); err != nil {
		return nil, fmt.Errorf("my_profile: %v", err)
	}

	for _, profile := range producers {
		return &profile, nil
	}

	return nil, fmt.Errorf("cannot find a profile")
}

type Producer struct {
	Prefix  string `json:"prefix"`
	Name    string `json:"name"`
	Website string `json:"website"`
}
type Query struct {
	Type  string      `json:"type"`
	Field string      `json:"field"`
	Value interface{} `json:"value"` // 用 interface{} 支持不同类型值（如 int、string）
}

type ListExtensionCriteria struct {
	Limit         int    `schema:"limit,omitempty"`
	Offset        int    `schema:"offset,omitempty"`
	OrderBy       string `schema:"orderBy,omitempty"`
	OrderSequence string `schema:"orderSequence,omitempty"`
	Query         *Query `json:"query,omitempty"`
}

func (e ProducerEndpoint) Extensions(ctx context.Context, criteria *ListExtensionCriteria) ([]Extension, error) {
	form := url.Values{}
	form.Set("producerId", strconv.FormatInt(int64(e.GetId()), 10))

	if criteria.Limit != 0 {
		form.Set("limit", strconv.Itoa(criteria.Limit))
	}
	if criteria.Offset != 0 {
		form.Set("offset", strconv.Itoa(criteria.Offset))
	}
	if criteria.OrderBy != "" {
		form.Set("orderBy", criteria.OrderBy)
	}
	if criteria.OrderSequence != "" {
		form.Set("orderSequence", criteria.OrderSequence)
	}

	if criteria.Query != nil {
		query, err := json.Marshal(criteria.Query)
		if err != nil {
			return nil, fmt.Errorf("list_extensions: %v", err)
		}
		form.Set("query", string(query))
	}

	r, err := e.c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/plugins?%s", ApiUrl, form.Encode()), nil)
	if err != nil {
		return nil, err
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, err
	}

	var extensions []Extension
	if err := json.Unmarshal(body, &extensions); err != nil {
		return nil, fmt.Errorf("list_extensions: %v", err)
	}

	return extensions, nil
}

func (e ProducerEndpoint) GetExtensionByName(ctx context.Context, name string) (*Extension, error) {
	criteria := ListExtensionCriteria{
		Query: &Query{
			Type:  "equals",
			Field: "productNumber",
			Value: name,
		},
	}

	extensions, err := e.Extensions(ctx, &criteria)
	if err != nil {
		return nil, err
	}

	for _, ext := range extensions {
		if strings.EqualFold(ext.Name, name) {
			return e.GetExtensionById(ctx, ext.Id)
		}
	}

	return nil, fmt.Errorf("cannot find Extension by name %s", name)
}

func (e ProducerEndpoint) GetExtensionById(ctx context.Context, id int) (*Extension, error) {
	errorFormat := "GetExtensionById: %v"

	// Create it
	r, err := e.c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/plugins/%d", ApiUrl, id), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	var extension Extension
	if err := json.Unmarshal(body, &extension); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return &extension, nil
}

type Extension struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	SubType        string `json:"subType"`
	StandardLocale Locale `json:"standardLocale"`
	Infos          []*struct {
		Locale             Locale       `json:"locale"`
		Name               string       `json:"name"`
		Description        string       `json:"description"`
		InstallationManual string       `json:"installationManual"`
		ShortDescription   string       `json:"shortDescription"`
		Highlights         string       `json:"highlights"`
		Features           string       `json:"features"`
		MetaTitle          string       `json:"metaTitle"`
		MetaDescription    string       `json:"metaDescription"`
		Tags               []StoreTag   `json:"tags"`
		Videos             []StoreVideo `json:"videos"`
		Faqs               []StoreFaq   `json:"faqs"`
	} `json:"infos"`
	Categories                          []StoreCategory   `json:"categories"`
	Category                            *StoreCategory    `json:"selectedFutureCategory"`
	Localizations                       []Locale          `json:"localizations"`
	AutomaticBugfixVersionCompatibility bool              `json:"automaticBugfixVersionCompatibility"`
	ProductType                         *StoreProductType `json:"productType"`
	Status                              struct {
		Name string `json:"name"`
	} `json:"status"`
	IconURL                                string `json:"iconUrl"`
	IsCompatibleWithLatestAllincartVersion bool   `json:"isCompatibleWithLatestAllincartVersion"`
}

type CreateExtensionRequest struct {
	Name       string `json:"name,omitempty"`
	SubType    string `json:"subType,omitempty"`
	ProducerID int    `json:"producerId"`
}

const (
	GenerationTheme  = "theme"
	GenerationApp    = "app"
	GenerationPlugin = "plugin"
)

func (e ProducerEndpoint) CreateExtension(ctx context.Context, newExtension CreateExtensionRequest) (*Extension, error) {
	requestBody, err := json.Marshal(newExtension)
	if err != nil {
		return nil, err
	}

	// Create it
	r, err := e.c.NewAuthenticatedRequest(ctx, "POST", fmt.Sprintf("%s/plugins", ApiUrl), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, err
	}

	var extension Extension
	if err := json.Unmarshal(body, &extension); err != nil {
		return nil, fmt.Errorf("create_extension: %v", err)
	}

	extension.Name = newExtension.Name

	// Patch the name
	err = e.UpdateExtension(ctx, &extension)

	if err != nil {
		return nil, err
	}

	return &extension, nil
}

func (e ProducerEndpoint) DeleteExtension(ctx context.Context, id int) error {
	r, err := e.c.NewAuthenticatedRequest(ctx, "DELETE", fmt.Sprintf("%s/plugins/%d", ApiUrl, id), nil)
	if err != nil {
		return err
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e ProducerEndpoint) UpdateExtension(ctx context.Context, extension *Extension) error {
	requestBody, err := json.Marshal(extension)

	if err != nil {
		return err
	}
	// Patch the name
	r, err := e.c.NewAuthenticatedRequest(ctx, "PUT", fmt.Sprintf("%s/plugins/%d", ApiUrl, extension.Id), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e ProducerEndpoint) GetSoftwareVersions(ctx context.Context, generation string) (*SoftwareVersionList, error) {
	errorFormat := "allincart_versions: %v"
	r, err := e.c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/pluginstatics/softwareVersions?filter=[{\"property\":\"pluginGeneration\",\"value\":\"%s\"},{\"property\":\"includeNonPublic\",\"value\":\"1\"}]", ApiUrl, generation), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	var versions SoftwareVersionList

	err = json.Unmarshal(body, &versions)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return &versions, nil
}

type SoftwareVersion struct {
	Name       string `json:"name"`
	Selectable bool   `json:"selectable"`
}

type Locale struct {
	Name string `json:"name"`
}

type StoreAvailablity struct {
	Name string `json:"name"`
}

type StoreCategory struct {
	Name string `json:"name"`
}

type StoreTag struct {
	Name string `json:"name"`
}

type StoreVideo struct {
	URL string `json:"url"`
}

type StoreProductType struct {
	Name string `json:"name"`
}

type StoreFaq struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type ExtensionGeneralInformation struct {
	FutureCategories []StoreCategory    `json:"futureCategories"`
	Locales          []Locale           `json:"locales"`
	ProductTypes     []StoreProductType `json:"productTypes"`
}

func (e ProducerEndpoint) GetExtensionGeneralInfo(ctx context.Context) (*ExtensionGeneralInformation, error) {
	r, err := e.c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/pluginstatics", ApiUrl), nil)
	if err != nil {
		return nil, fmt.Errorf("GetExtensionGeneralInfo: %v", err)
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf("GetExtensionGeneralInfo: %v", err)
	}

	var info *ExtensionGeneralInformation

	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, fmt.Errorf("allincart_versions: %v", err)
	}

	return info, nil
}
