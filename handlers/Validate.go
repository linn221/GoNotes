package handlers

import (
	"linn221/shop/utils"

	"gorm.io/gorm"
)

type Rule interface {
	Init() bool
	CountResults(*gorm.DB, *int64) *ServiceError
}

func Validate(db *gorm.DB, rules ...Rule) *ServiceError {
	var count int64
	for _, rule := range rules {
		if ok := rule.Init(); !ok {
			continue
		}
		err := rule.CountResults(db, &count)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateInBatch(db *gorm.DB, rules ...Rule) []*ServiceError {
	var count int64
	errors := make([]*ServiceError, 0)
	for _, rule := range rules {
		if ok := rule.Init(); !ok {
			continue
		}
		err := rule.CountResults(db, &count)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

type HasFilter struct {
	Cond         string
	FilterValues []interface{}
}

func (f *HasFilter) ApplyFilter(dbCtx *gorm.DB) {
	if f != nil {
		dbCtx.Where(f.Cond, f.FilterValues...)
	}
}

func NewFilter(cond string, values ...interface{}) *HasFilter {
	return &HasFilter{
		Cond:         cond,
		FilterValues: values,
	}
}

func NewShopFilter(shopId string) *HasFilter {
	return &HasFilter{
		Cond:         "shop_id = ?",
		FilterValues: []interface{}{shopId},
	}
}

// check if resource exists (where business_id = ?)
type ruleExists struct {
	statusCode *int
	table      string
	id         interface{}
	message    string
	do         *bool
	*HasFilter
}

// specifies When to validate
// if When is not specified, will validate by default
func (rule ruleExists) When(when bool) ruleExists {
	rule.do = &when
	return rule
}

func (rule ruleExists) OverrideStatusCode(i int) ruleExists {
	rule.statusCode = &i
	return rule
}

func (vr ruleExists) Init() bool {
	// skip validation if user specifies when
	return vr.do == nil || *vr.do
}

func (vr ruleExists) CountResults(dbCtx *gorm.DB, count *int64) *ServiceError {
	dbCtx = dbCtx.Table(vr.table).Where("id = ?", vr.id)
	vr.ApplyFilter(dbCtx)
	if err := dbCtx.Count(count).Error; err != nil {
		return systemErr(err)
	}
	if *count <= 0 {
		return clientErr(vr.message)
	}

	return nil
}

func NewExistsRule(table string, id interface{}, message string, filter *HasFilter) ruleExists {
	return ruleExists{
		table:     table,
		id:        id,
		message:   message,
		HasFilter: filter,
	}
}

// check if slice of resource id exists (where business_id IN ?)
type RuleMassExists[ID comparable] struct {
	Table         string
	Ids           []ID
	Message       string
	NoDuplicateID bool
	*HasFilter
}

func (r RuleMassExists[ID]) Init() bool {
	return len(r.Ids) > 0
}

func (r RuleMassExists[ID]) CountResults(dbCtx *gorm.DB, count *int64) *ServiceError {

	if r.NoDuplicateID {
		ids, duplicates := utils.UniqueSliceWithDuplicateCount(r.Ids)
		if duplicates > 0 {
			return clientErr("duplicate ids for " + r.Table)
		}
		dbCtx = dbCtx.Table(r.Table).Where("id IN ?", ids)
		err := dbCtx.Count(count).Error
		if err != nil {
			return systemErr(err)
		}
		if *count != int64(len(ids)) {
			return clientErr(r.Message)
		}

	} else {
		uniqIds := utils.UniqueSlice(r.Ids)
		dbCtx = dbCtx.Table(r.Table).Where("id IN ?", uniqIds)
		err := dbCtx.Count(count).Error
		if err != nil {
			return systemErr(err)
		}
		if *count != int64(len(uniqIds)) {
			return clientErr(r.Message)
		}

	}
	return nil
}

type ruleUnique struct {
	table    string
	message  string
	column   string
	value    interface{}
	exceptId int
	do       *bool

	*HasFilter
}

func (rule ruleUnique) When(cond bool) ruleUnique {
	rule.do = &cond
	return rule
}

func (rule ruleUnique) Filter(cond string, values ...interface{}) ruleUnique {
	rule.HasFilter = &HasFilter{
		Cond:         cond,
		FilterValues: values,
	}
	return rule
}

func (rule ruleUnique) Init() bool {

	if rule.do != nil && !*rule.do {
		return false
	}
	return true
}

func (r ruleUnique) CountResults(dbCtx *gorm.DB, count *int64) *ServiceError {
	dbCtx = dbCtx.Table(r.table).Where("`"+r.column+"`"+" = ?", r.value)
	if r.exceptId > 0 {
		dbCtx.Where("id != ?", r.exceptId)
	}
	r.ApplyFilter(dbCtx)
	err := dbCtx.Count(count).Error
	if err != nil {
		return systemErr(err)
	}

	if *count > 0 {
		return clientErr(r.message)
	}

	return nil
}

func NewUniqueRule(table string, column string, value interface{}, exceptId int, message string, filter *HasFilter) ruleUnique {
	// var v T
	return ruleUnique{
		table:     table,
		column:    column,
		value:     value,
		exceptId:  exceptId,
		message:   message,
		HasFilter: filter,
	}
}

type noResultRule struct {
	statusCode *int
	table      string
	message    string
	do         *bool
	*HasFilter
}

// specifies When to validate
// if When is not specified, will validate by default
func (rule noResultRule) When(when bool) noResultRule {
	rule.do = &when
	return rule
}

func (rule noResultRule) OverrideStatusCode(i int) noResultRule {
	rule.statusCode = &i
	return rule
}

func (vr noResultRule) Init() bool {
	// skip validation if user specifies when
	return vr.do == nil || *vr.do
}

func (vr noResultRule) CountResults(dbCtx *gorm.DB, count *int64) *ServiceError {
	dbCtx = dbCtx.Table(vr.table)
	vr.ApplyFilter(dbCtx)
	if err := dbCtx.Count(count).Error; err != nil {
		return systemErr(err)
	}
	if *count > 0 {
		return clientErr(vr.message)
	}

	return nil
}

func NewNoResultRule(table string, message string, filter *HasFilter) noResultRule {
	return noResultRule{
		table:     table,
		message:   message,
		HasFilter: filter,
	}
}
