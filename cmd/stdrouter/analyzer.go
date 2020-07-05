package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"
	"strings"

	"github.com/tetsuzawa/stdrouter/internal/stdrouter"
)

type AnalyzerConfig struct {
	fset                    *token.FileSet
	Node                    *stdrouter.Node
	ImportedPkgs            []string
	NotFoundHandler         *stdrouter.HandlerFunc
	MethodNotAllowedHandler *stdrouter.HandlerFunc
	PackageName             string
	RouterInstanceName      string
}

func Analyze(filename string) (*AnalyzerConfig, error) {
	cfg := &AnalyzerConfig{Node:new(stdrouter.Node)}
	cfg.fset = token.NewFileSet()
	var err error
	f, err := parser.ParseFile(cfg.fset, filename, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file -> %w", err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.File:
			SetPackageName(v, cfg)
		case *ast.GenDecl:
			if err = SetImportedPkg(v, cfg); err != nil {
				err = fmt.Errorf("SetImportedPkg -> %w", err)
				return false
			}
		case *ast.FuncDecl:
			if err = CheckFuncDecl(v); err != nil {
				err = fmt.Errorf("CheckFuncDecl -> %w", err)
				return false
			}
		case *ast.AssignStmt:
			if err = SetRouterInstance(v, cfg); err != nil {
				err = fmt.Errorf("SetRouterInstance -> %w", err)
				return false
			}
		case *ast.ExprStmt:
			if err = RegisterHandler(v, cfg); err != nil {
				err = fmt.Errorf("RegisterHandler -> %w", err)
				return false
			}
		default:
			return true
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func SetPackageName(file *ast.File, cfg *AnalyzerConfig) {
	cfg.PackageName = file.Name.Name
}

// SetImportedPkg adds imported packages to AnalyzerConfig. It does not handle CONST, TYPE and VAR.
func SetImportedPkg(genDecl *ast.GenDecl, cfg *AnalyzerConfig) error {
	if genDecl.Tok != token.IMPORT {
		return nil
	}
	for _, spec := range genDecl.Specs {
		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			return nil
		}
		name, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			return fmt.Errorf("strconv.Unquote -> %w", err)
		}
		cfg.ImportedPkgs = append(cfg.ImportedPkgs, name)
	}
	return nil
}

func CheckFuncDecl(funcDecl *ast.FuncDecl) error {
	funcName := funcDecl.Name.Name
	if funcName != "NewRouter" {
		return fmt.Errorf("invalid function declaration. want: NewRouter, got: %v", funcName)
	}
	return nil
}

func SetRouterInstance(assignStmt *ast.AssignStmt, cfg *AnalyzerConfig) error {
	routerIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return fmt.Errorf("syntax error: %s", cfg.fset.Position(assignStmt.Lhs[0].Pos()))
	}

	callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return nil
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	packageIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return fmt.Errorf("syntax error: %s", cfg.fset.Position(selectorExpr.X.Pos()))
	}
	if packageIdent.Name != "stdrouter" {
		return nil
	}
	if selectorExpr.Sel.Name != "NewRouter" {
		return nil
	}
	cfg.RouterInstanceName = routerIdent.Name
	return nil
}

func RegisterHandler(exprStmt *ast.ExprStmt, cfg *AnalyzerConfig) error {
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil
	}
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	routerIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return fmt.Errorf("syntax error: %s", cfg.fset.Position(selectorExpr.X.Pos()))
	}
	// Check if the Ident is router instance
	if routerIdent.Name != cfg.RouterInstanceName {
		return nil
	}
	methodName := selectorExpr.Sel.Name
	switch methodName {
	case "HandleFunc":
		if err := RegisterHandleFunc(callExpr.Args, cfg); err != nil {
			return fmt.Errorf("RegisterHandleFunc -> %w", err)
		}
	case "HandleNotFound":
		if err := RegisterHandleNotFound(callExpr.Args, cfg); err != nil {
			return fmt.Errorf("RegisterHandleNotFound -> %w", err)
		}
	case "HandleMethodNotAllowed":
		if err := RegisterHandleMethodNotAllowed(callExpr.Args, cfg); err !=nil {
			return fmt.Errorf("RegisterHandleMethodNotAllowed -> %w", err)
		}
	default:
		return fmt.Errorf("unknown method called: %s", methodName)
	}
	return nil
}

