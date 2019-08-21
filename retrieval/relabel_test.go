package retrieval

import (
	"reflect"
	"regexp"
	"testing"

	clientmodel "github.com/prometheus/client_golang/model"

	"github.com/prometheus/prometheus/config"
)

func TestRelabel(t *testing.T) {
	tests := []struct {
		input   clientmodel.LabelSet
		relabel []config.DefaultedRelabelConfig
		output  clientmodel.LabelSet
	}{
		{
			input: clientmodel.LabelSet{
				"a": "foo",
				"b": "bar",
				"c": "baz",
			},
			relabel: []config.DefaultedRelabelConfig{
				{
					SourceLabels: clientmodel.LabelNames{"a"},
					Regex:        &config.Regexp{*regexp.MustCompile("f(.*)")},
					TargetLabel:  clientmodel.LabelName("d"),
					Separator:    ";",
					Replacement:  "ch${1}-ch${1}",
					Action:       config.RelabelReplace,
				},
			},
			output: clientmodel.LabelSet{
				"a": "foo",
				"b": "bar",
				"c": "baz",
				"d": "choo-choo",
			},
		},
		{
			input: clientmodel.LabelSet{
				"a": "foo",
				"b": "bar",
				"c": "baz",
			},
			relabel: []config.DefaultedRelabelConfig{
				{
					SourceLabels: clientmodel.LabelNames{"a", "b"},
					Regex:        &config.Regexp{*regexp.MustCompile("^f(.*);(.*)r$")},
					TargetLabel:  clientmodel.LabelName("a"),
					Separator:    ";",
					Replacement:  "b${1}${2}m", // boobam
					Action:       config.RelabelReplace,
				},
				{
					SourceLabels: clientmodel.LabelNames{"c", "a"},
					Regex:        &config.Regexp{*regexp.MustCompile("(b).*b(.*)ba(.*)")},
					TargetLabel:  clientmodel.LabelName("d"),
					Separator:    ";",
					Replacement:  "$1$2$2$3",
					Action:       config.RelabelReplace,
				},
			},
			output: clientmodel.LabelSet{
				"a": "boobam",
				"b": "bar",
				"c": "baz",
				"d": "boooom",
			},
		},
		{
			input: clientmodel.LabelSet{
				"a": "foo",
			},
			relabel: []config.DefaultedRelabelConfig{
				{
					SourceLabels: clientmodel.LabelNames{"a"},
					Regex:        &config.Regexp{*regexp.MustCompile("o$")},
					Action:       config.RelabelDrop,
				}, {
					SourceLabels: clientmodel.LabelNames{"a"},
					Regex:        &config.Regexp{*regexp.MustCompile("f(.*)")},
					TargetLabel:  clientmodel.LabelName("d"),
					Separator:    ";",
					Replacement:  "ch$1-ch$1",
					Action:       config.RelabelReplace,
				},
			},
			output: nil,
		},
		{
			input: clientmodel.LabelSet{
				"a": "foo",
			},
			relabel: []config.DefaultedRelabelConfig{
				{
					SourceLabels: clientmodel.LabelNames{"a"},
					Regex:        &config.Regexp{*regexp.MustCompile("no-match")},
					Action:       config.RelabelDrop,
				},
			},
			output: clientmodel.LabelSet{
				"a": "foo",
			},
		},
		{
			input: clientmodel.LabelSet{
				"a": "foo",
			},
			relabel: []config.DefaultedRelabelConfig{
				{
					SourceLabels: clientmodel.LabelNames{"a"},
					Regex:        &config.Regexp{*regexp.MustCompile("no-match")},
					Action:       config.RelabelKeep,
				},
			},
			output: nil,
		},
		{
			input: clientmodel.LabelSet{
				"a": "foo",
			},
			relabel: []config.DefaultedRelabelConfig{
				{
					SourceLabels: clientmodel.LabelNames{"a"},
					Regex:        &config.Regexp{*regexp.MustCompile("^f")},
					Action:       config.RelabelKeep,
				},
			},
			output: clientmodel.LabelSet{
				"a": "foo",
			},
		},
		{
			// No replacement must be applied if there is no match.
			input: clientmodel.LabelSet{
				"a": "boo",
			},
			relabel: []config.DefaultedRelabelConfig{
				{
					SourceLabels: clientmodel.LabelNames{"a"},
					Regex:        &config.Regexp{*regexp.MustCompile("^f")},
					TargetLabel:  clientmodel.LabelName("b"),
					Replacement:  "bar",
					Action:       config.RelabelReplace,
				},
			},
			output: clientmodel.LabelSet{
				"a": "boo",
			},
		},
	}

	for i, test := range tests {
		var relabel []*config.RelabelConfig
		for _, rl := range test.relabel {
			relabel = append(relabel, &config.RelabelConfig{rl})
		}
		res, err := Relabel(test.input, relabel...)
		if err != nil {
			t.Errorf("Test %d: error relabeling: %s", i+1, err)
		}

		if !reflect.DeepEqual(res, test.output) {
			t.Errorf("Test %d: relabel output mismatch: expected %#v, got %#v", i+1, test.output, res)
		}
	}
}
