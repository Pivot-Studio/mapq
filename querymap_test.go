package mapq

import (
	"testing"
)

func TestQueryMap(t *testing.T) {
	type args struct {
		data  map[string]interface{}
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"basic true case", args{map[string]interface{}{"a": 1}, "a==1"}, true, false},
		{"basic false case", args{map[string]interface{}{"a": 1}, "a==2"}, false, false},
		{"and true case", args{map[string]interface{}{"a": 1, "b": 2}, "a==1&&b==2"}, true, false},
		{"and false case", args{map[string]interface{}{"a": 1, "b": 2}, "a==1&&b==3"}, false, false},
		{"or true case", args{map[string]interface{}{"a": 1, "b": 2}, "a==1||b==2"}, true, false},
		{"or true case2", args{map[string]interface{}{"a": 1, "b": 2}, "a==1||b==3"}, true, false},
		{"or false case", args{map[string]interface{}{"a": 1, "b": 2}, "a<=0||b>3"}, false, false},
		{"nested case", args{map[string]interface{}{"a": 1, "b": 2, "c": map[string]interface{}{"d": 3}},
			"a==1&&b==2&&c.d==3"}, true, false},
		{"add case", args{map[string]interface{}{"a": 1, "b": 2}, "a+b==3"}, true, false},
		{"add case false", args{map[string]interface{}{"a": 1, "b": 2}, "a+b<3"}, false, false},
		{"mul case", args{map[string]interface{}{"a": 3, "b": 2}, "a*b==6"}, true, false},
		{"div case", args{map[string]interface{}{"a": 3, "b": 2}, "a/b==1.5"}, true, false},
		{"parentheses case", args{map[string]interface{}{"a": 1, "b": 2}, "a==1&&!(b==2||b==3)"}, false, false},
		{"parentheses mul case", args{map[string]interface{}{"a": 3, "b": 2}, "(a+b)*b==10"}, true, false},
		{"parentheses mul false case", args{map[string]interface{}{"a": 3, "b": 2}, "(a+b)*b<5"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryMap(tt.args.data, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("QueryMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
