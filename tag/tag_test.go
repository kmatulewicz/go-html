package tag

import (
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	type args struct {
		doc   string
		tag   string
		match []Check
	}
	tests := []struct {
		name string
		args args
		want *Tag
	}{
		{
			name: "1",
			args: args{doc: `<some attr-1 = cont1 attr_2='cont2"' attr3="cont with space'">`, tag: "some"},
			want: &Tag{Name: "some", Attr: map[string]string{"attr-1": "cont1", "attr_2": "cont2\"", "attr3": "cont with space'"}, ContentIndex: 62, AfterClosureIndex: -1},
		},
		{
			name: "2",
			args: args{doc: `<someother></someother><some attr1="cont1"
			attr2="cont2"	attr3="cont
			with	space"><someother>some text</someother>`, tag: "some"},
			want: &Tag{Name: "some", Attr: map[string]string{"attr1": "cont1", "attr2": "cont2", "attr3": "cont\n\t\t\twith\tspace"}, ContentIndex: 87, AfterClosureIndex: -1},
		},
		{
			name: "3",
			args: args{doc: `<some attr0 attr1="cont1" attr2="cont2" attr3="cont with space" attr4>`, tag: "some"},
			want: &Tag{Name: "some", Attr: map[string]string{"attr0": "", "attr1": "cont1", "attr2": "cont2", "attr3": "cont with space", "attr4": ""}, ContentIndex: 70, AfterClosureIndex: -1},
		},
		{
			name: "4",
			args: args{doc: `<some></some><some attr1="cont1" attr2="cont2" attr3="cont with space">`, tag: "some", match: []Check{Has("attr2"), Contains("attr3", "with"), Equal("attr1", "cont1")}},
			want: &Tag{Name: "some", Attr: map[string]string{"attr1": "cont1", "attr2": "cont2", "attr3": "cont with space"}, ContentIndex: 71, AfterClosureIndex: -1},
		},
		{
			name: "5",
			args: args{doc: `<some attr1="cont1" attr2="cont2" attr3="cont with space">`, tag: "some", match: []Check{Has("attr5")}},
			want: nil,
		},
		{
			name: "6",
			args: args{doc: `<some attr1="cont1" attr2="cont2" attr3="cont with space">`, tag: "some", match: []Check{Contains("attr3", "cont1")}},
			want: nil,
		},
		{
			name: "7",
			args: args{doc: `<some aTTr1="cont1" attr2="cont2" attr3="cont with space">`, tag: "some", match: []Check{Equal("attr1", "cont1"), Equal("attr1", "cont")}},
			want: nil,
		},
		{
			name: "8",
			args: args{doc: `<br/>`, tag: "br", match: []Check{}},
			want: &Tag{Name: "br", Attr: map[string]string{"/": ""}, ContentIndex: 5, AfterClosureIndex: -1},
		},
		{
			name: "9",
			args: args{doc: `<br />`, tag: "br", match: []Check{}},
			want: &Tag{Name: "br", Attr: map[string]string{"/": ""}, ContentIndex: 6, AfterClosureIndex: -1},
		},
	}
	for i := range tests {
		if tests[i].want != nil {
			tests[i].want.checks = tests[i].args.match
			tests[i].want.doc = tests[i].args.doc
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Find(tt.args.doc, tt.args.tag, tt.args.match); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContent(t *testing.T) {
	type args struct {
		doc   string
		tag   string
		match []Check
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{"<a><a></a></a>", "a", []Check{}}, "<a></a>"},
		{"2", args{"<a><a></a>", "a", []Check{}}, ""},
		{"3", args{"<a><a></a></a><a></a>", "a", []Check{}}, "<a></a>"},
		{"4", args{"<a>", "a", []Check{}}, ""},
		{"5", args{"<ab><a><ab></ab></a>", "a", []Check{}}, "<ab></ab>"},
		{"6", args{"<ab><a ></ab></a >", "a", []Check{}}, "</ab>"},
		{"7", args{"<ab><a ></ab></a", "a", []Check{}}, ""},
		{"8", args{"<ab><a><a><ab></ab></ab></a ></ab>", "a", []Check{}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Find(tt.args.doc, tt.args.tag, tt.args.match).Content(); got != tt.want {
				t.Errorf("Find().Content() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTag_Next(t *testing.T) {
	type fields struct {
		Name              string
		Attr              map[string]string
		ContentIndex      int
		AfterClosureIndex int
		doc               string
		checks            []Check
	}
	tests := []struct {
		name   string
		fields fields
		want   *Tag
	}{
		{
			"1",
			fields{
				"a",
				map[string]string{"id": "1"},
				8,
				28,
				"<a id=1><a id=2><a id=3></a></a></a>",
				[]Check{Has("id")},
			},
			&Tag{"a",
				map[string]string{"id": "2"},
				16,
				32,
				"<a id=1><a id=2><a id=3></a></a></a>",
				[]Check{Has("id")}},
		},
		{
			"2",
			fields{
				"a",
				map[string]string{"id": "1"},
				8,
				32,
				"<a id=1><br><a id=2><a id=3></a></a></a>",
				[]Check{Has("id")},
			},
			&Tag{"a",
				map[string]string{"id": "2"},
				20,
				36,
				"<a id=1><br><a id=2><a id=3></a></a></a>",
				[]Check{Has("id")}},
		},
		{
			"3",
			fields{
				"a",
				map[string]string{"id": "3"},
				28,
				40,
				"<a id=1><br><a id=2><a id=3></a></a></a>",
				[]Check{Has("id")},
			},
			nil,
		},
	}
	for i := range tests {
		if tests[i].want != nil {
			tests[i].want.checks = tests[i].fields.checks
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tag{
				Name:              tt.fields.Name,
				Attr:              tt.fields.Attr,
				ContentIndex:      tt.fields.ContentIndex,
				AfterClosureIndex: tt.fields.AfterClosureIndex,
				doc:               tt.fields.doc,
				checks:            tt.fields.checks,
			}
			got := tr.Next()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tag.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}
