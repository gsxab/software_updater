package action

import (
	"fmt"
	"github.com/tebeka/selenium"
)

type Args struct {
	Elements []selenium.WebElement
	Strings  []string
	Update   bool
}

const (
	Any     int = -1
	Same    int = -2
	OneMore int = -3
)

func ElementsToArgs(elements []selenium.WebElement, input *Args) *Args {
	ret := &Args{
		Elements: elements,
	}
	if input != nil {
		ret.Strings = input.Strings
		ret.Update = input.Update
	}
	return ret
}

func ElementToArgs(element selenium.WebElement, input *Args) *Args {
	ret := &Args{
		Elements: []selenium.WebElement{element},
	}
	if input != nil {
		ret.Strings = input.Strings
		ret.Update = input.Update
	}
	return ret
}

func AnotherElementToArgs(element selenium.WebElement, input *Args) *Args {
	elements := []selenium.WebElement{element}
	elements = append(elements, input.Elements...)
	ret := &Args{
		Elements: elements,
	}
	if input != nil {
		ret.Strings = input.Strings
		ret.Update = input.Update
	}
	return ret
}

func StringsToArgs(strings []string, input *Args) *Args {
	ret := &Args{
		Strings: strings,
	}
	if input != nil {
		ret.Elements = input.Elements
		ret.Update = input.Update
	}
	return ret
}

func StringToArgs(str string, input *Args) *Args {
	ret := &Args{
		Strings: []string{str},
	}
	if input != nil {
		ret.Elements = input.Elements
		ret.Update = input.Update
	}
	return ret
}

func AnotherStringToArgs(str string, input *Args) *Args {
	strings := []string{str}
	strings = append(strings, input.Strings...)
	ret := &Args{
		Strings: strings,
	}
	if input != nil {
		ret.Elements = input.Elements
		ret.Update = input.Update
	}
	return ret
}

func MarkUpdateToArgs(input *Args) *Args {
	ret := &Args{
		Update: true,
	}
	if input != nil {
		ret.Elements = input.Elements
		ret.Strings = input.Strings
	}
	return ret
}

func StaticCheckInput(inSpec int, outSpec int, input int, id string) (output int, err error) {
	switch inSpec {
	case Any:
		break
	default:
		if inSpec != input {
			return 0, fmt.Errorf("static input check failed for %s, input specification: %d, input: %d", id, inSpec, input)
		}
	}
	switch outSpec {
	case Any:
		return Any, nil
	case Same:
		return input, nil
	case OneMore:
		if input == Any {
			return Any, nil
		}
		return input + 1, nil
	default:
		return outSpec, nil
	}
}

func DynamicCheckInput(inSpec int, input int, id string) error {
	switch inSpec {
	case Any:
		return nil
	default:
		if inSpec != input {
			return fmt.Errorf("dynamic input check failed for %s, input specification: %d, input: %d", id, inSpec, input)
		}
		return nil
	}
}

func DynamicCheckOutput(outSpec int, input int, output int, id string) error {
	switch outSpec {
	case Any:
		return nil
	case Same:
		if input != output {
			return fmt.Errorf("dynamic output check failed for %s, output specification: same, input: %d, output: %d", id, input, output)
		}
	case OneMore:
		if input+1 != output {
			return fmt.Errorf("dynamic output check failed for %s, output specification: one_more, input: %d, output: %d", id, input, output)
		}
	default:
		if outSpec != output {
			return fmt.Errorf("dynamic output check failed for %s, output specification: %d, output: %d", id, outSpec, output)
		}
	}
	return nil
}
