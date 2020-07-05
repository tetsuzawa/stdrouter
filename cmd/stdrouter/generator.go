package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/tetsuzawa/stdrouter/internal/stdrouter"
)

const stdrouterPkg = "github.com/tetsuzawa/stdrouter"

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) writeTpl(t *template.Template, data interface{}) error {
	if err := t.Execute(&g.buf, data); err != nil {
		return fmt.Errorf("failed to execute template -> %w", err)
	}
	return nil
}

func (g *Generator) generateHeadMsg() error {
	tplName := "head message"
	t, err := template.New(tplName).Parse(TplHeadMsg)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, strings.Join(os.Args, " "))
}

func (g *Generator) generatePackage(name string) error {
	tplName := "package"
	t, err := template.New(tplName).Parse(TplPackage)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, name)
}

func (g *Generator) generateImport() error {
	tplName := "import"
	t, err := template.New(tplName).Parse(TplImport)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, nil)
}

func (g *Generator) generateImportImpl(name string) error {
	tplName := "import content"
	t, err := template.New(tplName).Parse(TplImpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, strconv.Quote(name))
}

func (g *Generator) generateClosingBracket() error {
	tplName := "closing bracket"
	t, err := template.New(tplName).Parse(TplClosingBracket)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, nil)
}

func (g *Generator) generateRouter(routerInstanceName string) error {
	tplName := "router"
	t, err := template.New(tplName).Parse(TplRouter)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, routerInstanceName)
}

func (g *Generator) generateHandlerFunc(funcName string, pathParams []string) error {
	tplName := "handler func"
	t, err := template.New(tplName).Parse(TplHandlerFunc)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	data := struct {
		FuncName   string
		PathParams []string
	}{
		FuncName:   funcName,
		PathParams: pathParams,
	}

	return g.writeTpl(t, data)
}

func (g *Generator) generateSeparatePath(n int) error {
	tplName := "separate path"
	t, err := template.New(tplName).Parse(TplSeparatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	data := struct {
		Base string
		Num  int
		Tail string
	}{
		Base: "p",
		Num:  n,
		Tail: "p",
	}
	return g.writeTpl(t, data)
}

func (g *Generator) generateSeparateParam(n int) error {
	tplName := "separate param"
	t, err := template.New(tplName).Parse(TplSeparatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	data := struct {
		Base string
		Num  int
		Tail string
	}{
		Base: "endpoint",
		Num:  n,
		Tail: "param",
	}
	return g.writeTpl(t, data)
}

func (g *Generator) generateSwitch(target string) error {
	tplName := "switch"
	t, err := template.New(tplName).Parse(TplSwitch)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, target)
}

func (g *Generator) generateCasePath(path string) error {
	tplName := "case path"
	t, err := template.New(tplName).Parse(TplCase)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, path)
}

func (g *Generator) generateClosingCurlyBraces() error {
	tplName := "closing curly braces"
	t, err := template.New(tplName).Parse(TplClosingCurlyBraces)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, nil)
}

func (g *Generator) generateCaseMethod(httpMethod string) error {
	tplName := "case method"
	t, err := template.New(tplName).Parse(TplCase)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, "http.Method"+strings.Title(strings.ToLower(httpMethod)))
}

func (g *Generator) generateFunc(handlerFunc stdrouter.HandlerFunc, args []string) error {
	tplName := "function"
	t, err := template.New(tplName).Parse(TplImpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	sargs := "(w, r"
	for _, p := range args {
		sargs = fmt.Sprintf("%s, %s", sargs, p)
	}
	sargs += ")"

	if handlerFunc.Package == "" {
		if err = g.writeTpl(t, handlerFunc.Func+sargs); err != nil {
			return err
		}
		return nil
	}
	return g.writeTpl(t, handlerFunc.Package+"."+handlerFunc.Func+sargs)
}

func (g *Generator) generateDefault() error {
	tplName := "default"
	t, err := template.New(tplName).Parse(TplDefault)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, nil)
}

