package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type Action interface {
	Path() Path
	Icon() string
	InElmNum() int
	OutElmNum() int
	InStrNum() int
	OutStrNum() int
	Init(ctx context.Context, wg *sync.WaitGroup) error
	Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error)
	ToDTO() *DTO
}

type Factory interface {
	Path() Path
	Icon() string
	NewAction(args string) (Action, error)
	ToProtoDTO() *ProtoDTO
}

type ProtoDTO struct {
	Name     string   `json:"name"`
	Icon     string   `json:"icon"`
	ReadPage bool     `json:"read_page,omitempty"`
	OpenPage bool     `json:"open_page,omitempty"`
	Input    []string `json:"input,omitempty"`  // not used
	Output   []string `json:"output,omitempty"` // not used
}

type DTO struct {
	*ProtoDTO
	Values map[string]string `json:"values,omitempty"`
}

type HierarchyDTO struct {
	Name     string          `json:"name"`
	Path     string          `json:"path"`
	Level    int             `json:"level"`
	Leaf     bool            `json:"leaf"`
	Children []*HierarchyDTO `json:"children,omitempty"`
}

func ToDTO(action Action) *DTO {
	dto := action.ToDTO()
	if dto.ProtoDTO == nil {
		dto.ProtoDTO = &ProtoDTO{}
	}
	if len(dto.Name) == 0 {
		dto.Name = action.Path().Name()
	}
	if len(dto.Icon) == 0 {
		dto.Icon = action.Icon()
	}
	return dto
}

func ToProtoDTO(factory Factory) *ProtoDTO {
	dto := factory.ToProtoDTO()
	if dto == nil {
		dto = &ProtoDTO{}
	}
	if len(dto.Name) == 0 {
		dto.Name = factory.Path().Name()
	}
	if len(dto.Icon) == 0 {
		dto.Icon = factory.Icon()
	}
	return dto
}

type Result int

const (
	Finished   Result = iota // action exited with success or error (default)
	Cancelled                // action cancelled before exiting
	Again                    // action exited, and needs running again
	StopBranch               // action exited, and needs current flow to stop
	StopFlow                 // action exited, and needs current flow to stop
	Skipped                  // action exited, for checking actions, no need to check
)
