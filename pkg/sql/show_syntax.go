// Copyright 2017 The Cockroach Authors.
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

package sql

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/types"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
)

var showSyntaxColumns = sqlbase.ResultColumns{
	{Name: "field", Typ: types.String},
	{Name: "text", Typ: types.String},
}

type showSyntaxOpts struct {
	// autoWrap, when true, uses pretty.Pretty to render the statement.
	// When false the statement is rendered on a single line and the
	// other options are ignored.
	autoWrap bool

	// The following are options for pretty.Pretty.

	// lineWidth is the desired maximum line length.
	lineWidth int
	// simplify indicates whether to remove extraneous parentheses.
	simplify bool
	// useTabs indicates whether to use tabs.
	useTabs bool
	// indentWidth indicates how much spaces to use for indentation,
	// and, if useTabs is set, the width of a tab in spaces.
	indentWidth int
}

var defaultShowSyntaxOpts = showSyntaxOpts{
	autoWrap:  true,
	lineWidth: 60,
	// do not simplify by default because SHOW SYNTAX is typically used
	// to reproduce user input with as minimal changes as possible.
	simplify: false,
	// do not use tabs by default because SQL query results are usually
	// rendered in a context that can mess up tabs (e.g. ASCII art
	// tables in a terminal, or in GUIs).
	useTabs:     false,
	indentWidth: 4,
}

// ShowSyntax implements the plan for SHOW SYNTAX. This statement is
// usually handled as a special case in Executor, but for
// FROM [SHOW SYNTAX ...] we will arrive here too.
func (p *planner) ShowSyntax(ctx context.Context, n *tree.ShowSyntax) (planNode, error) {
	return newShowSyntaxPlan(n.Statement,
		func() (showSyntaxOpts, error) { return defaultShowSyntaxOpts, nil }), nil
}

// newShowSyntaxNode constructs a hookFnNode to run the SHOW SYNTAX
// statement for both the simple ShowSyntax node and
// ShowSyntaxExtended (SHOW SYNTAX ... WITH).
//
// The evalOpts argument is a deferred evaluator for the statement's
// options. The evaluation is deferred because it may contain
// references to placeholders or subqueries that cannot be known
// during prepare.
func newShowSyntaxPlan(stmt string, evalOpts func() (showSyntaxOpts, error)) planNode {
	// TODO(knz): in the call to runShowSyntax() below, reportErr is nil
	// although we might want to be able to capture (and report) these
	// errors as well.
	//
	// However, this code path is only used when SHOW SYNTAX is used as
	// a data source, i.e. a client actively uses a query of the form
	// SELECT ... FROM [SHOW SYNTAX ' ... '] WHERE ....  This is not
	// what `cockroach sql` does: the SQL shell issues a straight `SHOW
	// SYNTAX` that goes through the "statement observer" code
	// path. Since we care mainly about what users do in the SQL shell,
	// it's OK if we only deal with that case well for now and, for the
	// time being, forget/ignore errors when SHOW SYNTAX is used as data
	// source. This can be added later if deemed useful or necessary.
	fn := func(ctx context.Context, _ []planNode, res chan<- tree.Datums) error {
		opts, err := evalOpts()
		if err != nil {
			return err
		}
		return runShowSyntax(ctx, stmt, opts,
			func(ctx context.Context, field, msg string) error {
				res <- tree.Datums{tree.NewDString(field), tree.NewDString(msg)}
				return nil
			},
			nil /* reportErr */)
	}

	return &hookFnNode{header: showSyntaxColumns, f: fn}
}