func (g *Generator) generateIf(expr string) error {
	tplName := "if"
	t, err := template.New(tplName).Parse(TplIf)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, expr)
}

func (g *Generator) generateElse() error {
	tplName := "else"
	t, err := template.New(tplName).Parse(TplElse)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, nil)
}

func (g *Generator) generateSeparatePathFunc() error {
	tplName := "separate path function"
	t, err := template.New(tplName).Parse(TplSeparatePathFunc)
	if err != nil {
		return fmt.Errorf("failed to parse template: `%s` -> %w", tplName, err)
	}
	return g.writeTpl(t, nil)
}

func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to Analyze the error")
		return g.buf.Bytes()
	}
	return src
}

func (g *Generator) Generate(cfg *AnalyzerConfig) error {
	var err error

	// generate headers
	if err = g.generateHeadMsg(); err != nil {
		return fmt.Errorf("generateHeadMsg -> %w", err)
	}
	if err = g.generatePackage(cfg.PackageName); err != nil {
		return fmt.Errorf("generatePackage -> %w", err)
	}
	if err = g.generateImport(); err != nil {
		return fmt.Errorf("generateImport -> %w", err)
	}

	// use in SeparatePath func
	cfg.ImportedPkgs = append(cfg.ImportedPkgs, "path", "strings")
	// Drop stdrouter package
	cfg.ImportedPkgs = stdrouter.Drop(stdrouterPkg, cfg.ImportedPkgs)
	// Drop duplication
	cfg.ImportedPkgs = stdrouter.DropDuplication(cfg.ImportedPkgs)
	for _, v := range cfg.ImportedPkgs {
		if err = g.generateImportImpl(v); err != nil {
			return fmt.Errorf("generateImportImpl -> %w", err)
		}
	}
	if err = g.generateClosingBracket(); err != nil {
		return fmt.Errorf("generateClosingBracket -> %w", err)
	}
	if err = g.generateRouter(cfg.RouterInstanceName); err != nil {
		return fmt.Errorf("generateRouter -> %w", err)
	}

	// find hierarchy to path parameter and max depth
	var nodeHierarchies []stdrouter.Node
	sentinelNode := &stdrouter.Node{Depth: 0}
	stdrouter.Walk(cfg.Node, func(node *stdrouter.Node) bool {
		if node.IsPathParam {
			nodeHierarchies = append(nodeHierarchies, *node)
		}
		if node.Depth > sentinelNode.Depth {
			sentinelNode.Depth = node.Depth
		}
		return true
	})
	nodeHierarchies = append(nodeHierarchies, *sentinelNode)

	// generate functions
	handlerFuncName := "Base"
	var pathParams []string
	var startingNode = cfg.Node
	for i, pathParamBaseNode := range nodeHierarchies {
		if err = g.generateHandlerFunc("handle"+handlerFuncName, pathParams); err != nil {
			return fmt.Errorf("generateHandlerFunc -> %w", err)
		}
		if i == 0 {
			if err = g.generateSeparatePath(pathParamBaseNode.Depth); err != nil {
				return fmt.Errorf("generateSeparatePath -> %w", err)
			}
		} else {
			if err = g.generateSeparatePath(pathParamBaseNode.Depth - nodeHierarchies[i-1].Depth); err != nil {
				return fmt.Errorf("generateSeparatePath -> %w", err)
			}
		}

		if err = g.generateSwitch("endpoint"); err != nil {
			return fmt.Errorf("generateSwitch -> %w", err)
		}

		var pathParamHandler *stdrouter.HandlerFunc
		// generate switches of router
		stdrouter.Walk(startingNode, func(node *stdrouter.Node) bool {
			if node.Depth > pathParamBaseNode.Depth {
				return false
			}
			p := "/"
			if node != startingNode {
				p = path.Clean(stdrouter.BuildBasePath(node) + path.Clean("/"+node.Endpoint))
			}
			if node.Endpoint == pathParamBaseNode.Endpoint {
				pathParamHandler = &stdrouter.HandlerFunc{
					Func: "handle" + stdrouter.SnakeToCamel(pathParamBaseNode.Endpoint),
				}
				// to generate the function with name "handle<Param>"
				handlerFuncName = stdrouter.SnakeToCamel(node.Endpoint)
				startingNode = node
				pathParams = append(pathParams, stdrouter.ToLowerFirstLetter(stdrouter.SnakeToCamel(node.Endpoint)))
				return true
			}

			if i != 0 {
				p = strings.Replace(p, stdrouter.BuildBasePath(node), "", 1)
			}
			if err = g.generateCasePath(strconv.Quote(p)); err != nil {
				err = fmt.Errorf("generateCasePath -> %w", err)
				return false
			}
			if err = g.generateSwitch("r.Method"); err != nil {
				err = fmt.Errorf("generateSwitch -> %w", err)
				return false
			}
			for httpMethod, handlerFunc := range node.Methods {
				if err = g.generateCaseMethod(httpMethod); err != nil {
					err = fmt.Errorf("generateCaseMethod -> %w", err)
					return false
				}
				if err = g.generateFunc(handlerFunc, pathParams); err != nil {
					err = fmt.Errorf("generateFunc -> %w", err)
					return false
				}
			}
			if err = g.generateDefault(); err != nil {
				err = fmt.Errorf("generateDefault -> %w", err)
				return false
			}
			if err = g.generateFunc(*cfg.MethodNotAllowedHandler, nil); err != nil {
				err = fmt.Errorf("generateFunc -> %w", err)
				return false
			}
			if err = g.generateClosingCurlyBraces(); err != nil {
				err = fmt.Errorf("generateClosingCurlyBraces -> %w", err)
				return false
			}
			return true
		})
		if err != nil {
			return fmt.Errorf("stdrouter.Walk -> %w", err)
		}

		if err = g.generateDefault(); err != nil {
			return fmt.Errorf("generateDefault -> %w", err)
		}
		if pathParamHandler != nil {
			if i == 0 {
				if err = g.generateSeparateParam(pathParamBaseNode.Depth - 1); err != nil {
					return fmt.Errorf("generateSeparateParam -> %w", err)
				}
			} else {
				if err = g.generateSeparateParam(pathParamBaseNode.Depth - nodeHierarchies[i-1].Depth - 1); err != nil {
					return fmt.Errorf("generateSeparateParam -> %w", err)
				}
			}
			if err = g.generateIf(fmt.Sprintf(`endpoint == "%s"`, stdrouter.BuildBasePath(&pathParamBaseNode))); err != nil {
				return fmt.Errorf("generateIf -> %w", err)
			}
			args := append([]string{"p"}, pathParams[:len(pathParams)-1]...)
			args = append(args, "param[1:]")
			if err = g.generateFunc(*pathParamHandler, args); err != nil {
				return fmt.Errorf("generateFunc -> %w", err)
			}
			if err = g.generateElse(); err != nil {
				return fmt.Errorf("generateElse -> %w", err)
			}
			if err = g.generateFunc(*cfg.NotFoundHandler, nil); err != nil {
				return fmt.Errorf("generateFunc -> %w", err)
			}
			if err = g.generateClosingCurlyBraces(); err != nil {
				return fmt.Errorf("generateClosingCurlyBraces -> %w", err)
			}
		} else {
			if err = g.generateFunc(*cfg.NotFoundHandler, nil); err != nil {
				return fmt.Errorf("generateFunc -> %w", err)
			}
		}
		// end switch
		if err = g.generateClosingCurlyBraces(); err != nil {
			return fmt.Errorf("generateClosingCurlyBraces -> %w", err)
		}
		// end func
		if err = g.generateClosingCurlyBraces(); err != nil {
			return fmt.Errorf("generateClosingCurlyBraces -> %w", err)
		}
	}

	// generate helper func
	if err = g.generateSeparatePathFunc(); err != nil {
		return fmt.Errorf("generateSeparatePathFunc -> %w", err)
	}

	return nil
}
