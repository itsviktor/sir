package postgres

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

func parseExpr(aExpr *parser.A_exprContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	return parseQual(aExpr.A_expr_qual().(*parser.A_expr_qualContext), scope, tctx)
}

func parseQual(
	aExpr *parser.A_expr_qualContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseLessLess(
		aExpr.A_expr_lessless().(*parser.A_expr_lesslessContext),
		scope,
		tctx,
	)
}

func parseLessLess(
	aExpr *parser.A_expr_lesslessContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	left := parseOr(
		aExpr.A_expr_or(0).(*parser.A_expr_orContext),
		scope,
		tctx,
	)

	for i := 1; i < len(aExpr.AllA_expr_or()); i++ {
		rightExpr := aExpr.A_expr_or(i)
		if rightExpr == nil {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid or context: nil right at index %d", i)
		}
		right := parseOr(rightExpr.(*parser.A_expr_orContext), scope, tctx)

		left = &ir.BinaryExpr{
			Left:  left,
			Right: right,
			Op:    ir.NewOp(ir.Or, "OR"),
		}
	}

	return left
}

func parseOr(
	aExpr *parser.A_expr_orContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	left := parseAnd(aExpr.A_expr_and(0).(*parser.A_expr_andContext), scope, tctx)

	for i := 1; i < len(aExpr.AllA_expr_and()); i++ {
		rightExpr := aExpr.A_expr_and(i)
		if rightExpr == nil {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid and context: nil right at index %d", i)
		}
		right := parseAnd(rightExpr.(*parser.A_expr_andContext), scope, tctx)

		left = &ir.BinaryExpr{
			Left:  left,
			Right: right,
			Op:    ir.NewOp(ir.And, "AND"),
		}
	}

	return left
}

func parseAnd(
	aExpr *parser.A_expr_andContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	left := parseBetween(aExpr.A_expr_between(0).(*parser.A_expr_betweenContext), scope, tctx)

	if aExpr.A_expr_between(1) == nil {
		return left
	}

	right := parseBetween(aExpr.A_expr_between(1).(*parser.A_expr_betweenContext), scope, tctx)

	return &ir.LogicalExpr{
		Right: right,
		Left:  left,
		Op:    ir.NewOp(ir.And, "AND"),
	}
}

func parseBetween(
	aExpr *parser.A_expr_betweenContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	left := parseIn(
		aExpr.A_expr_in(0).(*parser.A_expr_inContext),
		scope,
		tctx,
	)

	if aExpr.BETWEEN() == nil {
		return left
	}

	isNot := aExpr.NOT() != nil

	from := parseIn(
		aExpr.A_expr_in(1).(*parser.A_expr_inContext),
		scope,
		tctx,
	)
	to := parseIn(
		aExpr.A_expr_in(2).(*parser.A_expr_inContext),
		scope,
		tctx,
	)

	return &ir.BetweenExpr{
		Left: left,
		From: from,
		To:   to,
		Not:  isNot,
	}
}

func parseIn(
	aExpr *parser.A_expr_inContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	left := parseUnaryNot(
		aExpr.A_expr_unary_not().(*parser.A_expr_unary_notContext),
		scope,
		tctx,
	)

	if aExpr.IN_P() == nil {
		return left
	}

	isNot := aExpr.NOT() != nil

	right := aExpr.In_expr()
	if right == nil {
		tctx.ErrOnToken(aExpr.GetStart(), "invalid in expr: no right part")
	}

	switch rightNode := right.(type) {
	case *parser.In_expr_listContext:
		listCtx := rightNode.Expr_list()

		var exprs []ir.Expr

		for _, aExpr := range listCtx.AllA_expr() {
			exprs = append(exprs, parseExpr(aExpr.(*parser.A_exprContext), scope, tctx))
		}

		return &ir.InExpr{
			Left: left,
			Not:  isNot,
			Right: ir.InList{
				Items: exprs,
			},
		}
	case *parser.In_expr_selectContext:
		selectClause := rightNode.Select_with_parens().Select_no_parens()
		subqueryScope, ok := scope.FindChildrenByNode(selectClause)
		if !ok {
			tctx.ErrOnToken(selectClause.GetStart(), "cannot found scope for this query")
		}

		subquery := parseSelectQuery(selectClause.(*parser.Select_no_parensContext), subqueryScope, tctx)

		return &ir.InExpr{
			Left: left,
			Not:  isNot,
			Right: ir.InSubquery{
				Query: subquery,
			},
		}
	case *parser.In_expr_dsqlContext:
		expr := parseDsql(rightNode.Dsql_param().(*parser.Dsql_paramContext), scope, tctx)

		return &ir.InExpr{
			Left: left,
			Not:  isNot,
			Right: ir.InDsql{
				Expr: expr,
			},
		}
	}

	return left
}

