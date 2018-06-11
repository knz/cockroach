// Copyright 2018 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package parser

import (
	"fmt"
	"testing"

	"github.com/mjibson/pretty"

	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

func TestPretty(t *testing.T) {
	const S = `
		SELECT count(*) count, winner, counter * 60 * 5 as counter
			FROM (
				SELECT winner, round(length / 60 / 5) as counter
				FROM players
				WHERE build = $1 AND hero = $2
			)
			GROUP BY winner, counter
	`
	// SELECT s FROM config WHERE key = $1

	stmt, err := ParseOne(S)
	if err != nil {
		t.Fatal(err)
	}
	d := tree.Doc(stmt)
	for _, i := range []int{1, 9, 15, 20, 1000} {
		s, _ := pretty.PrettyString(d, i)
		fmt.Printf("%d: %s\n\n", i, s)
	}
}
