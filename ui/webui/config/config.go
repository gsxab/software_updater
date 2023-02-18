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

package config

type Config struct {
	Secret     *string `json:"secret,omitempty"`
	Addr       string  `json:"addr,omitempty"`
	DateFormat string  `json:"date_format"`
	TimeFormat string  `json:"time_format"`
}

var WebUIConfig *Config

func DefaultConfig() *Config {
	return &Config{
		Secret:     nil,
		Addr:       "127.0.0.1:8080",
		DateFormat: "2006-01-02",
		TimeFormat: "2006-01-02 03:04:05",
	}
}