func parseUnaryNot(
	aExpr *parser.A_expr_unary_notContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	expr := parseIsNull(
		aExpr.A_expr_isnull().(*parser.A_expr_isnullContext),
		scope,
		tctx,
	)

	if aExpr.NOT() == nil {
		return expr
	}

	return &ir.UnaryExpr{
		Expr: expr,
		Op:   ir.NewOp(ir.Not, "NOT"),
	}
}

func parseIsNull(
	aExpr *parser.A_expr_isnullContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	left := parseIsNot(
		aExpr.A_expr_is_not().(*parser.A_expr_is_notContext),
		scope,
		tctx,
	)

	if aExpr.ISNULL() != nil {
		return &ir.IsExpr{
			Expr:      left,
			Predicate: ir.IsNull,
			Not:       false,
		}
	} else if aExpr.NOTNULL() != nil {
		return &ir.IsExpr{
			Expr:      left,
			Predicate: ir.IsNull,
			Not:       true,
		}
	}

	return left
}

func parseIsNot(
	aExpr *parser.A_expr_is_notContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	expr := parseCompare(
		aExpr.A_expr_compare().(*parser.A_expr_compareContext),
		scope,
		tctx,
	)

	if aExpr.IS() == nil {
		return expr
	}

	isNot := aExpr.NOT() != nil

	if aExpr.DISTINCT() != nil {
		right := parseExpr(
			aExpr.A_expr().(*parser.A_exprContext),
			scope,
			tctx,
		)

		return &ir.IsDistinctExpr{
			Left:  expr,
			Right: right,
			Not:   isNot,
		}
	}

	var predicate ir.IsPredicate
	switch {
	case aExpr.NULL_P() != nil:
		predicate = ir.IsNull
	case aExpr.TRUE_P() != nil:
		predicate = ir.IsTrue
	case aExpr.FALSE_P() != nil:
		predicate = ir.IsFalse
	case aExpr.UNKNOWN() != nil:
		predicate = ir.IsUnknown
	case aExpr.DOCUMENT_P() != nil:
		predicate = ir.IsDocument
	case aExpr.NORMALIZED() != nil:
		predicate = ir.IsNormalized
	default:
		tctx.ErrOnToken(aExpr.GetStart(), "invalid is expr: no predicate")
	}

	return &ir.IsExpr{
		Expr:      expr,
		Predicate: predicate,
		Not:       isNot,
	}
}

