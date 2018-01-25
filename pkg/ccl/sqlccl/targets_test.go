// Copyright 2016 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt

package sqlccl

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
)

func TestDescriptorsMatchingTargets(t *testing.T) {
	defer leaktest.AfterTest(t)()

	descriptors := []sqlbase.Descriptor{
		*sqlbase.WrapDescriptor(&sqlbase.DatabaseDescriptor{ID: 0, Name: "system"}),
		*sqlbase.WrapDescriptor(&sqlbase.TableDescriptor{ID: 1, Name: "foo", ParentID: 0}),
		*sqlbase.WrapDescriptor(&sqlbase.TableDescriptor{ID: 2, Name: "bar", ParentID: 0}),
		*sqlbase.WrapDescriptor(&sqlbase.DatabaseDescriptor{ID: 3, Name: "data"}),
		*sqlbase.WrapDescriptor(&sqlbase.TableDescriptor{ID: 4, Name: "baz", ParentID: 3}),
		*sqlbase.WrapDescriptor(&sqlbase.DatabaseDescriptor{ID: 5, Name: "empty"}),
	}

	tests := []struct {
		sessionDatabase string
		pattern         string
		expected        []string
		expectedDBs     []string
		err             string
	}{
		{"", "DATABASE system", []string{"system", "foo", "bar"}, []string{"system"}, ``},
		{"", "DATABASE system, noexist", nil, nil, `database "noexist" does not exist`},
		{"", "DATABASE data", []string{"data", "baz"}, []string{"data"}, ``},
		{"", "DATABASE system, data", []string{"system", "foo", "bar", "data", "baz"}, []string{"data", "system"}, ``},
		{"", "DATABASE system, data, noexist", nil, nil, `database "noexist" does not exist`},
		{"system", "DATABASE system", []string{"system", "foo", "bar"}, []string{"system"}, ``},
		{"system", "DATABASE system, noexist", nil, nil, `database "noexist" does not exist`},
		{"system", "DATABASE data", []string{"data", "baz"}, []string{"data"}, ``},
		{"system", "DATABASE system, data", []string{"system", "foo", "bar", "data", "baz"}, []string{"data", "system"}, ``},
		{"system", "DATABASE system, data, noexist", nil, nil, `database "noexist" does not exist`},

		{"", "TABLE foo", nil, nil, `table "foo" does not exist`},
		{"system", "TABLE foo", []string{"system", "foo"}, nil, ``},
		{"data", "TABLE foo", nil, nil, `table "foo" does not exist`},

		{"", "TABLE *", nil, nil, `no database specified for wildcard`},
		{"", "TABLE *, system.public.foo", nil, nil, `no database specified for wildcard`},
		{"noexist", "TABLE *", nil, nil, `database "noexist" does not exist`},
		{"system", "TABLE *", []string{"system", "foo", "bar"}, nil, ``},
		{"data", "TABLE *", []string{"data", "baz"}, nil, ``},
		{"empty", "TABLE *", []string{"empty"}, nil, ``},

		{"", "TABLE foo, baz", nil, nil, `table "(foo|baz)" does not exist`},
		{"system", "TABLE foo, baz", nil, nil, `table "baz" does not exist`},
		{"data", "TABLE foo, baz", nil, nil, `table "foo" does not exist`},

		{"", "TABLE system.public.foo", []string{"system", "foo"}, nil, ``},
		{"", "TABLE system.public.foo, foo", []string{"system", "foo"}, nil, `table "foo" does not exist`},

		{"", "TABLE system.public.foo, bar", []string{"system", "foo"}, nil, `table "bar" does not exist`},
		{"system", "TABLE system.public.foo, bar", []string{"system", "foo", "bar"}, nil, ``},

		{"", "TABLE noexist.*", nil, nil, `database "noexist" does not exist`},
		{"", "TABLE empty.*", []string{"empty"}, nil, ``},
		{"", "TABLE system.public.*", []string{"system", "foo", "bar"}, nil, ``},
		{"", "TABLE system.public.*, foo, baz", nil, nil, `table "(foo|baz)" does not exist`},
		{"system", "TABLE system.public.*, foo, baz", nil, nil, `table "baz" does not exist`},
		{"data", "TABLE system.public.*, baz", []string{"system", "foo", "bar", "data", "baz"}, nil, ``},
		{"data", "TABLE system.public.*, foo, baz", nil, nil, `table "(foo|baz)" does not exist`},

		{"", "TABLE system.public.FoO", []string{"system", "foo"}, nil, ``},

		{"", `TABLE system.public."foo"`, []string{"system", "foo"}, nil, ``},
		{"system", `TABLE "foo"`, []string{"system", "foo"}, nil, ``},
		// TODO(dan): Enable these tests once #8862 is fixed.
		// {"", `TABLE system."FOO"`, []string{"system"}},
		// {"system", `TABLE "FOO"`, []string{"system"}},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			sql := fmt.Sprintf(`GRANT ALL ON %s TO ignored`, test.pattern)
			stmt, err := parser.ParseOne(sql)
			if err != nil {
				t.Fatal(err)
			}
			targets := stmt.(*tree.Grant).Targets

			matched, err := descriptorsMatchingTargets(test.sessionDatabase, descriptors, targets)
			if test.err != "" {
				if !testutils.IsError(err, test.err) {
					t.Fatalf("expected error matching '%v', but got '%v'", test.err, err)
				}
			} else if err != nil {
				t.Fatal(err)
			} else {
				var matchedNames []string
				for _, m := range matched.descs {
					matchedNames = append(matchedNames, m.GetName())
				}
				var matchedDBNames []string
				for _, m := range matched.requestedDBs {
					matchedDBNames = append(matchedDBNames, m.GetName())
				}
				sort.Strings(test.expected)
				sort.Strings(test.expectedDBs)
				sort.Strings(matchedNames)
				sort.Strings(matchedDBNames)
				if !reflect.DeepEqual(test.expected, matchedNames) {
					t.Fatalf("expected %q got %q", test.expected, matchedNames)
				}
				if !reflect.DeepEqual(test.expectedDBs, matchedDBNames) {
					t.Fatalf("expected %q got %q", test.expectedDBs, matchedDBNames)
				}
			}
		})
	}
}
