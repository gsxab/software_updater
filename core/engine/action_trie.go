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

package engine

import (
	"fmt"
	"software_updater/core/action"
	"software_updater/core/hook"
	"strings"
)

type ActionTrie interface {
	Search(action.Path) (action.Factory, *hook.ActionHooks, map[string]ActionTrieNode, error)
	SearchLeaf(action.Path) (action.Factory, *hook.ActionHooks, error)
	SearchLeafAllHooks(action.Path) (action.Factory, []*hook.ActionHooks, error)
	PutHook(action.Path, hook.Event, hook.Hook, *hook.Position) error
	PutHookLeaf(action.Path, hook.Event, hook.Hook, *hook.Position) error
	PutFactLeaf(action.Path, action.Factory) (action.Factory, error)
	GetPath(string) action.Path
	DTO() (*action.HierarchyDTO, error)
}

func NewActionTrie() ActionTrie {
	return &ActionTrieImpl{
		ActionTrieNode: &ActionTrieInternalNode{
			children: make(map[string]ActionTrieNode),
		},
		namePath: make(map[string]action.Path),
	}
}

type ActionTrieImpl struct {
	ActionTrieNode
	namePath map[string]action.Path
}

func (a *ActionTrieImpl) getNode(path action.Path) (node ActionTrieNode, err error) {
	node = a.ActionTrieNode
	for _, key := range path {
		node, err = node.Child(key)
		if err != nil {
			return
		}
	}
	return
}

func (a *ActionTrieImpl) getNodeAllowAll(path action.Path) (node ActionTrieNode, err error) {
	node = a.ActionTrieNode
	for _, key := range path {
		if key == action.All {
			return
		}
		node, err = node.Child(key)
		if err != nil {
			return
		}
	}
	return
}

func (a *ActionTrieImpl) getNodeAllHooks(path action.Path) (node ActionTrieNode, hooks []*hook.ActionHooks, err error) {
	node = a.ActionTrieNode
	hooks = append(hooks, node.Hooks())
	for _, key := range path {
		node, err = node.Child(key)
		if err != nil {
			return
		}
		hooks = append(hooks, node.Hooks())
	}
	return
}

func (a *ActionTrieImpl) getOrCreateLeafNode(path action.Path) (node ActionTrieNode, err error) {
	node = a.ActionTrieNode
	var newNode ActionTrieNode
	last := len(path) - 1
	for idx, key := range path[:last] {
		newNode, err = node.Child(key)
		if err != nil {
			newNode = &ActionTrieInternalNode{
				ActionTrieNodeBase: ActionTrieNodeBase{
					path: strings.Join(path[:idx+1], "."),
				},
				children: make(map[string]ActionTrieNode),
			}
			err := node.AddChild(key, newNode)
			if err != nil {
				return nil, err
			}
		}
		node = newNode
	}
	leaf, err := node.Child(path[last])
	if err != nil {
		newNode = &ActionTrieLeaf{
			ActionTrieNodeBase: ActionTrieNodeBase{
				path: path.String(),
			},
		}
		err := node.AddChild(path[last], newNode)
		if err != nil {
			return nil, err
		}
		return newNode, nil
	}
	return leaf, nil
}

func (a *ActionTrieImpl) Search(keys action.Path) (action.Factory, *hook.ActionHooks, map[string]ActionTrieNode, error) {
	node, err := a.getNode(keys)
	if err != nil {
		return nil, nil, nil, err
	}
	if node.IsLeaf() {
		return node.Content(), node.Hooks(), nil, nil
	}
	return nil, node.Hooks(), node.Children(), nil
}

func (a *ActionTrieImpl) SearchLeaf(keys action.Path) (action.Factory, *hook.ActionHooks, error) {
	node, err := a.getNode(keys)
	if err != nil {
		return nil, nil, err
	}
	return node.Content(), node.Hooks(), nil
}

func (a *ActionTrieImpl) SearchLeafAllHooks(keys action.Path) (action.Factory, []*hook.ActionHooks, error) {
	node, hooks, err := a.getNodeAllHooks(keys)
	if err != nil {
		return nil, nil, err
	}
	return node.Content(), hooks, nil
}

func (a *ActionTrieImpl) PutHook(keys action.Path, event hook.Event, h hook.Hook, pos *hook.Position) error {
	node, err := a.getNodeAllowAll(keys)
	if err != nil {
		return err
	}
	err = node.Hooks().PutAt(event, h, pos)
	if err != nil {
		return err
	}
	return nil
}

