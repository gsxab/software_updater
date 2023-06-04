/*
 * SPDX-License-Identifier: GPL-3.0-or-later
 *
 * Copyright (c) 2023. gsxab.
 *
 * This file is part of Software Update Watcher, a.k.a. Zhixin Robot.
 *
 * Software Update Watcher is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
 *
 * Software Update Watcher is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with Software Update Watcher. If not, see <https://www.gnu.org/licenses/>.
 */

package action

import (
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/tools/web"
)

type Args struct {
	Elements web.Elements
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

func IndexedElementToArgs(element selenium.WebElement, index int, input *Args) *Args {
	elements := []selenium.WebElement{element}
	elements = append(elements, input.Elements[:index]...)
	elements = append(elements, element)
	elements = append(elements, input.Elements[index+1:]...)
	ret := &Args{
		Elements: elements,
	}
	if input != nil {
		ret.Strings = input.Strings
		ret.Update = input.Update
	}
	return ret
}

func AnotherElementToArgs(element selenium.WebElement, input *Args) *Args {
	elements := input.Elements
	elements = append(elements, element)
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

func IndexedStringToArgs(str string, index int, input *Args) *Args {
	var strings []string
	strings = append(strings, input.Strings[:index]...)
	strings = append(strings, str)
	strings = append(strings, input.Strings[index+1:]...)
	ret := &Args{
		Strings: strings,
	}
	if input != nil {
		ret.Strings = input.Strings
		ret.Update = input.Update
	}
	return ret
}

func AnotherStringToArgs(str string, input *Args) *Args {
	strings := input.Strings
	strings = append(strings, str)
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
