package ror

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// OrganizationQueryBuilder helps build queries for searching organizations.
type OrganizationQueryBuilder struct {
	client *Client

	pageSize int

	organizationContinentCodeList []string
	organizationContinentNameList []string
	organizationCountryCodeList   []string
	organizationCountryNameList   []string
	statusList                    []OrganizationStatus
	typeList                      []OrganizationType

	createdDateFrom   *time.Time
	createdDateUntil  *time.Time
	modifiedDateFrom  *time.Time
	modifiedDateUntil *time.Time

	numberOfResults *int
	query           *string
}

// NewOrganizationQueryBuilder creates a new OrganizationQueryBuilder.
func NewOrganizationQueryBuilder(client *Client) *OrganizationQueryBuilder {
	return &OrganizationQueryBuilder{
		client:   client,
		pageSize: 20,
	}
}

// WithStatus adds a status filter.
func (b *OrganizationQueryBuilder) WithStatus(status OrganizationStatus) *OrganizationQueryBuilder {
	b.statusList = append(b.statusList, status)
	return b
}

// WithType adds a type filter.
func (b *OrganizationQueryBuilder) WithType(t OrganizationType) *OrganizationQueryBuilder {
	b.typeList = append(b.typeList, t)
	return b
}

// WithCountryCode adds a country code filter.
func (b *OrganizationQueryBuilder) WithCountryCode(countryCode string) *OrganizationQueryBuilder {
	b.organizationCountryCodeList = append(b.organizationCountryCodeList, countryCode)
	return b
}

// WithCountryName adds a country name filter.
func (b *OrganizationQueryBuilder) WithCountryName(countryName string) *OrganizationQueryBuilder {
	b.organizationCountryNameList = append(b.organizationCountryNameList, countryName)
	return b
}

// WithContinentCode adds a continent code filter.
func (b *OrganizationQueryBuilder) WithContinentCode(continentCode string) *OrganizationQueryBuilder {
	b.organizationContinentCodeList = append(b.organizationContinentCodeList, continentCode)
	return b
}

// WithContinentName adds a continent name filter.
func (b *OrganizationQueryBuilder) WithContinentName(continentName string) *OrganizationQueryBuilder {
	b.organizationContinentNameList = append(b.organizationContinentNameList, continentName)
	return b
}

// CreatedDateFrom sets the start date for creation filter.
func (b *OrganizationQueryBuilder) CreatedDateFrom(createdDateFrom time.Time) *OrganizationQueryBuilder {
	b.createdDateFrom = &createdDateFrom
	return b
}

// CreatedDateUntil sets the end date for creation filter.
func (b *OrganizationQueryBuilder) CreatedDateUntil(createdDateUntil time.Time) *OrganizationQueryBuilder {
	b.createdDateUntil = &createdDateUntil
	return b
}

// ModifiedDateFrom sets the start date for modified filter.
func (b *OrganizationQueryBuilder) ModifiedDateFrom(modifiedDateFrom time.Time) *OrganizationQueryBuilder {
	b.modifiedDateFrom = &modifiedDateFrom
	return b
}

// ModifiedDateUntil sets the end date for modified filter.
func (b *OrganizationQueryBuilder) ModifiedDateUntil(modifiedDateUntil time.Time) *OrganizationQueryBuilder {
	b.modifiedDateUntil = &modifiedDateUntil
	return b
}

// WithQuery sets the search query.
func (b *OrganizationQueryBuilder) WithQuery(query string) *OrganizationQueryBuilder {
	b.query = &query
	return b
}

// WithNumberOfResults sets the max number of results.
func (b *OrganizationQueryBuilder) WithNumberOfResults(numberOfResults int) *OrganizationQueryBuilder {
	b.numberOfResults = &numberOfResults
	return b
}

// Execute performs the query.
func (b *OrganizationQueryBuilder) Execute(ctx context.Context) (*OrganizationsResult, error) {
	queries, err := b.BuildQueries()
	if err != nil {
		return nil, err
	}

	var results []*OrganizationsResult
	for _, q := range queries {
		result, err := b.client.PerformQuery(ctx, q)
		if err != nil {
			return nil, err
		}
		if result != nil {
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		return nil, nil
	}

	first := results[0]
	for i := 1; i < len(results); i++ {
		first = first.Combine(results[i])
	}
	return first, nil
}

// BuildQueries builds the query strings for all pages.
func (b *OrganizationQueryBuilder) BuildQueries() ([]string, error) {
	results := b.pageSize
	if b.numberOfResults != nil {
		results = *b.numberOfResults
		if results <= 0 {
			return nil, errors.New("number of results must be greater than 0")
		}
	}

	if results <= b.pageSize {
		return []string{b.buildQuery(nil)}, nil
	}

	pages := results / b.pageSize
	if results%b.pageSize > 0 {
		pages++
	}

	var queries []string
	for i := 1; i <= pages; i++ {
		page := i
		queries = append(queries, b.buildQuery(&page))
	}
	return queries, nil
}

func (b *OrganizationQueryBuilder) buildQuery(page *int) string {
	var components []string

	if page != nil {
		components = append(components, fmt.Sprintf("page=%d", *page))
	}

	var filters []string
	for _, s := range b.statusList {
		filters = append(filters, "status:"+string(s))
	}
	for _, t := range b.typeList {
		filters = append(filters, "types:"+string(t))
	}
	for _, cc := range b.organizationCountryCodeList {
		filters = append(filters, "country.country_code:"+url.QueryEscape(cc))
	}
	for _, cn := range b.organizationCountryNameList {
		filters = append(filters, "locations.geonames_details.country_name:"+url.QueryEscape(cn))
	}
	for _, cc := range b.organizationContinentCodeList {
		filters = append(filters, "locations.geonames_details.continent_code:"+url.QueryEscape(cc))
	}
	for _, cn := range b.organizationContinentNameList {
		filters = append(filters, "locations.geonames_details.continent_name:"+url.QueryEscape(cn))
	}

	if len(filters) > 0 {
		components = append(components, "filter="+strings.Join(filters, ","))
	}

	var advancedQuery []string
	if b.createdDateFrom != nil || b.createdDateUntil != nil {
		advancedQuery = append(advancedQuery, "admin.created.date:"+getFormattedDateRange(b.createdDateFrom, b.createdDateUntil))
	}
	if b.modifiedDateFrom != nil || b.modifiedDateUntil != nil {
		advancedQuery = append(advancedQuery, "admin.last_modified.date:"+getFormattedDateRange(b.modifiedDateFrom, b.modifiedDateUntil))
	}

	if b.query != nil && len(advancedQuery) > 0 {
		advancedQuery = append(advancedQuery, url.QueryEscape(*b.query))
	} else if b.query != nil {
		components = append(components, "query="+url.QueryEscape(*b.query))
	}

	if len(advancedQuery) > 0 {
		components = append(components, "query.advanced="+strings.Join(advancedQuery, "%20AND%20"))
	}

	return strings.Join(components, "&")
}

func getFormattedDateRange(from *time.Time, until *time.Time) string {
	dateStr := func(t time.Time) string {
		return t.Format("2006-01-02")
	}

	f := ""
	if from != nil {
		f = dateStr(*from)
	}
	u := ""
	if until != nil {
		u = dateStr(*until)
	}

	if f == "" {
		f = "0001-01-01"
	}
	if u == "" {
		u = "9999-12-31"
	}

	return fmt.Sprintf("%%5B%s%%20TO%20%s%%5D", f, u)
}
