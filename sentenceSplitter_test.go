package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSentenceSplitter(t *testing.T) {
	splitter, _ := newSentenceSplitter()
	actual, _ := splitter.split(`人魚は、南の方の海にばかり棲んでいるのではあ
	りません。北の海にも棲んでいたのであります。
	　北方の海うみの色は、青うございました。ある
	とき、岩の上に、女の人魚があがって、あたりの景
	色をながめながら休んでいました。

	小川未明作 赤い蝋燭と人魚より`)

	if diff := cmp.Diff(
		[]string{
			"人魚は、南の方の海にばかり棲んでいるのではありません。",
			"北の海にも棲んでいたのであります。",
			"北方の海うみの色は、青うございました。",
			"あるとき、岩の上に、女の人魚があがって、あたりの景色をながめながら休んでいました。",
			"小川未明作赤い蝋燭と人魚より",
		},
		actual,
	); diff != "" {
		t.Errorf(diff)
	}
}