func parseCompare(
	aExpr *parser.A_expr_compareContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Expr {
	left := parseLike(aExpr.A_expr_like(0).(*parser.A_expr_likeContext), scope, tctx)

	rightExpr := aExpr.A_expr_like(1)
	if rightExpr == nil {
		return left
	}
	right := parseLike(rightExpr.(*parser.A_expr_likeContext), scope, tctx)

	op := aExpr.GetChild(1)
	if op == nil {
		tctx.ErrOnToken(aExpr.GetStart(), "invalid compare expr: no operator")
	}

	opText := op.(antlr.TerminalNode).GetText()
	opType, ok := ir.StringToOpType(opText)
	if !ok {
		tctx.ErrOnToken(aExpr.GetStart(), "unknown binary op: %s", opText)
	}

	return &ir.BinaryExpr{
		Left:  left,
		Right: right,
		Op:    ir.NewOp(opType, opText),
	}
}

func parseLike(aExpr *parser.A_expr_likeContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	left := parseQualOp(aExpr.A_expr_qual_op(0).(*parser.A_expr_qual_opContext), scope, tctx)

	rightExpr := aExpr.A_expr_qual_op(1)
	if rightExpr == nil {
		return left
	}
	right := parseQualOp(rightExpr.(*parser.A_expr_qual_opContext), scope, tctx)

	var opText string
	not := aExpr.NOT()
	ilike := aExpr.ILIKE()

	var op ir.OpType
	switch {
	case not != nil && ilike != nil:
		op = ir.NotILike
		opText = "NOT ILIKE"
	case not != nil:
		op = ir.NotLike
		opText = "NOT LIKE"
	case ilike != nil:
		op = ir.ILike
		opText = "ILIKE"
	default:
		op = ir.Like
		opText = "LIKE"
	}

	return &ir.BinaryExpr{
		Left:  left,
		Right: right,
		Op:    ir.NewOp(op, opText),
	}
}

func parseQualOp(aExpr *parser.A_expr_qual_opContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	left := parseUnaryQualOp(aExpr.A_expr_unary_qualop(0).(*parser.A_expr_unary_qualopContext), scope, tctx)

	for i := 1; i < len(aExpr.AllA_expr_unary_qualop()); i++ {
		rightExpr := aExpr.A_expr_unary_qualop(i)
		if rightExpr == nil {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid qual op context: nil right at index %d", i)
		}
		right := parseUnaryQualOp(rightExpr.(*parser.A_expr_unary_qualopContext), scope, tctx)

		opChild := aExpr.GetChild(2*i - 1)
		if opChild == nil {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid qual op context: no op at index %d", i)
		}
		opChildNode, ok := opChild.(*parser.Qual_opContext)
		if !ok {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid qual op context: invalid op at index %d, wait for Qual_opContext, got %T", i, opChild)
		}

		opText := opChildNode.GetText()
		op, ok := ir.StringToOpType(opText)
		if !ok {
			op = ir.CustomBinary
		}

		left = &ir.BinaryExpr{
			Left:  left,
			Right: right,
			Op:    ir.NewOp(op, opText),
		}
	}

	return left
}

func parseUnaryQualOp(aExpr *parser.A_expr_unary_qualopContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	expr := parseAdd(aExpr.A_expr_add().(*parser.A_expr_addContext), scope, tctx)

	qualOp := aExpr.Qual_op()
	if qualOp == nil {
		return expr
	}

	opNode, ok := qualOp.(*parser.Qual_opContext)
	if !ok {
		tctx.ErrOnToken(aExpr.GetStart(), "invalid unary qual op context: invalid op, wait for Qual_opContext, got %T", qualOp)
	}

	opText := opNode.GetText()
	op, ok := ir.StringToOpType(opText)
	if !ok {
		op = ir.CustomUnary
	}

	return &ir.UnaryExpr{
		Expr: expr,
		Op:   ir.NewOp(op, opNode.GetText()),
	}
}

func parseAdd(aExpr *parser.A_expr_addContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	left := parseMul(aExpr.A_expr_mul(0).(*parser.A_expr_mulContext), scope, tctx)

	for i := 1; i < len(aExpr.AllA_expr_mul()); i++ {
		rightExpr := aExpr.A_expr_mul(i)
		if rightExpr == nil {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid add context: nil right at index %d", i)
		}
		right := parseMul(rightExpr.(*parser.A_expr_mulContext), scope, tctx)

		opChild := aExpr.GetChild(2*i - 1)
		if opChild == nil {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid add context: no op at index %d", i)
		}
		opChildNode, ok := opChild.(antlr.TerminalNode)
		if !ok {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid add context: invalid op at index %d, wait for TerminalNode, got %T", i, opChild)
		}

		opText := opChildNode.GetText()
		op, ok := ir.StringToOpType(opText)
		if !ok {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid add context: invalid op %s at index %d", opChildNode.GetText(), i)
		}

		left = &ir.BinaryExpr{
			Left:  left,
			Right: right,
			Op:    ir.NewOp(op, opText),
		}
	}

	return left
}

func parseMul(aExpr *parser.A_expr_mulContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	left := parseCaret(aExpr.A_expr_caret(0).(*parser.A_expr_caretContext), scope, tctx)

	for i := 1; i < len(aExpr.AllA_expr_caret()); i++ {
		rightExpr := aExpr.A_expr_caret(i)
		if rightExpr == nil {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid mul context: nil right at index %d", i)
		}
		right := parseCaret(rightExpr.(*parser.A_expr_caretContext), scope, tctx)

		opChild := aExpr.GetChild(2*i - 1)
		if opChild == nil {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid mul context: no op at index %d", i)
		}
		opChildNode, ok := opChild.(antlr.TerminalNode)
		if !ok {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid mul context: invalid op at index %d, wait for TerminalNode, got %T", i, opChild)
		}

		opText := opChildNode.GetText()
		op, ok := ir.StringToOpType(opText)
		if !ok {
			tctx.ErrOnToken(aExpr.GetStart(), "invalid mul context: invalid op %s at index %d", opChildNode.GetText(), i)
		}

		left = &ir.BinaryExpr{
			Left:  left,
			Right: right,
			Op:    ir.NewOp(op, opText),
		}
	}

	return left
}

func parseCaret(aExpr *parser.A_expr_caretContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	left := parseUnarySign(aExpr.A_expr_unary_sign(0).(*parser.A_expr_unary_signContext), scope, tctx)

	rightExpr := aExpr.A_expr_unary_sign(1)
	if rightExpr == nil {
		return left
	}
	right := parseUnarySign(rightExpr.(*parser.A_expr_unary_signContext), scope, tctx)

	if aExpr.CARET() == nil {
		tctx.ErrOnToken(aExpr.GetStart(), "invalid caret expr: no caret op")
	}

	return &ir.BinaryExpr{
		Left:  left,
		Right: right,
		Op:    ir.NewOp(ir.Caret, "^"),
	}
}

func parseUnarySign(aExpr *parser.A_expr_unary_signContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	right := parseAtTimeZone(aExpr.A_expr_at_time_zone().(*parser.A_expr_at_time_zoneContext), scope, tctx)

	if aExpr.MINUS() != nil {
		return &ir.UnaryExpr{
			Expr: right,
			Op:   ir.NewOp(ir.UnaryMinus, "-"),
		}
	} else if aExpr.PLUS() != nil {
		return &ir.UnaryExpr{
			Expr: right,
			Op:   ir.NewOp(ir.UnaryPlus, "+"),
		}
	}

	return right
}

func parseAtTimeZone(aExpr *parser.A_expr_at_time_zoneContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	return parseCollate(aExpr.A_expr_collate().(*parser.A_expr_collateContext), scope, tctx)
}

func parseCollate(ctx *parser.A_expr_collateContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	return parseTypecast(ctx.A_expr_typecast().(*parser.A_expr_typecastContext), scope, tctx)
}

func parseTypecast(ctx *parser.A_expr_typecastContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	cCtx := ctx.C_expr()
	return parseCExpr(cCtx, scope, tctx)
}

func parseCExpr(ctx parser.IC_exprContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	cExpr, ok := ctx.(*parser.C_expr_exprContext)
	if ok {
		switch {
		case cExpr.A_expr() != nil:
			return parseExpr(cExpr.A_expr().(*parser.A_exprContext), scope, tctx)
		case cExpr.Columnref() != nil:
			return parseColumnRef(cExpr.Columnref().(*parser.ColumnrefContext), scope, tctx)
		default:
			return &ir.LiteralExpr{Value: ctx.GetText()}
		}
	}

	dsqlExpr, ok := ctx.(*parser.C_expr_dsqlparamContext)
	if ok {
		return parseDsql(dsqlExpr.Dsql_param().(*parser.Dsql_paramContext), scope, tctx)
	}

	tctx.ErrOnToken(ctx.GetStart(), "unknown c expression type %T", ctx)

	return nil
}

func parseColumnRef(ctx *parser.ColumnrefContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	colid, ok := ctx.Colid().(*parser.ColidContext)
	if !ok {
		tctx.ErrOnToken(ctx.GetStart(), "expected colid, got %T", ctx.Colid())
	}

	parts := parseQualifiedParts(colid, ctx.Indirection())

	switch len(parts) {
	case 1:
		return resolveUnqualifiedColumn(ctx, parts[0], scope, tctx)
	case 2:
		return resolveQualifiedColumn(ctx, parts[1], parts[0], scope, tctx)
	case 3:
		return resolveQualifiedColumn(ctx, parts[2], parts[1], scope, tctx)
	default:
		tctx.ErrOnToken(ctx.GetStart(), "invalid column ref qualified parts length: %d", len(parts))
		return nil
	}
}

func resolveUnqualifiedColumn(ctx *parser.ColumnrefContext, name string, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	candidates := scope.FindRelationCandidatesNames(name)

	switch len(candidates) {
	case 0:
		tctx.ErrOnToken(ctx.GetStart(), "column %q does not found in any relation in the current scope, available relations: %s", name, strings.Join(scope.RelationNames(), ", "))
	case 1:
		return resolveQualifiedColumn(ctx, candidates[0], name, scope, tctx)
	default:
		tctx.ErrOnToken(ctx.GetStart(), "column reference %q is ambiguous, found in multiple relations: %s", name, strings.Join(candidates, ", "))
	}
	return nil
}

func resolveQualifiedColumn(ctx *parser.ColumnrefContext, relationName, name string, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	relation, table, ok := scope.FindRelation(relationName)
	if !ok {
		tctx.ErrOnToken(ctx.GetStart(), "relation %q does not found in the current scope, available: %s", relationName, strings.Join(scope.RelationNames(), ", "))
	}

	if name == "*" {
		return &ir.WildcardColumnExpr{
			Relation: relation,
		}
	}

	if !table.HasColumn(name) {
		tctx.ErrOnToken(ctx.GetStart(), "relation %q has no column %q, available: %s", relation.Name, name, strings.Join(table.ColumnNames(), ", "))
	}

	return &ir.ColumnExpr{Name: name, Relation: relation}
}

func parseDsql(dsqlExpr *parser.Dsql_paramContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	rawName := dsqlExpr.DSQL_PARAM().GetText()
	paramName := strings.TrimPrefix(rawName, "$")

	var field *string
	if attrName := dsqlExpr.Attr_name(); attrName != nil {
		f := attrName.GetText()
		field = &f
	}

	return &ir.DsqlExpr{
		Name:  paramName,
		Field: field,
	}
}
