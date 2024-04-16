package corednsyaegi

import (
	"fmt"
	"os"
	"regexp"

	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

// PluginFactoryAPI is the function signature that plugins must implement.
// This function should return the plugin instance.
// The method should be call `NewPlugin`.
type NewPluginAPISignature = func(next corednsplugin.Handler) corednsplugin.Handler

const pluginAPIName = "NewPlugin"

var packageRegexp = regexp.MustCompile(`(?m)^package +([^\s]+) *$`)

func LoadPlugin(src string) (NewPluginAPISignature, error) {
	// Load the plugin in a new interpreter.
	// For each plugin we need to use an independent interpreter to avoid name collisions.
	yaegiInterp, err := newPluginYaegiInterpreter()
	if err != nil {
		return nil, fmt.Errorf("could not create a new Yaegi interpreter: %w", err)
	}

	_, err = yaegiInterp.Eval(src)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate plugin source code: %w", err)
	}

	// Discover package name.
	packageMatch := packageRegexp.FindStringSubmatch(src)
	if len(packageMatch) != 2 {
		return nil, fmt.Errorf("invalid plugin source code, could not get package name")
	}
	packageName := packageMatch[1]

	// Get plugin logic.
	pluginFuncTmp, err := yaegiInterp.Eval(fmt.Sprintf("%s.%s", packageName, pluginAPIName))
	if err != nil {
		return nil, fmt.Errorf("could not get plugin: %w", err)
	}

	pluginFunc, ok := pluginFuncTmp.Interface().(NewPluginAPISignature)
	if !ok {
		return nil, fmt.Errorf("invalid plugin type")
	}

	return pluginFunc, nil
}

// newPluginReadyYaegiInterpreter will:
// - Create a new Yaegi interpreter.
// - Add the required libraries available (standard library and our own library).
func newPluginYaegiInterpreter() (*interp.Interpreter, error) {
	// Create interpreter
	i := interp.New(interp.Options{Env: os.Environ()})

	// Add standard library.
	err := i.Use(stdlib.Symbols)
	if err != nil {
		return nil, fmt.Errorf("yaegi could not use stdlib symbols: %w", err)
	}

	// Add unsafe library.
	err = i.Use(unsafe.Symbols)
	if err != nil {
		return nil, fmt.Errorf("yaegi could not use stdlib unsafe symbols: %w", err)
	}

	// Add our own plugin library.
	err = i.Use(Symbols)
	if err != nil {
		return nil, fmt.Errorf("yaegi could not use custom symbols: %w", err)
	}

	return i, nil
}
