package engine

import (
	"fmt"
	"software_updater/core/action"
	"software_updater/core/hook"
)

type ActionManager struct {
	categories ActionTrie
}

func NewActionManager() *ActionManager {
	actionManager := &ActionManager{}
	actionManager.categories = NewActionTrie()
	return actionManager
}

func (m *ActionManager) Register(factory action.Factory) bool {
	path := factory.Path()

	_, err := m.categories.PutFactLeaf(path, factory)
	if err != nil {
		return false
	}

	return true
}

func (m *ActionManager) RegisterHook(info *hook.RegisterInfo) error {
	position := info.Position
	if position == nil {
		position = &hook.Position{Cmd: hook.LastCmd}
	}
	err := m.categories.PutHook(info.Action, info.Event, info.Hook, position)
	return err
}

func (m *ActionManager) Action(storedAction *StoredAction) (action.Action, []*hook.ActionHooks, error) {
	path := action.Path(storedAction.Path)
	args := storedAction.JSON
	if storedAction.Path == nil {
		path = m.categories.GetPath(storedAction.Name)
	}
	factory, hooks, err := m.categories.SearchLeafAllHooks(path)
	if err != nil {
		return nil, nil, fmt.Errorf("action not found, path: %s, error: %w", path, err)
	}
	if len(args) == 0 {
		args = "{}"
	}
	a, err := factory.NewAction(args)
	if err != nil {
		return nil, nil, fmt.Errorf("action creation failed, path: %s, error: %w", path, err)
	}
	return a, hooks, err
}

type StoredAction struct {
	Path []string `json:"path,omitempty"`
	Name string   `json:"name,omitempty"`
	JSON string   `json:"json,omitempty"`
}

type StoredBranch struct {
	Actions []StoredAction `json:"actions,omitempty"`
	Next    []StoredBranch `json:"next,omitempty"`
}

type StoredFlow struct {
	Root StoredBranch `json:"root,omitempty"`
}
