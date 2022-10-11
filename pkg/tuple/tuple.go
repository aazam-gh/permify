package tuple

import (
	"fmt"
	"strings"

	"github.com/Permify/permify/pkg/errors"
	"github.com/Permify/permify/pkg/helper"
	base "github.com/Permify/permify/pkg/pb/base/v1"
)

const (
	ENTITY   = "%s:%s"
	RELATION = "#%s"
)

const (
	ELLIPSIS = "..."
)

const (
	USER = "user"
)

const (
	SEPARATOR = "."
)

// IsSubjectUser -
func IsSubjectUser(subject *base.Subject) bool {
	if subject.Type == USER {
		return true
	}
	return false
}

// AreSubjectsEqual -
func AreSubjectsEqual(s1 *base.Subject, s2 *base.Subject) bool {
	if IsSubjectUser(s1) {
		return s1.GetId() == s2.GetId()
	}
	return s1.GetRelation() == s2.GetRelation() && s1.GetId() == s2.GetId() && s1.GetType() == s1.GetType()
}

// EntityAndRelationToString -
func EntityAndRelationToString(entityAndRelation *base.EntityAndRelation) string {
	return EntityToString(entityAndRelation.GetEntity()) + fmt.Sprintf(RELATION, entityAndRelation.GetRelation())
}

// EntityToString -
func EntityToString(entity *base.Entity) string {
	return fmt.Sprintf(ENTITY, entity.GetType(), entity.GetId())
}

// SubjectToString -
func SubjectToString(subject *base.Subject) string {
	if IsSubjectUser(subject) {
		return fmt.Sprintf(ENTITY, subject.GetType(), subject.GetId())
	}
	return fmt.Sprintf("%s"+RELATION, fmt.Sprintf(ENTITY, subject.GetType(), subject.GetId()), subject.GetRelation())
}

// IsEntityAndSubjectEquals -
func IsEntityAndSubjectEquals(t *base.Tuple) bool {
	return t.GetEntity().GetType() == t.GetSubject().GetType() && t.GetEntity().GetType() == t.GetSubject().GetType() && t.GetRelation() == t.GetSubject().GetRelation()
}

// ValidateSubjectType -
func ValidateSubjectType(subject *base.Subject, relationTypes []string) (err errors.Error) {
	key := subject.GetType()
	if subject.GetRelation() != "" {
		if !IsSubjectUser(subject) {
			if subject.GetRelation() != ELLIPSIS {
				key += "#" + subject.GetRelation()
			}
		}
	}

	if !helper.InArray(key, relationTypes) {
		return errors.NewError(errors.Validation).SetMessage("subject type is not found in defined types")
	}
	return nil
}

// SplitRelation -
func SplitRelation(relation string) (a []string) {
	s := strings.Split(relation, SEPARATOR)
	for _, b := range s {
		a = append(a, b)
	}
	if len(a) == 1 {
		a = append(a, "")
	}
	return
}

// IsRelationComputed -
func IsRelationComputed(relation string) bool {
	sp := strings.Split(relation, SEPARATOR)
	if len(sp) == 1 {
		return true
	}
	return false
}

// IsSubjectValid -
func IsSubjectValid(subject *base.Subject) bool {
	if subject.GetType() == "" {
		return false
	}

	if subject.GetId() == "" {
		return false
	}

	if IsSubjectUser(subject) {
		return subject.GetRelation() == ""
	} else {
		return subject.GetRelation() != ""
	}
}
