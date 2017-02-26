package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"unicode"
)

func main() {
	if len(os.Args) < 3 {
		errOut(fmt.Sprintf("usage: %s <input> <outputdir>", os.Args[0]))
	}
	ifile := os.Args[1]
	dl, err := parse(ifile)
	check(err)

	ext := filepath.Ext(ifile)
	var basename string
	if len(os.Args) > 3 {
		basename = os.Args[3]
	} else {
		basename = filepath.Base(ifile[:len(ifile)-len(ext)])
	}
	basetarget := filepath.Join(os.Args[2], basename)

	process(basetarget, basename, dl)
}

func check(err error) {
	if err != nil {
		errOut(err.Error())
	}
}

func errOut(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func parse(fname string) (defList, error) {
	contents, err := ioutil.ReadFile(fname)
	check(err)
	fileSet := token.NewFileSet()
	f := fileSet.AddFile(fname, -1, len(contents))
	s := Scanner{
		convPos: func(pos token.Pos) string {
			return fileSet.PositionFor(pos, true).String()
		},
	}
	s.s.Init(f, contents, nil, 0 /* skip comments */)
	p := irgenParserImpl{}
	if p.Parse(&s) != 0 {
		return nil, fmt.Errorf("%s: %s", s.lastPos, s.lastError)
	}
	if s.results == nil {
		return nil, fmt.Errorf("%s: no definition until end of input", s.lastPos)
	}
	return s.results, nil
}

func process(tbase string, pkg string, defs defList) {
	m := make(map[string]def)
	recs := make(map[string]struct{})
	unions := make(map[string]struct{})
	allnames := make([]string, len(defs))
	for i, d := range defs {
		allnames[i] = d.name
		if _, ok := m[d.name]; ok {
			errOut(fmt.Sprintf("%s: duplicate name: %s", d.pos, d.name))
		}

		m[d.name] = d
		switch d.t {
		case union:
			unions[d.name] = struct{}{}
		case rec:
			recs[d.name] = struct{}{}
		}
	}
	sort.Strings(allnames)

	genproto(pkg, tbase, allnames, m, recs, unions)
	genbridge(pkg, tbase, allnames, m, recs, unions)
	genalloc(pkg, tbase, allnames, m, recs, unions)
	genwalk(pkg, tbase, allnames, m, recs, unions)
	gensql(pkg, tbase, allnames, m, recs, unions)
}

func gensql(
	pkg, tbase string, allnames []string, m map[string]def, recs, unions map[string]struct{},
) {
	p, err := os.Create(tbase + ".sql.go")
	check(err)
	defer p.Close()
	w := bufio.NewWriter(p)
	defer w.Flush()

	fmt.Fprintf(w, `// Code generated by irgen.
// GENERATED FILE DO NOT EDIT
package %s
import "text/template"
`, pkg)

	tmpl := make(map[string]string)
	for name := range recs {
		d := m[name]
		if d.sql != "" {
			tmpl[name] = d.sql
		}
	}

	if len(tmpl) > 0 {
		fmt.Fprintln(w, `import "bytes"`)
	}

	for name := range tmpl {
		fmt.Fprintf(w, "func (n *%s) FormatSQL(buf *bytes.Buffer) {", name)
		fmt.Fprintf(w, " if err := tmpl%s.Execute(buf, n); err != nil { panic(err) } }\n", name)
	}
	for name := range tmpl {
		fmt.Fprintf(w, "func (n *%s) SQL() string { var buf bytes.Buffer; n.FormatSQL(&buf); return buf.String() }\n", name)
	}

	lastTemplate := "template"
	for name, sql := range tmpl {
		fmt.Fprintf(w, "var tmpl%s = func () *template.Template { ret, err := %s.New(%q).Parse(%s); if err != nil { panic(err) }; return ret }()\n",
			name, lastTemplate, name, sql)
		lastTemplate = "tmpl" + name
	}
	for name := range unions {
		d := m[name]
		fmt.Fprintf(w, "var tmpl%s = func() *template.Template { ret, err := %s.New(%q).Parse(`",
			name, lastTemplate, name)
		lastTemplate = "tmpl" + name
		for i, membername := range d.u {
			elsestr := ""
			if i > 0 {
				elsestr = "else "
			}
			fmt.Fprintf(w, `{{%sif .Get%s}}`, elsestr, membername)
			if _, ok := tmpl[membername]; ok {
				fmt.Fprintf(w, `{{template "%s" .Get%s}}`, membername, membername)
			} else {
				fmt.Fprintf(w, `{{.Get%s}}`, membername)
			}
			fmt.Fprintf(w, "`+\n`")
		}
		fmt.Fprintf(w, "{{end}}`); if err != nil { panic(err) }; return ret }()\n")
	}

}

func genwalk(
	pkg, tbase string, allnames []string, m map[string]def, recs, unions map[string]struct{},
) {
	p, err := os.Create(tbase + ".walk.go")
	check(err)
	defer p.Close()
	w := bufio.NewWriter(p)
	defer w.Flush()

	fmt.Fprintf(w, `// Code generated by irgen.
// GENERATED FILE DO NOT EDIT
package %s

`, pkg)

	// We need the set of all rec types which are also members of unions.
	unionized := make(map[string]string)
	for name := range unions {
		for _, membername := range m[name].u {
			unionized[membername] = name
		}
	}
	// Compute the invert set.
	nonunionized := make(map[string]string)
	for name := range recs {
		if _, ok := unionized[name]; !ok {
			nonunionized[name] = "*" + name
		}
	}

	// Generate the visitor interface.
	fmt.Fprintf(w, "type Visitor struct {\n")
	for name := range nonunionized {
		fmt.Fprintf(w, " VisitPre%s func(node *%s, copy *%s) (recurse bool, newNode *%s)\n", name, name, name, name)
		fmt.Fprintf(w, " VisitPost%s func(node *%s, copy *%s) (newNode *%s)\n", name, name, name, name)
	}
	for name := range unions {
		fmt.Fprintf(w, " VisitPre%s func(node %s) (recurse bool, newNode %s)\n", name, name, name)
		fmt.Fprintf(w, " VisitPost%s func(node %s) (newNode %s)\n", name, name, name)
	}
	fmt.Fprintf(w, "}\n")

	// Generate the (*Visitor).WalkXXX() jump methods.
	for name := range nonunionized {
		fmt.Fprintf(w, "\nfunc (v *Visitor) Walk%s(c *AllocContext, node %s) (newNode %s, changed bool) {\n",
			name, name, name)
		fmt.Fprintf(w, " recurse, nodeRef := true, &node\n")
		fmt.Fprintf(w, " if v.VisitPre%s != nil { recurse, nodeRef = v.VisitPre%s(nodeRef, &newNode) }\n",
			name, name)
		fmt.Fprintf(w, " if recurse {\n  nodeRef = nodeRef.Walk(c, v, &newNode)\n")
		fmt.Fprintf(w, "  if v.VisitPost%s != nil { nodeRef = v.VisitPost%s(nodeRef, &newNode) }\n }\n", name, name)
		fmt.Fprintf(w, " return newNode, nodeRef != &node\n}\n")
	}
	for name := range unions {
		fmt.Fprintf(w, "\nfunc (v *Visitor) Walk%s(c *AllocContext, node %s) (newNode %s, changed bool) {\n",
			name, name, name)
		fmt.Fprintf(w, " recurse, newNode := true, node\n")
		fmt.Fprintf(w, " if v.VisitPre%s != nil { recurse, newNode = v.VisitPre%s(node) }\n",
			name, name)
		fmt.Fprintf(w, " if recurse {\n  newNode = newNode.Walk(c, v)\n")
		fmt.Fprintf(w, "  if v.VisitPost%s != nil { newNode = v.VisitPost%s(newNode) }\n }\n", name, name)
		fmt.Fprintf(w, " return newNode, node != newNode\n}\n")
	}

	// Generate the (*Node).Walk() methods.
	for name := range recs {
		d := m[name]

		// The return type of the method is the interface if the record belongs
		// to a union, otherwise it is the record type itself.
		targettype := "(ret *" + name + ")"
		isunionized := false
		if newtarget, ok := unionized[name]; ok {
			isunionized = true
			targettype = newtarget
		}

		if isunionized {
			fmt.Fprintf(w, "\nfunc (n *%s) Walk(c *AllocContext, v *Visitor) %s {\n", name, targettype)
			fmt.Fprintf(w, " ret := n\n")
		} else {
			fmt.Fprintf(w, "\nfunc (n *%s) Walk(c *AllocContext, v *Visitor, copy *%s) %s {\n", name, name, targettype)
			fmt.Fprintf(w, " ret = n\n")
		}

		// numFields is the number of non-immediate fields (either record or union types).
		numFields := 0

		// cond is the accumulated predicate that determines whether any
		// of the non-immediate fields were changed.
		cond := ""

		// The following loop constructs the recursive calls to Walk on
		// the array fields only.
		for _, f := range d.f {
			if isImmediateType(f.t) || !f.isArray {
				continue
			}

			// The field name in the struct is the uppercased version of the
			// name in the input.
			funame := upperName(f.name)

			getFieldName := funame
			if _, ok := unions[f.t]; ok {
				getFieldName = "X" + funame
			}
			fmt.Fprintf(w, " for i := range n.%s {\n", getFieldName)
			fmt.Fprintf(w, "  newVal, changed := v.Walk%s(c, n.%s_At(i))\n", f.t, funame)
			fmt.Fprintf(w, "  if changed {\n")
			if isunionized {
				fmt.Fprintf(w, "   if ret == n { ret = n.Clone(c) }\n")
			} else {
				fmt.Fprintf(w, "   if ret != copy { n.CopyTo(c, copy); ret = copy }\n")
			}

			fmt.Fprintf(w, "   ret.Set%s_At(c, i, newVal)\n  }\n }\n", funame)
		}

		// The following loop constructs the recursive calls to Walk on
		// the non-array fields only.
		for _, f := range d.f {
			if isImmediateType(f.t) || f.isArray {
				continue
			}

			// Accumulate the condition for below.
			cond = cond + "changed_" + f.name + " || "
			numFields++

			// The field name in the struct is the uppercased version of the
			// name in the input.
			funame := upperName(f.name)

			// To access the field we either use the field directly (for
			// types that are not part of unions) or use the accessor.
			getFieldName := funame
			if _, ok := unions[f.t]; ok {
				getFieldName = funame + "()"
			}
			fmt.Fprintf(w, " new_%s, changed_%s := v.Walk%s(c, n.%s)\n",
				f.name, f.name, f.t, getFieldName)
		}

		// The following generates the condition on whether any of the
		// fields were changed.
		if numFields > 0 {
			fmt.Fprintf(w, " if %sfalse {\n", cond)
			if isunionized {
				fmt.Fprintf(w, "  if ret == n { ret = n.Clone(c) }\n")
			} else {
				fmt.Fprintf(w, "  if ret != copy { n.CopyTo(c, copy); ret = copy }\n")
			}
			for _, f := range d.f {
				if f.isArray {
					continue
				}
				funame := upperName(f.name)
				if _, ok := unions[f.t]; ok {
					fmt.Fprintf(w, "  ret.Set%s(c, new_%s)\n", funame, f.name)
				} else if isNodeType(f.t) {
					fmt.Fprintf(w, "  ret.%s = new_%s\n", funame, f.name)
				} else {
					fmt.Fprintf(w, "  ret.%s = n.%s\n", funame, funame)
				}
			}
			fmt.Fprintf(w, " }\n")
		}

		fmt.Fprintf(w, " return ret\n}\n")
	}
}

func genalloc(
	pkg, tbase string, allnames []string, m map[string]def, recs, unions map[string]struct{},
) {
	p, err := os.Create(tbase + ".alloc.go")
	check(err)
	defer p.Close()
	w := bufio.NewWriter(p)
	defer w.Flush()

	// We need to detect the types for which we need array constructors
	// and converters.
	arrayTypes := make(map[string]struct{})
	for name := range recs {
		d := m[name]
		for _, f := range d.f {
			if f.isArray {
				arrayTypes[f.t] = struct{}{}
			}
		}
	}

	fmt.Fprintf(w, `// Code generated by irgen.
// GENERATED FILE DO NOT EDIT
package %s

type AllocContext struct {
`, pkg)
	for name := range arrayTypes {
		name = realname(unions, name)
		fmt.Fprintf(w, "  pool%s []%s\n", name, name)
	}
	for _, name := range allnames {
		if _, ok := arrayTypes[name]; ok {
			continue
		}
		name = realname(unions, name)
		fmt.Fprintf(w, "  pool%s []%s\n", name, name)
	}
	for name := range unions {
		d := m[name]
		for _, membername := range d.u {
			fmt.Fprintf(w, "  poolU_%s_%s []U_%s_%s\n", name, membername, name, membername)
		}
	}
	fmt.Fprintln(w, "}")
	for _, name := range allnames {
		fmt.Fprintf(w, "\n// Allocator interface for %s.\n", name)
		_, isUnion := unions[name]
		pub := "N"
		origName := name
		if isUnion {
			pub = "n"
			name = "U_" + name
		}
		fmt.Fprintf(w, "func (c *AllocContext) %sew%s() *%s { return c.%sew%s_copy(%s{}) }\n",
			pub, name, name, pub, name, name)
		if !isUnion {
			d := m[name]
			fmt.Fprintf(w, "func (c *AllocContext) %s(", name)
			for _, f := range d.f {
				array := ""
				if f.isArray {
					array = "[]"
				}
				fmt.Fprintf(w, "f_%s %s%s, ", f.name, array, f.t)
			}
			fmt.Fprintf(w, ") *%s {\n return c.%sew%s_copy(%s{", name, pub, name, name)
			for _, f := range d.f {
				funame := upperName(f.name)
				if _, ok := unions[f.t]; ok {
					if f.isArray {
						fmt.Fprintf(w, "X%s:c.convArray%s(f_%s), ", funame, f.t, f.name)
					} else {
						fmt.Fprintf(w, "X%s:c.wrap%s(f_%s), ", funame, f.t, f.name)
					}
				} else {
					fmt.Fprintf(w, "%s:f_%s, ", funame, f.name)
				}
			}
			fmt.Fprintf(w, "})\n}\n")
		}
		genalloc1(w, pub, name)
		if isUnion {
			fmt.Fprintf(w, "func (c *AllocContext) wrap%s(v %s) (r %s) {", origName, origName, name)
			fmt.Fprintf(w, " if v != nil { r.set(c, v) }; return r }\n")
		}
	}

	fmt.Fprintf(w, "\n// Allocators for the intermediary structs between a setof and the concrete types.\n")
	for name := range unions {
		d := m[name]
		for _, membername := range d.u {
			genalloc1(w, "n", "U_"+name+"_"+membername)
		}
	}

	fmt.Fprintf(w, "\n// Clone methods.\n")
	for name := range recs {
		d := m[name]
		fmt.Fprintf(w, "func (n *%s) Clone(c *AllocContext) (ret *%s) {\n", name, name)
		fmt.Fprintf(w, " ret = c.New%s()\n", name)
		fmt.Fprintf(w, " n.CopyTo(c, ret)\n return ret\n}\n")
		fmt.Fprintf(w, "func (n *%s) CopyTo(c *AllocContext, ret *%s) {\n", name, name)
		for _, f := range d.f {
			funame := upperName(f.name)
			fieldName := funame
			allocType := f.t
			if _, ok := unions[f.t]; ok {
				fieldName = "X" + fieldName
				allocType = "U_" + allocType
			}
			if f.isArray {
				fmt.Fprintf(w, " ret.%s = c.newArray%s(len(n.%s), len(n.%s))\n",
					fieldName, allocType, fieldName, fieldName)
				fmt.Fprintf(w, " for i := range n.%s {\n", fieldName)
				fmt.Fprintf(w, "  ret.Set%s_At(c, i, n.%s_At(i))\n", funame, funame)
				fmt.Fprintf(w, " }\n")
			} else {
				fmt.Fprintf(w, " ret.%s = n.%s\n", fieldName, fieldName)
			}
		}
		fmt.Fprintf(w, "}\n")
	}

	fmt.Fprintf(w, "\n// Allocators for the arrays.\n")
	for name := range arrayTypes {
		allocType := name
		if _, ok := unions[name]; ok {
			fmt.Fprintf(w, "func (c *AllocContext) convArray%s(a []%s) (r []U_%s) {\n", name, name, name)
			fmt.Fprintf(w, " if len(a) == 0 { return nil }\n r = c.newArrayU_%s(len(a), len(a))\n", name)
			fmt.Fprintf(w, " for i, v := range a { r[i].set(c, v) }\n return r\n}\n")
			allocType = "U_" + allocType
		}
		fmt.Fprintf(w, "func (c *AllocContext) newArray%s(l, max int) []%s {\n", allocType, allocType)
		fmt.Fprintf(w, "  if len(c.pool%s) >= cap(c.pool%s)-max {\n", allocType, allocType)
		fmt.Fprintf(w, "    sz := POOL_ALLOC_SIZE; if (max > sz) { sz = max };\n")
		fmt.Fprintf(w, "    c.pool%s = make([]%s, 0, sz)\n", allocType, allocType)
		fmt.Fprintf(w, "  }\n  c.pool%s = c.pool%s[:len(c.pool%s)+max]\n", allocType, allocType, allocType)
		fmt.Fprintf(w, "  return c.pool%s[len(c.pool%s)-max:len(c.pool%s)-max+l]\n", allocType, allocType, allocType)
		fmt.Fprintf(w, "}\n")
	}
}

func genalloc1(w io.Writer, pub, name string) {
	fmt.Fprintf(w, "func (c *AllocContext) %sew%s_copy(model %s) *%s {\n", pub, name, name, name)
	fmt.Fprintf(w, "  if len(c.pool%s) == cap(c.pool%s) {\n", name, name)
	fmt.Fprintf(w, "    c.pool%s = make([]%s, 0, POOL_ALLOC_SIZE)\n", name, name)
	fmt.Fprintf(w, "  }\n  c.pool%s = append(c.pool%s, model)\n", name, name)
	fmt.Fprintf(w, "  return &c.pool%s[len(c.pool%s)-1]\n}\n", name, name)
}

func genbridge(
	pkg, tbase string, allnames []string, m map[string]def, recs, unions map[string]struct{},
) {
	p, err := os.Create(tbase + ".bridge.go")
	check(err)
	defer p.Close()
	w := bufio.NewWriter(p)
	defer w.Flush()

	fmt.Fprintf(w, `// Code generated by irgen.
// GENERATED FILE DO NOT EDIT
package %s

`, pkg)
	fmt.Fprintf(w, "\n// Union interfaces and type relationships.\n")
	for name := range unions {
		// Generate the interface binding the union type to its instances.
		// sum A = B | C -> interface A {}; func (*B) isA() {}; func(*C) isA() {}
		fmt.Fprintf(w, "\ntype %s interface { is%s(); Walk(c *AllocContext, v *Visitor) %s }\n", name, name, name)
		d := m[name]
		for _, membername := range d.u {
			fmt.Fprintf(w, "func (* %s) is%s() {}\n", realname(unions, membername), name)
		}

		// Generate the field getter and setters for all nodes that have members
		// of the type of this union.
		// We need this extra hoop because go-protobuf generates intermediate
		// structs between a setof declaration and the actual members.
		fmt.Fprintf(w, "\n// Low-level getter and setter for union %s.\n", name)
		fmt.Fprintf(w, "func (n *U_%s) get() %s {\n switch v := n.V.(type) {\n", name, name)
		for _, membername := range d.u {
			mname := realname(unions, membername)
			fmt.Fprintf(w, "  case *U_%s_%s: if v.%s == nil { return nil }; return v.%s\n", name, mname, mname, mname)
		}
		fmt.Fprintf(w, " }\n return nil\n}\n")
		fmt.Fprintf(w, "func (n *U_%s) set(c *AllocContext, val %s) {\n switch v := val.(type) {\n", name, name)
		for _, membername := range d.u {
			mname := realname(unions, membername)
			fmt.Fprintf(w, "  case *%s: n.V = c.newU_%s_%s_copy(U_%s_%s{%s:v})\n", mname, name, mname, name, mname, mname)
		}
		fmt.Fprintf(w, " }\n}\n")
	}

	for name := range recs {
		for _, f := range m[name].f {
			_, isUnion := unions[f.t]
			if !isUnion && !f.isArray {
				continue
			}
			funame := upperName(f.name)
			fieldName := funame
			allocType := f.t
			itemGet := ""
			if isUnion {
				fieldName = "X" + fieldName
				allocType = "U_" + allocType
				itemGet = ".get()"
			}

			fmt.Fprintf(w, "\n// Accessors for %s.%s.\n", name, funame)
			if f.isArray {
				fmt.Fprintf(w, "func (n *%s) %s_Len() int { return len(n.%s) }\n",
					name, funame, fieldName)
				fmt.Fprintf(w, "func (n *%s) %s_At(i int) %s { return n.%s[i]%s }\n",
					name, funame, f.t, fieldName, itemGet)
				fmt.Fprintf(w, "func (n *%s) Reserve%s(c *AllocContext, sz int) {\n",
					name, funame)
				fmt.Fprintf(w, " if sz == 0 { n.%s = nil\n", fieldName)
				fmt.Fprintf(w, " } else { n.%s = c.newArray%s(0, sz) }\n}\n", fieldName, allocType)
				fmt.Fprintf(w, "func (n *%s) Set%s(c *AllocContext, val []%s) {\n",
					name, funame, f.t)
				if isUnion {
					fmt.Fprintf(w, " if val == nil { n.%s = nil\n", fieldName)
					fmt.Fprintf(w, " } else { n.%s = c.convArray%s(val) }\n}\n", fieldName, f.t)
				} else {
					fmt.Fprintf(w, " n.%s = val\n}\n", fieldName)
				}
				fmt.Fprintf(w, "func (n *%s) Set%s_At(c *AllocContext, i int, val %s) {\n",
					name, funame, f.t)
				if isUnion {
					fmt.Fprintf(w, " if val == nil { n.%s[i].Reset()\n", fieldName)
					fmt.Fprintf(w, " } else { n.%s[i].set(c, val) }\n}\n", fieldName)
				} else {
					fmt.Fprintf(w, " n.%s[i] = val\n}\n", fieldName)
				}
				fmt.Fprintf(w, "func (n *%s) %s_Append(c *AllocContext, val %s) {\n",
					name, funame, f.t)
				if isUnion {
					fmt.Fprintf(w, " n.%s = append(n.%s, %s{})\n", fieldName, fieldName, allocType)
					fmt.Fprintf(w, " n.Set%s_At(c, len(n.%s)-1, val)\n}\n", funame, fieldName)
				} else {
					fmt.Fprintf(w, " n.%s = append(n.%s, val)\n}\n", fieldName, fieldName)
				}
			} else if isUnion {
				fmt.Fprintf(w, "func (n *%s) %s() %s { return n.%s.get() }\n",
					name, funame, f.t, fieldName)
				fmt.Fprintf(w, "func (n *%s) Set%s(c *AllocContext, val %s) {\n",
					name, funame, f.t)
				fmt.Fprintf(w, " if val == nil { n.%s.Reset()\n", fieldName)
				fmt.Fprintf(w, " } else { n.%s.set(c, val) }\n}\n", fieldName)
			}
		}
	}
}

func genproto(pkg, tbase string, allnames []string, m map[string]def, recs, unions map[string]struct{}) {
	p, err := os.Create(tbase + ".proto")
	check(err)
	defer p.Close()
	w := bufio.NewWriter(p)
	defer w.Flush()

	fmt.Fprintf(w, `// Code generated by irgen.
// GENERATED FILE DO NOT EDIT
syntax = "proto2";
package cockroach.sql.ir;
option go_package = "%s";
import "gogoproto/gogo.proto";

`, pkg)

	for _, name := range allnames {
		d := m[name]
		name = realname(unions, name)
		fmt.Fprintf(w, "message %s {\n", name)
		switch d.t {
		case union:
			fmt.Fprintf(w, " oneof v {\n")
			for i, n := range d.u {
				fmt.Fprintf(w, "  %s %s = %d;\n", realname(unions, n), lowerName(n), i+1)
			}
			fmt.Fprintf(w, " }\n")
		case rec:
			for i, f := range d.f {
				rep := "optional"
				if f.isArray {
					rep = "repeated"
				}
				nullable := " [(gogoproto.nullable) = false]"
				if f.isArray && !(f.t[0] >= 'A' && f.t[0] <= 'Z') {
					nullable = ""
				}
				fname, ftype := f.name, f.t
				if _, ok := unions[ftype]; ok {
					fname = "_" + fname
					ftype = "U_" + ftype
				}
				fmt.Fprintf(w, " %s %s %s = %d%s;\n", rep, ftype, fname, i+1, nullable)
			}
		}
		fmt.Fprint(w, "}\n\n")
	}
}

func realname(unions map[string]struct{}, name string) string {
	if _, ok := unions[name]; ok {
		return "U_" + name
	}
	return name
}

func upperName(name string) string {
	// blah_blah -> BlahBlah
	var buf bytes.Buffer
	seenUnderscore := false
	for i, c := range name {
		if i == 0 || seenUnderscore {
			buf.WriteRune(unicode.ToUpper(c))
			continue
		}
		seenUnderscore = false
		if c == '_' {
			seenUnderscore = true
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

func isNodeType(t string) bool {
	return t[0] >= 'A' && t[0] <= 'Z'
}

func isImmediateType(t string) bool {
	return !isNodeType(t)
}

func lowerName(name string) string {
	// BlahBlah -> blah_blah
	var buf bytes.Buffer
	for i, c := range name {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				buf.WriteRune('_')
			}
			c = c - 'A' + 'a'
		}
		buf.WriteRune(c)
	}
	return buf.String()
}
