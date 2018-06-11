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

package tree

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/mjibson/sqlfmt/pretty"
)

type Docer interface {
	Doc() pretty.Doc
}

func Doc(s Statement) pretty.Doc {
	if d, ok := s.(Docer); ok {
		return d.Doc()
	}
	spew.Fdump(os.Stderr, "DOCSTATEMENT", s)
	return pretty.Text(s.String())
}

func DocExpr(node Expr) pretty.Doc {
	if d, ok := node.(Docer); ok {
		return d.Doc()
	}
	spew.Fdump(os.Stderr, "DOCEXPR", node)
	return pretty.Text(node.String())
}

func DocNode(f NodeFormatter) pretty.Doc {
	if d, ok := f.(Docer); ok {
		return d.Doc()
	}
	spew.Fdump(os.Stderr, "DOCNODE", f)
	return pretty.Text(AsString(f))
}

func (s SelectExprs) Doc() pretty.Doc {
	d := make([]pretty.Doc, len(s))
	for i, e := range s {
		d[i] = e.Doc()
	}
	return pretty.Group(pretty.Join(",", d...))
}

func (node SelectExpr) Doc() pretty.Doc {
	d := DocExpr(node.Expr)
	if node.As != "" {
		// This is pretty expensive. Use a cheaper solution for now.
		/*
			d = pretty.Group(pretty.Concat(
				d,
				pretty.Nest(1, pretty.Concat(
					pretty.Line,
					pretty.Group(pretty.Fold(pretty.Concat,
						pretty.Text("AS"),
						pretty.Line,
						pretty.Text(node.As.String()),
					)),
				)),
			))
		*/
		d = pretty.Group(pretty.Concat(
			d,
			pretty.Nest(1, pretty.Concat(
				pretty.Line,
				pretty.Text("AS "+node.As.String()),
			)),
		))
	}
	return d
}

func (t TableExprs) Doc() pretty.Doc {
	if len(t) == 0 {
		return pretty.Nil
	}
	d := make([]pretty.Doc, len(t))
	for i, e := range t {
		d[i] = DocNode(e)
	}
	return pretty.JoinGroup("FROM", ",", d...)
}

func (w *Where) Doc() pretty.Doc {
	if w == nil {
		return pretty.Nil
	}
	return pretty.Concat(pretty.Line, pretty.JoinGroup(w.Type, "", DocExpr(w.Expr)))
}

func (node GroupBy) Doc() pretty.Doc {
	if len(node) == 0 {
		return pretty.Nil
	}
	d := make([]pretty.Doc, len(node))
	for i, e := range node {
		d[i] = DocExpr(e)
	}
	return pretty.Concat(pretty.Line, pretty.JoinGroup("GROUP BY", ",", d...))
}

func (node *NormalizableTableName) Doc() pretty.Doc {
	return DocNode(node.TableNameReference)
}

func (node *AndExpr) Doc() pretty.Doc {
	return docBinaryOp(node.Left, node.Right, "AND")
}

func (node *OrExpr) Doc() pretty.Doc {
	return docBinaryOp(node.Left, node.Right, "OR")
}

func docBinaryOp(l, r Expr, op string) pretty.Doc {
	return pretty.Group(pretty.Concat(
		DocExpr(l),
		pretty.Concat(
			pretty.Line,
			pretty.Group(pretty.Fold(pretty.Concat,
				pretty.Text(op),
				pretty.Line,
				DocExpr(r),
			)),
		),
	))
}

func (node *Exprs) Doc() pretty.Doc {
	if node == nil || len(*node) == 0 {
		return pretty.Nil
	}
	d := make([]pretty.Doc, len(*node))
	for i, e := range *node {
		d[i] = DocExpr(e)
	}
	return pretty.Join(",", d...)
}

func (node *BinaryExpr) Doc() pretty.Doc {
	var pad pretty.Doc = pretty.Nil
	if node.Operator.isPadded() {
		pad = pretty.Line
	}
	return pretty.Group(pretty.Concat(
		DocExpr(node.Left),
		pretty.Concat(
			pad,
			pretty.Group(pretty.Fold(pretty.Concat,
				pretty.Text(node.Operator.String()),
				pad,
				DocExpr(node.Right),
			)),
		),
	))
}

func (node *ParenExpr) Doc() pretty.Doc {
	return pretty.Bracket("(", DocExpr(node.Expr), ")")
}

func (node *ParenSelect) Doc() pretty.Doc {
	return pretty.Bracket("(", Doc(node.Select), ")")
}

func (node *Limit) Doc() pretty.Doc {
	if node == nil {
		return pretty.Nil
	}
	var count, offset pretty.Doc
	if node.Count != nil {
		count = pretty.Concat(
			pretty.Line,
			pretty.JoinGroup("LIMIT", "", DocExpr(node.Count)),
		)
	}
	if node.Offset != nil {
		offset = pretty.Concat(
			pretty.Line,
			pretty.JoinGroup("OFFSET", "", DocExpr(node.Offset)),
		)
	}
	return pretty.Concat(count, offset)
}

func (node *OrderBy) Doc() pretty.Doc {
	if node == nil || len(*node) == 0 {
		return pretty.Nil
	}
	d := make([]pretty.Doc, len(*node))
	for i, e := range *node {
		d[i] = DocNode(e)
	}
	return pretty.Concat(
		pretty.Line,
		pretty.JoinGroup("ORDER BY", ",", d...),
	)
}