func RegisterHandleFunc(args []ast.Expr, cfg *AnalyzerConfig) error {
	if len(args) != 3 {
		return fmt.Errorf("invalid number of arguments to HandleFunc. got %d, want 1", len(args))
	}
	// check path
	basicLit, ok := args[0].(*ast.BasicLit)
	if !ok {
		return fmt.Errorf("the first argument is not BasicLit")
	}
	if basicLit.Kind != token.STRING {
		return fmt.Errorf("the type of first argument is invalid. want: string, got: %s", strings.ToLower(basicLit.Kind.String()))
	}
	path, err := strconv.Unquote(basicLit.Value)
	if err != nil {
		return fmt.Errorf("strconv.Unquote -> %w", err)
	}

	// check method
	methodSelectorExpr, ok := args[1].(*ast.SelectorExpr)
	if !ok {
		return fmt.Errorf("method must be chosen from the http package")
	}
	ident, ok := methodSelectorExpr.X.(*ast.Ident)
	if !ok {
		return fmt.Errorf("SelectorExpr.X is not Ident")
	}
	if ident.Name != "http" {
		return fmt.Errorf("method must be chosen from the http package")
	}
	httpMethod := methodSelectorExpr.Sel.Name
	if !stdrouter.Contains(httpMethod, stdrouter.HTTPMethods) {
		return fmt.Errorf("method not found. got: %v", httpMethod)
	}
	httpMethod = strings.TrimLeft(httpMethod, "Method")

	// check handler func
	var handlerIdent *ast.Ident
	var packageName, funcName string
	handlerSelectorExpr, ok := args[2].(*ast.SelectorExpr)
	if ok {
		handlerIdent, ok = handlerSelectorExpr.X.(*ast.Ident)
		if !ok {
			return nil
		}
		packageName = handlerIdent.Name
		funcName = handlerSelectorExpr.Sel.Name
	} else {
		handlerIdent, ok = args[2].(*ast.Ident)
		if !ok {
			return nil
		}
		packageName = ""
		funcName = handlerIdent.Name
	}
	if err := cfg.Node.Add(path, httpMethod, stdrouter.HandlerFunc{packageName, funcName}); err != nil {
		return fmt.Errorf("Node.Add ->")
	}
	return nil
}

func RegisterHandleNotFound(args []ast.Expr, cfg *AnalyzerConfig) error {
	if len(args) != 1 {
		log.Fatalf("invalid number of arguments to HandleNotFound. got %d, want 1", len(args))
	}

	var handlerIdent *ast.Ident
	var packageName, funcName string
	handlerSelectorExpr, ok := args[0].(*ast.SelectorExpr)
	if ok {
		handlerIdent, ok = handlerSelectorExpr.X.(*ast.Ident)
		if !ok {
			return nil
		}
		packageName = handlerIdent.Name
		funcName = handlerSelectorExpr.Sel.Name
	} else {
		handlerIdent, ok = args[0].(*ast.Ident)
		if !ok {
			return nil
		}
		packageName = ""
		funcName = handlerIdent.Name
	}
	if cfg.NotFoundHandler != nil {
		return fmt.Errorf("duplicate declaration: HandleNotFound")
	}
	cfg.NotFoundHandler = &stdrouter.HandlerFunc{packageName, funcName}
	return nil
}

func RegisterHandleMethodNotAllowed(args []ast.Expr, cfg *AnalyzerConfig) error {
	if len(args) != 1 {
		log.Fatalf("invalid number of arguments to HandleMethodNotAllowed. got %d, want 1", len(args))
	}
	var handlerIdent *ast.Ident
	var packageName, funcName string
	handlerSelectorExpr, ok := args[0].(*ast.SelectorExpr)
	if ok {
		handlerIdent, ok = handlerSelectorExpr.X.(*ast.Ident)
		if !ok {
			return nil
		}
		packageName = handlerIdent.Name
		funcName = handlerSelectorExpr.Sel.Name
	} else {
		handlerIdent, ok = args[0].(*ast.Ident)
		if !ok {
			return nil
		}
		packageName = ""
		funcName = handlerIdent.Name
	}
	if cfg.MethodNotAllowedHandler != nil {
		return fmt.Errorf("duplicate declaration: MethodNotAllowed")
	}
	cfg.MethodNotAllowedHandler = &stdrouter.HandlerFunc{packageName, funcName}
	return nil
}