func (a *ActionTrieImpl) PutHookLeaf(keys action.Path, event hook.Event, h hook.Hook, pos *hook.Position) error {
	node, err := a.getOrCreateLeafNode(keys)
	if err != nil {
		return err
	}
	err = node.Hooks().PutAt(event, h, pos)
	if err != nil {
		return err
	}
	return nil
}

func (a *ActionTrieImpl) PutFactLeaf(keys action.Path, factory action.Factory) (action.Factory, error) {
	node, err := a.getOrCreateLeafNode(keys)
	if err != nil {
		return nil, err
	}
	old := node.SetContent(factory)
	a.namePath[keys.Name()] = keys
	return old, nil
}

func (a *ActionTrieImpl) GetPath(name string) action.Path {
	return a.namePath[name]
}

func (a *ActionTrieImpl) DTO() (*action.HierarchyDTO, error) {
	return a.SubTreeDTO(0), nil
}

type ActionTrieNode interface {
	Path() string
	SetPath(path string) error
	Hooks() *hook.ActionHooks
	IsLeaf() bool
	Content() action.Factory
	SetContent(action.Factory) action.Factory
	Children() map[string]ActionTrieNode
	Child(key string) (ActionTrieNode, error)
	AddChild(key string, node ActionTrieNode) error
	SubTreeDTO(int) *action.HierarchyDTO
}

type ActionTrieNodeBase struct {
	path  string
	hooks hook.ActionHooks
}

func (a *ActionTrieNodeBase) Path() string {
	return a.path
}

func (a *ActionTrieNodeBase) SetPath(path string) error {
	if len(a.path) > 0 {
		return fmt.Errorf("path already set: %s", a.path)
	}
	a.path = path
	return nil
}

func (a *ActionTrieNodeBase) Hooks() *hook.ActionHooks {
	return &a.hooks
}

type ActionTrieInternalNode struct {
	ActionTrieNodeBase
	children map[string]ActionTrieNode
}

func (a *ActionTrieInternalNode) IsLeaf() bool {
	return false
}

func (a *ActionTrieInternalNode) Content() action.Factory {
	return nil
}

func (a *ActionTrieInternalNode) SetContent(action.Factory) action.Factory {
	return nil
}

func (a *ActionTrieInternalNode) Children() map[string]ActionTrieNode {
	return a.children
}

func (a *ActionTrieInternalNode) Child(key string) (ActionTrieNode, error) {
	node, ok := a.children[key]
	if !ok {
		return nil, fmt.Errorf("child not found: %s", key)
	}
	return node, nil
}

func (a *ActionTrieInternalNode) AddChild(key string, node ActionTrieNode) error {
	if _, ok := a.children[key]; ok {
		return fmt.Errorf("duplicated key: %s.%s", a.path, key)
	}
	a.children[key] = node
	return nil
}

func (a *ActionTrieInternalNode) SubTreeDTO(level int) *action.HierarchyDTO {
	result := &action.HierarchyDTO{
		Name:  a.path[strings.LastIndex(a.path, action.Delim)+1:],
		Path:  a.Path(),
		Level: level,
		Leaf:  false,
	}
	children := make([]*action.HierarchyDTO, 0, len(a.children))
	for _, child := range a.children {
		children = append(children, child.SubTreeDTO(level+1))
	}
	result.Children = children
	return result
}

type ActionTrieLeaf struct {
	ActionTrieNodeBase
	factory action.Factory
}

func (a *ActionTrieLeaf) IsLeaf() bool {
	return true
}

func (a *ActionTrieLeaf) Content() action.Factory {
	return a.factory
}

func (a *ActionTrieLeaf) SetContent(factory action.Factory) action.Factory {
	tmp := a.factory
	a.factory = factory
	return tmp
}

func (a *ActionTrieLeaf) Children() map[string]ActionTrieNode {
	panic("accessing a leaf as an internal node")
}

func (a *ActionTrieLeaf) Child(string) (ActionTrieNode, error) {
	panic("accessing a leaf as an internal node")
}

func (a *ActionTrieLeaf) AddChild(string, ActionTrieNode) error {
	panic("accessing a leaf as an internal node")
}

func (a *ActionTrieLeaf) SubTreeDTO(level int) *action.HierarchyDTO {
	result := &action.HierarchyDTO{
		Name:  a.path[strings.LastIndex(a.path, action.Delim)+1:],
		Path:  a.path,
		Level: level,
		Leaf:  true,
	}
	return result
}
