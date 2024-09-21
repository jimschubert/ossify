package conventions

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/jimschubert/ossify/internal/model"
)

func Test_indexOf(t *testing.T) {
	type args struct {
		data   []string
		search string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "first index", args: args{data: []string{"asdf", "jkl;"}, search: "asdf"}, want: 0},
		{name: "last index", args: args{data: []string{"asdf", "jkl;"}, search: "jkl;"}, want: 1},
		{name: "other index", args: args{data: []string{"asdf", "aaaa", "bbbb", "cccc", "jkl;"}, search: "bbbb"}, want: 2},
		{name: "not found", args: args{data: []string{"asdf", "aaaa", "bbbb", "cccc"}, search: "jkl;"}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Index(tt.args.data, tt.args.search); got != tt.want {
				t.Errorf("indexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRule_MarshalJSON(t *testing.T) {
	type fields struct {
		Level model.StrictnessLevel
		Type  model.RuleType
		Value string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"valid input 0/0/.bak", fields{model.StrictnessLevel(0), model.RuleType(0), ".bak"}, `{"level":"prohibited","type":"unspecified","value":".bak"}`, false},
		{"valid input 1/1/src", fields{model.StrictnessLevel(1), model.RuleType(1), "src"}, `{"level":"optional","type":"directory","value":"src"}`, false},
		{"valid input 2/2/LICENSE", fields{model.StrictnessLevel(2), model.RuleType(2), "LICENSE"}, `{"level":"preferred","type":"file","value":"LICENSE"}`, false},
		{"valid input mixed/other", fields{model.StrictnessLevel(3), model.RuleType(1), "tools"}, `{"level":"required","type":"directory","value":"tools"}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &model.Rule{
				Level: tt.fields.Level,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			got, err := r.MarshalJSON()
			actual := string(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rule.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr {
				return
			}
			if !reflect.DeepEqual(actual, tt.want) {
				t.Errorf("Rule.MarshalJSON() = %v, want %v", actual, tt.want)
			}
		})
	}
}

func TestRule_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Level model.StrictnessLevel
		Type  model.RuleType
		Value string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"valid input 0/0/.bak", fields{model.StrictnessLevel(0), model.RuleType(0), ".bak"}, args{[]byte(`{"level":"prohibited","type":"unspecified","value":".bak"}`)}, false},
		{"invalid input 0/0/.bak", fields{model.StrictnessLevel(1), model.RuleType(1), ".bak"}, args{[]byte(`{"level":"unspecified","type":"unspecified","value":".bak"}`)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &model.Rule{
				Level: tt.fields.Level,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			unMarshaled := &model.Rule{}
			err := json.Unmarshal(tt.args.data, unMarshaled)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rule.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr {
				return
			}
			if !reflect.DeepEqual(unMarshaled, r) {
				t.Errorf("Rule.UnmarshalJSON() = %v, want %v", unMarshaled, r)
			}
		})
	}
}

// func TestLoad(t *testing.T) {
//	tests := []struct {
//		name    string
//		want    *[]Convention
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := Load()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Load() = %v, want %v", got, tt.want)
//			}
//		})
//	}
// }
//
// func TestConvention_Evaluate(t *testing.T) {
//	type fields struct {
//		Name  string
//		Rules []Rule
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Convention{
//				Name:  tt.fields.Name,
//				Rules: tt.fields.Rules,
//			}
//			if err := c.Evaluate(); (err != nil) != tt.wantErr {
//				t.Errorf("Convention.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
// }
