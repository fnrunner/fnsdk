package fn

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// AsMain evaluates the ResourceContext from STDIN to STDOUT.
// `input` can be
// - a `ResourceContextProcessor` which implements `Process` method
// - a function `Runner` which implements `Run` method
func AsMain(input interface{}) error {
	err := func() error {
		var p ResourceContextProcessor
		switch input := input.(type) {
		case runnerProcessor:
			p = input
		case ResourceContextProcessorFunc:
			p = input
		default:
			return fmt.Errorf("unknown input type %T", input)
		}
		// read stdin
		in, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("unable to read from stdin: %v", err)
		}
		//return fmt.Errorf("received bytes: %s", string(in))
		out, err := Run(p, in)
		// If there is an error, we don't return the error immediately.
		// We write out to stdout before returning any error.
		_, outErr := os.Stdout.Write(out)
		if outErr != nil {
			return outErr
		}
		return err
	}()
	if err != nil {
		Logf("failed to evaluate function: %v", err)
	}
	return err
}

// Run evaluates the function. input must be a resourceContext in json format. An
// updated resourceContext will be returned.
func Run(p ResourceContextProcessor, in []byte) ([]byte, error) {
	// parse input as a resource context
	rCtx, err := ParseResourceContext(in)
	if err != nil {
		return nil, err
	}
	// calls the external program which implements the Run interface
	// based on how AsMain is defined it either calls
	// - the raw processor implementation (processor.go)
	// - the more abstracted sdk through the fn runner (runner_processor.go)
	success, fnErr := p.Process(rCtx)
	// marshal the renewed rctx to return as stdout
	out, jsonErr := json.MarshalIndent(rCtx, "", "  ")
	if jsonErr != nil {
		return out, jsonErr
	}
	if fnErr != nil {
		return out, fnErr
	}
	if !success {
		return out, fmt.Errorf("error: function failure")
	}
	return out, nil
}
