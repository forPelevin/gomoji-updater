package unicodefile

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/forPelevin/gomoji"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestProviderAllEmojisSkipsMalformedRows(t *testing.T) {
	t.Helper()

	const (
		sequencesBody = `# group: Smileys & Emotion
# subgroup: face-smiling

1F600      ; fully-qualified     # 😀 E1.0 grinning face
1F601      fully-qualified       # 😁 E1.0 beaming face with smiling eyes
1F3FB      ; component           # 🏻 E2.0 light skin tone
`
		zwjBody = `# group: People & Body
# subgroup: person-gesture

1F44B 200D 1F3FB ; fully-qualified  # 👋🏻 E2.0 waving hand: light skin tone
`
		testBody = `# group: Animals & Nature
# subgroup: animal-mammal

1F43C      ; fully-qualified     # 🐼 E0.6 panda face
`
	)

	responses := map[string]string{
		"emoji-sequences":     sequencesBody,
		"emoji-zwj-sequences": zwjBody,
		"emoji-test":          testBody,
	}

	originalTransport := http.DefaultClient.Transport
	originalDefaultTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultClient.Transport = originalTransport
		http.DefaultTransport = originalDefaultTransport
	})

	handler := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		name := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
		if body, ok := responses[strings.TrimSuffix(name, ".txt")]; ok {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		}

		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		}, nil
	})

	http.DefaultClient.Transport = handler
	http.DefaultTransport = handler

	provider := NewProvider()
	emojis, err := provider.AllEmojis(context.Background())
	if err != nil {
		t.Fatalf("AllEmojis returned error: %v", err)
	}

	if len(emojis) != 3 {
		t.Fatalf("expected 3 emojis, got %d", len(emojis))
	}

	got := make(map[string]gomoji.Emoji)
	for _, em := range emojis {
		got[em.Slug] = em
	}

	if _, ok := got["grinning-face"]; !ok {
		t.Errorf("missing grinning-face entry: %+v", got)
	}
	if _, ok := got["waving-hand-light-skin-tone"]; !ok {
		t.Errorf("missing waving-hand-light-skin-tone entry: %+v", got)
	}
	if _, ok := got["panda-face"]; !ok {
		t.Errorf("missing panda-face entry: %+v", got)
	}

	if _, ok := got["beaming-face-with-smiling-eyes"]; ok {
		t.Errorf("unexpected emoji parsed from malformed line: %+v", got["beaming-face-with-smiling-eyes"])
	}
}
