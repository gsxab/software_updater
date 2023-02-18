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

package hook

type PositionCommand int

const (
	FirstCmd PositionCommand = iota
	LastCmd
	PrevCmd
	NextCmd
)

type Position struct {
	Cmd PositionCommand
	Ref string
}

const (
	FirstStr = "first"
	LastStr  = "last"
	PrevStr  = "prev_to"
	NextStr  = "next_to"
)

var PositionCommandMap = map[PositionCommand]string{
	FirstCmd: FirstStr,
	LastCmd:  LastStr,
	PrevCmd:  PrevStr,
	NextCmd:  NextStr,
}

func First() *Position {
	return &Position{FirstCmd, ""}
}

func Last() *Position {
	return &Position{FirstCmd, ""}
}

func PrevTo(ref string) *Position {
	return &Position{PrevCmd, ref}
}

func NextTo(ref string) *Position {
	return &Position{NextCmd, ref}
}
