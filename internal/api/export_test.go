package api

import "net/url"

// BuildQueryForTest exposes the internal buildQuery for white-box tests.
func BuildQueryForTest(p ListParams) url.Values {
	return buildQuery(p)
}
