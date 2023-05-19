package conditionsets

//go:generate go run ../../hack/codegen/conditionsets -apiPath ../../api/v1alpha1 -apiImportAlias operatorsv1alpha1 -apiImport github.com/operator-framework/operator-controller/api/v1alpha1 -fromPrefixToSlice Type:ConditionTypes -fromPrefixToSlice Reason:ConditionReasons -targetFilePath ./zz_generated.conditionsets.go -targetPackageName $GOPACKAGE
