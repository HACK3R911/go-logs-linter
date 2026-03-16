package detector

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	LoggerSlog LoggerType = iota
	LoggerZap
)

type LogCall struct {
	Logger        LoggerType
	Method        string
	Messages      []LogMessage
	Concatenation ConcatenationInfo
}

type LogMessage struct {
	Text     string
	Position token.Pos
}

type ConcatenationInfo struct {
	IsConcatenation bool
	VarNames        []string
	ZapKeys         []string
}

type LoggerType int

type LoggerDetector struct {
	slogIdent   string
	zapIdent    string
	slogMethods map[string]bool
	zapMethods  map[string]bool
	impPaths    map[string]bool
	files       []*ast.File
}

func NewLoggerDetector() *LoggerDetector {
	return &LoggerDetector{
		slogIdent: "slog",
		zapIdent:  "zap",
		slogMethods: map[string]bool{
			"Debug":        true,
			"Info":         true,
			"Warn":         true,
			"Error":        true,
			"DebugContext": true,
			"InfoContext":  true,
			"WarnContext":  true,
			"ErrorContext": true,
		},
		zapMethods: map[string]bool{
			"Debug":  true,
			"Info":   true,
			"Warn":   true,
			"Error":  true,
			"Debugw": true,
			"Infow":  true,
			"Warnw":  true,
			"Errorw": true,
		},
		impPaths: make(map[string]bool),
	}
}

func (d *LoggerDetector) Detect(pass *analysis.Pass, call *ast.CallExpr) *LogCall {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return nil
	}

	methodName := sel.Sel.Name

	d.initImpPaths(pass)

	if d.isSlogImport(ident) {
		return d.detectSlogCall(sel, call)
	}

	if d.isZapImport(ident) {
		return d.detectZapCall(call, methodName)
	}

	if d.isLogMethod(methodName) {
		if d.isZapPackageUsed() {
			return d.detectZapCall(call, methodName)
		}
	}

	return nil
}

func (d *LoggerDetector) initImpPaths(pass *analysis.Pass) {
	d.impPaths = make(map[string]bool)

	if pass.Pkg != nil && pass.Pkg.Imports() != nil {
		for _, imp := range pass.Pkg.Imports() {
			if imp != nil {
				d.impPaths[imp.Path()] = true
			}
		}
		return
	}

	if pass.Files != nil {
		for _, file := range pass.Files {
			if file != nil && file.Imports != nil {
				for _, imp := range file.Imports {
					if imp != nil && imp.Path != nil {
						path := strings.Trim(imp.Path.Value, `"`)
						d.impPaths[path] = true
					}
				}
			}
		}
	}
}

func (d *LoggerDetector) isSlogImport(ident *ast.Ident) bool {
	if d.impPaths["log/slog"] || d.impPaths["golang.org/x/exp/slog"] {
		if ident.Name == d.slogIdent || ident.Name == "_" {
			return true
		}
	}
	return false
}

func (d *LoggerDetector) isZapImport(ident *ast.Ident) bool {
	if d.impPaths["go.uber.org/zap"] {
		if ident.Name == d.zapIdent || ident.Name == "_" {
			return true
		}
	}
	return false
}

func (d *LoggerDetector) isZapPackageUsed() bool {
	return d.impPaths["go.uber.org/zap"]
}

func (d *LoggerDetector) isLogMethod(methodName string) bool {
	return d.slogMethods[methodName] || d.zapMethods[methodName]
}

func (d *LoggerDetector) detectSlogCall(sel *ast.SelectorExpr, call *ast.CallExpr) *LogCall {
	method := sel.Sel.Name
	if !d.slogMethods[method] {
		return nil
	}

	messages := d.extractSlogMessages(call)
	if len(messages) == 0 {
		return nil
	}

	concat := d.extractConcatenation(call)

	return &LogCall{
		Logger:        LoggerSlog,
		Method:        method,
		Messages:      messages,
		Concatenation: concat,
	}
}

func (d *LoggerDetector) extractConcatenation(call *ast.CallExpr) ConcatenationInfo {
	var info ConcatenationInfo

	if len(call.Args) == 0 {
		return info
	}

	firstArg := call.Args[0]
	if binExpr, ok := firstArg.(*ast.BinaryExpr); ok && binExpr.Op == token.ADD {
		info.IsConcatenation = true
		info.VarNames = append(info.VarNames, d.extractVarNames(binExpr)...)
	}

	return info
}

func (d *LoggerDetector) extractVarNames(expr ast.Expr) []string {
	var names []string

	switch e := expr.(type) {
	case *ast.BinaryExpr:
		names = append(names, d.extractVarNames(e.X)...)
		names = append(names, d.extractVarNames(e.Y)...)
	case *ast.Ident:
		names = append(names, e.Name)
	case *ast.ParenExpr:
		names = append(names, d.extractVarNames(e.X)...)
	}

	return names
}

func (d *LoggerDetector) extractSlogMessages(call *ast.CallExpr) []LogMessage {
	var messages []LogMessage

	if len(call.Args) == 0 {
		return messages
	}

	firstArg := call.Args[0]
	if lit, ok := firstArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		msg := strings.Trim(lit.Value, `"`)
		messages = append(messages, LogMessage{
			Text:     msg,
			Position: lit.Pos(),
		})
		return messages
	}

	if binExpr, ok := firstArg.(*ast.BinaryExpr); ok && binExpr.Op == token.ADD {
		if lit, ok := binExpr.X.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			msg := strings.Trim(lit.Value, `"`)
			messages = append(messages, LogMessage{
				Text:     msg,
				Position: lit.Pos(),
			})
		}
		return messages
	}

	if len(call.Args) >= 2 {
		secondArg := call.Args[1]
		if lit, ok := secondArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			msg := strings.Trim(lit.Value, `"`)
			messages = append(messages, LogMessage{
				Text:     msg,
				Position: lit.Pos(),
			})
		}
	}

	return messages
}

func (d *LoggerDetector) detectZapCall(call *ast.CallExpr, method string) *LogCall {
	if !d.zapMethods[method] {
		return nil
	}

	messages := d.extractZapMessages(call, method)
	if len(messages) == 0 {
		return nil
	}

	concat := ConcatenationInfo{}
	if strings.HasSuffix(method, "w") {
		concat.ZapKeys = d.extractZapKeys(call)
	}

	return &LogCall{
		Logger:        LoggerZap,
		Method:        method,
		Messages:      messages,
		Concatenation: concat,
	}
}

func (d *LoggerDetector) extractZapKeys(call *ast.CallExpr) []string {
	var keys []string

	if len(call.Args) < 3 {
		return keys
	}

	for i := 1; i < len(call.Args)-1; i += 2 {
		if lit, ok := call.Args[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
			key := strings.Trim(lit.Value, `"`)
			keys = append(keys, key)
		}
	}

	return keys
}

func (d *LoggerDetector) extractZapMessages(call *ast.CallExpr, method string) []LogMessage {
	var messages []LogMessage

	isInfowMethod := strings.HasSuffix(method, "w")

	if isInfowMethod && len(call.Args) > 0 {
		if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
			msg := strings.Trim(lit.Value, `"`)
			messages = append(messages, LogMessage{
				Text:     msg,
				Position: lit.Pos(),
			})
		}
	}

	if !isInfowMethod && len(call.Args) > 0 {
		if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
			msg := strings.Trim(lit.Value, `"`)
			messages = append(messages, LogMessage{
				Text:     msg,
				Position: lit.Pos(),
			})
		}
	}

	return messages
}
