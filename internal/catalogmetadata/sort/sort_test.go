package sort_test

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/operator-framework/operator-controller/internal/catalogmetadata"
	metadatasort "github.com/operator-framework/operator-controller/internal/catalogmetadata/sort"
	"github.com/operator-framework/operator-registry/alpha/declcfg"
	"github.com/operator-framework/operator-registry/alpha/property"
	"github.com/stretchr/testify/assert"
)

func TestSortByPackageName(t *testing.T) {
	b1 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-a", "version": "1.0.0"}`)},
	}}}
	b2 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-b", "version": "1.0.0"}`)},
	}}}
	b3 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-c", "version": "1.0.0"}`)},
	}}}

	bundles := []*catalogmetadata.Bundle{&b2, &b3, &b1}
	sort.Slice(bundles, func(i, j int) bool {
		return metadatasort.ByPackageAndVersion(bundles[i], bundles[j])
	})
	assert.Equal(t, bundles[0], &b1)
	assert.Equal(t, bundles[1], &b2)
	assert.Equal(t, bundles[2], &b3)
}

func TestSortByVersion(t *testing.T) {
	b1 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-a", "version": "1.0.0"}`)},
		{Type: property.TypeChannel, Value: json.RawMessage(`{"channelName":"alpha","priority":0}`)},
	}}}
	b2 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-a", "version": "1.0.1"}`)},
		{Type: property.TypeChannel, Value: json.RawMessage(`{"channelName":"alpha","priority":0}`)},
	}}}
	b3 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-a", "version": "2.0.0"}`)},
		{Type: property.TypeChannel, Value: json.RawMessage(`{"channelName":"alpha","priority":0}`)},
	}}}

	bundles := []*catalogmetadata.Bundle{&b2, &b3, &b1}
	sort.Slice(bundles, func(i, j int) bool {
		return metadatasort.ByPackageAndVersion(bundles[i], bundles[j])
	})
	assert.Equal(t, bundles[0], &b3)
	assert.Equal(t, bundles[1], &b2)
	assert.Equal(t, bundles[2], &b1)
}

func TestSortWithMissingProperties(t *testing.T) {
	b1 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Name: "b1", Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-a", "version": "1.0.0"}`)},
	}}}
	b2 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Name: "b2", Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-a"}`)},
	}}}
	b3 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Name: "b3", Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-a", "version": "2.0.0"}`)},
	}}}
	b4 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Name: "b4", Properties: []property.Property{
		{Type: property.TypePackage, Value: json.RawMessage(`{"packageName": "package-b", "version": "2.0.0"}`)},
	}}}
	b5 := catalogmetadata.Bundle{Bundle: declcfg.Bundle{Name: "b5", Properties: []property.Property{}}}

	bundles := []*catalogmetadata.Bundle{&b2, &b3, &b1, &b4, &b5}
	sort.Slice(bundles, func(i, j int) bool {
		return metadatasort.ByPackageAndVersion(bundles[i], bundles[j])
	})
	assert.Equal(t, bundles[0], &b3) // alphabetically highest with highest version
	assert.Equal(t, bundles[1], &b1) // alphabetically highest with lowest version
	assert.Equal(t, bundles[2], &b4) // alphabetically lowest
	assert.Equal(t, bundles[3], &b2) // no version; missing fields get sorted last
	assert.Equal(t, bundles[4], &b5) // empty; missing fields get sorted last
}
