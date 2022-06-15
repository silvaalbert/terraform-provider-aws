package kendra_test

import "testing"

// Serialize to limit service quota exceeded errors.
func TestAccKendra_serial(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"Thesaurus": {
			"basic":        testAccThesaurus_basic,
			"disappears":   testAccThesaurus_disappears,
			"tags":         testAccThesaurus_tags,
			"Description":  testAccThesaurus_Description,
			"Name":         testAccThesaurus_Name,
			"RoleARN":      testAccThesaurus_RoleARN,
			"SourceS3Path": testAccThesaurus_SourceS3Path,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}
