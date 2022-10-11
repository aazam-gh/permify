package services

import (
	"context"

	"github.com/Permify/permify/internal/commands"
	internalErrors "github.com/Permify/permify/internal/errors"
	"github.com/Permify/permify/internal/managers"
	"github.com/Permify/permify/pkg/dsl/schema"
	"github.com/Permify/permify/pkg/errors"
	base "github.com/Permify/permify/pkg/pb/base/v1"
)

// PermissionService -
type PermissionService struct {
	manager managers.IEntityConfigManager

	// commands
	check       commands.ICheckCommand
	expand      commands.IExpandCommand
	lookupQuery commands.ILookupQueryCommand
}

// NewPermissionService -
func NewPermissionService(cc commands.ICheckCommand, ec commands.IExpandCommand, lq commands.ILookupQueryCommand, en managers.IEntityConfigManager) *PermissionService {
	return &PermissionService{
		manager:     en,
		check:       cc,
		expand:      ec,
		lookupQuery: lq,
	}
}

// Check -
func (service *PermissionService) Check(ctx context.Context, subject *base.Subject, action string, entity *base.Entity, version string, d int32) (response commands.CheckResponse, err errors.Error) {
	var en *base.EntityDefinition
	en, err = service.manager.Read(ctx, entity.GetType(), version)
	if err != nil {
		return response, internalErrors.EntityTypeCannotFoundError
	}

	var a *base.ActionDefinition
	a, err = schema.GetAction(en, action)
	if err != nil {
		return response, internalErrors.ActionCannotFoundError
	}

	child := a.Child

	q := &commands.CheckQuery{
		Entity:  entity,
		Subject: subject,
	}

	q.SetDepth(d)

	return service.check.Execute(ctx, q, child)
}

// Expand -
func (service *PermissionService) Expand(ctx context.Context, entity *base.Entity, action string, version string) (response commands.ExpandResponse, err errors.Error) {
	var en *base.EntityDefinition
	en, err = service.manager.Read(ctx, entity.GetType(), version)
	if err != nil {
		return response, internalErrors.EntityTypeCannotFoundError
	}

	var a *base.ActionDefinition
	a, err = schema.GetAction(en, action)
	if err != nil {
		return response, internalErrors.ActionCannotFoundError
	}

	child := a.Child

	q := &commands.ExpandQuery{
		Entity: entity,
	}

	return service.expand.Execute(ctx, q, child)
}

// LookupQuery -
func (service *PermissionService) LookupQuery(ctx context.Context, entityType string, subject *base.Subject, action string, version string) (response commands.LookupQueryResponse, err errors.Error) {
	var sch *base.Schema
	sch, err = service.manager.All(ctx, version)
	if err != nil {
		return response, internalErrors.SchemaCannotFoundError
	}

	//entityType
	var en *base.EntityDefinition
	en, err = schema.GetEntityByName(sch, entityType)
	if err != nil {
		return response, internalErrors.EntityTypeCannotFoundError
	}

	var a *base.ActionDefinition
	a, err = schema.GetAction(en, action)
	if err != nil {
		return response, internalErrors.ActionCannotFoundError
	}

	child := a.Child

	q := &commands.LookupQueryQuery{
		EntityType: entityType,
		Action:     action,
		Subject:    subject,
	}

	q.SetSchema(sch)

	return service.lookupQuery.Execute(ctx, q, child)
}
