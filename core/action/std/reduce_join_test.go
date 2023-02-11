package std

import (
	"github.com/tebeka/selenium"
	"software_updater/core/action"
)

import (
	"context"
	"reflect"
	"software_updater/core/db/po"
	"sync"
	"testing"
)

func TestReduceJoin_Do(t *testing.T) {
	type fields struct {
		Val string
	}
	type args struct {
		ctx     context.Context
		driver  selenium.WebDriver
		input   *action.Args
		version *po.Version
		wg      *sync.WaitGroup
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput *action.Args
		wantExit   action.Result
		wantErr    bool
	}{
		{
			name:   "reduce_join",
			fields: fields{Val: `{"sep":","}`},
			args: args{
				ctx:     nil,
				driver:  nil,
				input:   &action.Args{Strings: []string{"1", "2"}},
				version: nil,
				wg:      nil,
			},
			wantOutput: &action.Args{Strings: []string{"1,2"}},
			wantExit:   0,
			wantErr:    false,
		},
		{
			name:   "reduce_join",
			fields: fields{Val: `{"sep":",","skip":[0]}`},
			args: args{
				ctx:     nil,
				driver:  nil,
				input:   &action.Args{Strings: []string{"1", "2"}},
				version: nil,
				wg:      nil,
			},
			wantOutput: &action.Args{Strings: []string{"2", "1"}},
			wantExit:   0,
			wantErr:    false,
		},
		{
			name:   "reduce_join",
			fields: fields{Val: `{"sep":",","skip":[1]}`},
			args: args{
				ctx:     nil,
				driver:  nil,
				input:   &action.Args{Strings: []string{"1", "2", "3"}},
				version: nil,
				wg:      nil,
			},
			wantOutput: &action.Args{Strings: []string{"1,3", "2"}},
			wantExit:   0,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := (&ReduceJoin{}).NewAction(tt.fields.Val)
			gotOutput, gotExit, err := a.Do(tt.args.ctx, tt.args.driver, tt.args.input, tt.args.version, tt.args.wg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Do() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
			if gotExit != tt.wantExit {
				t.Errorf("Do() gotExit = %v, want %v", gotExit, tt.wantExit)
			}
		})
	}
}
