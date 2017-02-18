%{
package main
%}

%union {
    str string
    nl nameList
    f field
    fl fieldList
    d def
    dl defList
}

%type <d> tdecl
%type <dl> top tdecl_list
%type <fl> field_list
%type <f> field
%type <nl> name_list
%type <str> opt_sql
%token <str> IDENT STR
%token SUM DEF IDENT ERROR SQL

%%

top : tdecl_list { irgenlex.(*Scanner).results = $1 }

tdecl_list:
{ $$ = defList(nil) }
| tdecl_list tdecl { $$ = append($1, $2) }
;

tdecl:
SUM IDENT '=' name_list ';'
{ $$ = def{name:$2, pos:irgenlex.(*Scanner).lastPos, t:union, u:$4, f:nil, sql:""} }
| DEF IDENT '{' field_list opt_sql '}' opt_semi
{ $$ = def{name:$2, pos:irgenlex.(*Scanner).lastPos, t:rec, u:nil, f:$4, sql:$5} }
;

opt_sql: { $$ = "" } | SQL STR opt_semi { $$ = $2 } ;

opt_semi: ';' | ;

name_list:
IDENT { $$ = append(nameList(nil), $1) }
| name_list '|' IDENT { $$ = append($1, $3) }
;

field_list:
 { $$ = fieldList(nil) }
| field_list field { $$ = append($1, $2) };
field:
IDENT IDENT '*' ';' { $$ = field{$1, $2, true} }
| IDENT IDENT ';' { $$ = field{$1, $2, false} }
 ;
