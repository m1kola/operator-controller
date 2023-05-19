package conditionsets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	internalast "github.com/operator-framework/operator-controller/internal/ast"
	"github.com/operator-framework/operator-controller/internal/conditionsets"
)

const apiPath = "../../api/v1alpha1"

func TestConditionTypes(t *testing.T) {
	assert.NotEmpty(t, conditionsets.ConditionTypes)

	constNames, constValues, err := internalast.ParseTopLevelConstants(apiPath, "Type")
	assert.Nil(t, err)

	for i, condType := range constValues {
		assert.Containsf(t, conditionsets.ConditionTypes, condType, "%s is missing from ConditionTypes", constNames[i])
	}
	for _, condType := range conditionsets.ConditionTypes {
		assert.Containsf(t, constValues, condType, "There must be a Type%[1]s string literal constant for type %[1]q (i.e. 'const Type%[1]s = %[1]q')", condType)
	}
}

func TestConditionReasons(t *testing.T) {
	assert.NotEmpty(t, conditionsets.ConditionReasons)

	constNames, constValues, err := internalast.ParseTopLevelConstants(apiPath, "Reason")
	assert.Nil(t, err)

	for i, condReason := range constValues {
		assert.Containsf(t, conditionsets.ConditionReasons, condReason, "%s is missing from ConditionReasons", constNames[i])
	}
	for _, condReason := range conditionsets.ConditionReasons {
		assert.Containsf(t, constValues, condReason, "There must be a Reason%[1]s string literal constant for reason %[1]q (i.e. 'const Reason%[1]s = %[1]q')", condReason)
	}
}