// runShowSyntax analyzes the syntax and reports its structure as data
// for the client. Even an error is reported as data.
//
// Since errors won't propagate to the client as an error, but as
// a result, the usual code path to capture and record errors will not
// be triggered. Instead, the caller can pass a reportErr closure to
// capture errors instead. May be nil.
func runShowSyntax(
	ctx context.Context,
	stmt string,
	opts showSyntaxOpts,
	report func(ctx context.Context, field, msg string) error,
	reportErr func(err error),
) error {
	stmts, err := parser.Parse(stmt)
	if err != nil {
		if reportErr != nil {
			reportErr(err)
		}

		pqErr, ok := pgerror.GetPGCause(err)
		if !ok {
			return pgerror.NewErrorf(pgerror.CodeInternalError, "unknown parser error: %v", err)
		}
		if err := report(ctx, "error", pqErr.Message); err != nil {
			return err
		}
		if err := report(ctx, "code", pqErr.Code); err != nil {
			return err
		}
		if pqErr.Source != nil {
			if pqErr.Source.File != "" {
				if err := report(ctx, "file", pqErr.Source.File); err != nil {
					return err
				}
			}
			if pqErr.Source.Line > 0 {
				if err := report(ctx, "line", fmt.Sprintf("%d", pqErr.Source.Line)); err != nil {
					return err
				}
			}
			if pqErr.Source.Function != "" {
				if err := report(ctx, "function", pqErr.Source.Function); err != nil {
					return err
				}
			}
		}
		if pqErr.Detail != "" {
			if err := report(ctx, "detail", pqErr.Detail); err != nil {
				return err
			}
		}
		if pqErr.Hint != "" {
			if err := report(ctx, "hint", pqErr.Hint); err != nil {
				return err
			}
		}
	} else {
		for _, stmt := range stmts {
			var s string
			if opts.autoWrap {
				// We do not use the usual value for fmtFlags, FmtRoundtrip
				// which includes FmtShowPasswords, because we don't want
				// passwords to show up in user logs.
				s = tree.PrettyWithOpts(stmt,
					opts.lineWidth, opts.useTabs, opts.indentWidth, opts.simplify, tree.FmtParsable)
			} else {
				s = tree.AsStringWithFlags(stmt, tree.FmtParsable)
			}
			if err := report(ctx, "sql", s); err != nil {
				return err
			}
		}
	}
	return nil
}

// ShowSyntaxExtended implements the plan for SHOW SYNTAX WITH.
func (p *planner) ShowSyntaxExtended(
	ctx context.Context, n *tree.ShowSyntaxExtended,
) (planNode, error) {
	var autoWrapOpt, lineWidthOpt, indentWidthOpt, simplifyOpt, useTabsOpt func() (tree.Datum, error)

	var err error
	for _, opt := range n.Options {
		switch opt.Key {
		case "auto_wrap":
			autoWrapOpt, err = p.TypeAs(opt.Value, types.Bool, string(opt.Key))
			if err != nil {
				return nil, err
			}
		case "width":
			lineWidthOpt, err = p.TypeAs(opt.Value, types.Int, string(opt.Key))
			if err != nil {
				return nil, err
			}
		case "indent":
			indentWidthOpt, err = p.TypeAs(opt.Value, types.Int, string(opt.Key))
			if err != nil {
				return nil, err
			}
		case "simplify":
			simplifyOpt, err = p.TypeAs(opt.Value, types.Bool, string(opt.Key))
			if err != nil {
				return nil, err
			}
		case "tabs":
			useTabsOpt, err = p.TypeAs(opt.Value, types.Bool, string(opt.Key))
			if err != nil {
				return nil, err
			}
		default:
			return nil, pgerror.NewErrorf(pgerror.CodeUndefinedParameterError,
				"unknown option: %q", opt.Key)
		}
	}

	optFn := func() (showSyntaxOpts, error) {
		opts := defaultShowSyntaxOpts
		if autoWrapOpt != nil {
			d, err := autoWrapOpt()
			if err != nil {
				return opts, err
			}
			opts.autoWrap = bool(*(d.(*tree.DBool)))
		}
		if lineWidthOpt != nil {
			d, err := lineWidthOpt()
			if err != nil {
				return opts, err
			}
			opts.lineWidth = int(*(d.(*tree.DInt)))
		}
		if indentWidthOpt != nil {
			d, err := indentWidthOpt()
			if err != nil {
				return opts, err
			}
			opts.indentWidth = int(*(d.(*tree.DInt)))
		}
		if simplifyOpt != nil {
			d, err := simplifyOpt()
			if err != nil {
				return opts, err
			}
			opts.simplify = bool(*(d.(*tree.DBool)))
		}
		if useTabsOpt != nil {
			d, err := useTabsOpt()
			if err != nil {
				return opts, err
			}
			opts.useTabs = bool(*(d.(*tree.DBool)))
		}
		return opts, nil
	}

	return newShowSyntaxPlan(n.ShowSyntax.Statement, optFn), nil
}
