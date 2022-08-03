package casbinpgadapter

import (
	"database/sql"
	"log"

	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/nrfta/go-casbin-pg-adapter/db/migrations"
	// no-lint
	_ "github.com/lib/pq"

	"github.com/nrfta/go-casbin-pg-adapter/pkg/model"
	"github.com/nrfta/go-casbin-pg-adapter/pkg/repository"
)

// Adapter is a postgresql adaptor for casbin
type Adapter struct {
	db                   *sql.DB
	dbSchema             string
	tableName            string
	casbinRuleRepository *repository.CasbinRuleRepository
}

var _ = persist.Adapter(&Adapter{})
var _ = persist.BatchAdapter(&Adapter{})

// NewAdapter returns a new casbin postgresql adapter
func NewAdapter(db *sql.DB, tableName string) (*Adapter, error) {
	return NewAdapterWithDBSchema(db, "public", tableName)
}

// NewAdapterWithDBSchema returns a new casbin postgresql adapter with the schema named dbSchema
func NewAdapterWithDBSchema(db *sql.DB, dbSchema string, tableName string) (*Adapter, error) {
	casbinRuleRepository := repository.NewCasbinRuleRepository(dbSchema, tableName, db)
	adapter := &Adapter{
		db,
		dbSchema,
		tableName,
		casbinRuleRepository,
	}

	if err := migrations.Migrate(adapter.dbSchema, adapter.tableName, adapter.db); err != nil {
		log.Println("casbin pg migrations filed:", err)
		return nil, err
	}

	return adapter, nil
}

// LoadPolicy loads all policy rules from the storage.
func (adapter *Adapter) LoadPolicy(cmodel casbinModel.Model) error {
	casbinRules, err := adapter.casbinRuleRepository.LoadAllCasbinRules()
	if err != nil {
		return err
	}

	for _, casbinRule := range casbinRules {
		persist.LoadPolicyLine(casbinRule.ToPolicyLine(), cmodel)
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (adapter *Adapter) SavePolicy(cmodel casbinModel.Model) error {
	casbinRules := make([]model.CasbinRule, 0)
	for pType, ast := range cmodel["p"] {
		for _, rule := range ast.Policy {
			casbinRule := model.NewCasbinRuleFromPTypeAndRule(pType, rule)
			casbinRules = append(casbinRules, casbinRule)
		}
	}
	for pType, ast := range cmodel["g"] {
		for _, rule := range ast.Policy {
			casbinRule := model.NewCasbinRuleFromPTypeAndRule(pType, rule)
			casbinRules = append(casbinRules, casbinRule)
		}
	}
	if err := adapter.casbinRuleRepository.ReplaceAllCasbinRules(casbinRules); err != nil {
		return err
	}
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (adapter *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	casbinRule := model.NewCasbinRuleFromPTypeAndRule(ptype, rule)
	err := adapter.casbinRuleRepository.InsertCasbinRule(casbinRule)
	return err
}

// AddPolicies adds policy rules to the storage.
// This is part of the Auto-Save feature.
func (adapter *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	for _, rule := range rules {
		if err := adapter.AddPolicy(sec, ptype, rule); err != nil {
			return err
		}
	}
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (adapter *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	casbinRule := model.NewCasbinRuleFromPTypeAndRule(ptype, rule)
	err := adapter.casbinRuleRepository.DeleteCasbinRule(casbinRule)
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (adapter *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	casbinRule := model.NewCasbinRuleFromPTypeAndFilter(ptype, fieldIndex, fieldValues...)
	err := adapter.casbinRuleRepository.DeleteCasbinRule(casbinRule)
	return err
}

// RemovePolicies removes policy rules from the storage.
// This is part of the Auto-Save feature.
func (adapter *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	for _, rule := range rules {
		if err := adapter.RemovePolicy(sec, ptype, rule); err != nil {
			return err
		}
	}
	return nil
}
