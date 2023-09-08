package sort

import (
	"strings"

	"github.com/operator-framework/operator-controller/internal/catalogmetadata"
)

// ByPackageAndVersion is a sort function that orders bundles by package
// and inverse version (higher versions on top).
// If a property does not exist for one of the entities, the one missing the property
// is pushed down; if both entities are missing the same property they are ordered by id.
func ByPackageAndVersion(b1, b2 *catalogmetadata.Bundle) bool {
	// first sort package lexical order
	pkgOrder := packageOrder(b1, b2)
	if pkgOrder != 0 {
		return pkgOrder < 0
	}

	// order version from highest to lowest (favor the latest release)
	versionOrder := versionOrder(b1, b2)
	return versionOrder > 0
}

func compareErrors(err1 error, err2 error) int {
	if err1 != nil && err2 == nil {
		return 1
	}

	if err1 == nil && err2 != nil {
		return -1
	}
	return 0
}

func packageOrder(b1, b2 *catalogmetadata.Bundle) int {
	name1, err1 := b1.PackageName()
	name2, err2 := b2.PackageName()
	errComp := compareErrors(err1, err2)
	if errComp != 0 {
		return errComp
	}
	return strings.Compare(name1, name2)
}

func versionOrder(b1, b2 *catalogmetadata.Bundle) int {
	ver1, err1 := b1.Version()
	ver2, err2 := b2.Version()
	errComp := compareErrors(err1, err2)
	if errComp != 0 {
		// the sign gets inverted because version is sorted
		// from highest to lowest
		return -1 * errComp
	} else if ver1 == nil && ver2 == nil {
		// if neither bundle has version they are considered equal
		return 0
	}
	return ver1.Compare(*ver2)
}