func (node Select) Doc() pretty.Doc {
	return pretty.Fold(pretty.Concat,
		node.With.Doc(),
		Doc(node.Select),
		node.OrderBy.Doc(),
		node.Limit.Doc(),
	)
}

func (node SelectClause) Doc() pretty.Doc {
	if node.TableSelect {
		return pretty.JoinGroup("TABLE", "", DocNode(node.From.Tables[0]))
	}
	var distinct pretty.Doc
	if node.Distinct {
		distinct = pretty.Text(" DISTINCT")
	}
	return pretty.Fold(pretty.Concat,
		pretty.Group(pretty.Fold(pretty.Concat,
			pretty.Text("SELECT"),
			distinct,
			pretty.Nest(1, pretty.Concat(
				pretty.Line,
				pretty.Group(node.Exprs.Doc()),
			)),
		)),
		node.From.Doc(),
		node.Where.Doc(),
		node.GroupBy.Doc(),
		node.Having.Doc(),
		node.Window.Doc(),
	)
}

func (node *From) Doc() pretty.Doc {
	if node == nil || len(node.Tables) == 0 {
		return pretty.Nil
	}
	d := pretty.Concat(pretty.Line, node.Tables.Doc())
	if node.AsOf.Expr != nil {
		d = pretty.Concat(
			d,
			pretty.Group(pretty.Nest(1, pretty.Concat(
				pretty.Line,
				DocNode(&node.AsOf),
			))),
		)
	}
	return d
}

func (node *Window) Doc() pretty.Doc {
	if node == nil || len(*node) == 0 {
		return pretty.Nil
	}
	d := make([]pretty.Doc, len(*node))
	for i, e := range *node {
		d[i] = pretty.Fold(pretty.Concat,
			pretty.Text(e.Name.String()),
			pretty.Text(" AS "),
			DocNode(e),
		)
	}
	return pretty.Concat(
		pretty.Line,
		pretty.JoinGroup("WINDOW", ",", d...),
	)
}

func (node *With) Doc() pretty.Doc {
	if node == nil {
		return pretty.Nil
	}
	return pretty.Concat(
		DocNode(node),
		pretty.Line,
	)
}

func (node *Subquery) Doc() pretty.Doc {
	d := pretty.Text("<unknown>")
	if node.Select != nil {
		d = Doc(node.Select)
	}
	if node.Exists {
		d = pretty.Concat(
			pretty.Text("EXISTS"),
			d,
		)
	}
	return d
}

func (node *AliasedTableExpr) Doc() pretty.Doc {
	d := DocNode(node.Expr)
	if node.Hints != nil {
		d = pretty.Concat(
			d,
			DocNode(node.Hints),
		)
	}
	if node.Ordinality {
		d = pretty.Concat(
			d,
			pretty.Text(" WITH ORDINALITY"),
		)
	}
	if node.As.Alias != "" {
		d = pretty.Fold(pretty.Concat,
			d,
			pretty.Text(" AS "),
			DocNode(&node.As),
		)
	}
	return d
}

func (node *FuncExpr) Doc() pretty.Doc {
	//return pretty.Text(AsString(node))
	d := node.Exprs.Doc()
	if node.Type != 0 {
		d = pretty.Concat(
			pretty.Text(funcTypeName[node.Type]+" "),
			d,
		)
	}

	d = pretty.Bracket(
		AsString(&node.Func)+"(",
		d,
		")",
	)

	if window := node.WindowDef; window != nil {
		var over pretty.Doc
		if window.Name != "" {
			over = DocNode(&window.Name)
		} else {
			over = DocNode(window)
		}
		d = pretty.Concat(
			d,
			pretty.Concat(
				pretty.Text(" OVER "),
				over,
			),
		)
	}
	if node.Filter != nil {
		d = pretty.Fold(pretty.Concat,
			d,
			pretty.Text(" FILTER (WHERE "),
			DocNode(node.Filter),
			pretty.Text(")"),
		)
	}
	return d
}

func (node *ComparisonExpr) Doc() pretty.Doc {
	if node.Operator.hasSubOperator() {
		return DocNode(node)
	}

	opStr := node.Operator.String()
	if node.Operator == IsDistinctFrom && (node.Right == DNull || node.Right == DBoolTrue || node.Right == DBoolFalse) {
		opStr = "IS NOT"
	} else if node.Operator == IsNotDistinctFrom && (node.Right == DNull || node.Right == DBoolTrue || node.Right == DBoolFalse) {
		opStr = "IS"
	}

	return pretty.Group(pretty.ConcatLine(
		DocExpr(node.Left),
		pretty.Group(pretty.ConcatLine(
			pretty.Text(opStr),
			DocExpr(node.Right),
		)),
	))
}

func (node *NumVal) Doc() pretty.Doc         { return pretty.Text(node.String()) }
func (node *Placeholder) Doc() pretty.Doc    { return pretty.Text(node.String()) }
func (node UnqualifiedStar) Doc() pretty.Doc { return pretty.Text(node.String()) }
func (node *UnresolvedName) Doc() pretty.Doc { return pretty.Text(node.String()) }
