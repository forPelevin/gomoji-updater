package unicodefile

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/forPelevin/gomoji"
)

type Provider struct{}

// NewProvider creates a new instance of Provider.
func NewProvider() *Provider {
	return &Provider{}
}

// AllEmojis downloads emoji data from unicode.org text feeds.
func (o *Provider) AllEmojis(ctx context.Context) ([]gomoji.Emoji, error) {
	sequences, err := o.allEmojis(ctx, "emoji-sequences")
	if err != nil {
		return nil, fmt.Errorf("get emoji-sequences: %w", err)
	}

	zwjSequences, err := o.allEmojis(ctx, "emoji-zwj-sequences")
	if err != nil {
		return nil, fmt.Errorf("get emoji-zwj-sequences: %w", err)
	}

	testEmojis, err := o.allEmojis(ctx, "emoji-test")
	if err != nil {
		return nil, fmt.Errorf("get emoji-test: %w", err)
	}

	allEmojis := append(sequences, zwjSequences...)
	allEmojis = append(allEmojis, testEmojis...)

	return allEmojis, nil
}

func (o *Provider) allEmojis(_ context.Context, name string) ([]gomoji.Emoji, error) {
	resp, err := http.Get(fmt.Sprintf("https://unicode.org/Public/emoji/latest/%s.txt", name))
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	var (
		emojis   []gomoji.Emoji
		group    string
		subgroup string
	)

	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		row := strings.TrimSpace(sc.Text())
		if row == "" {
			continue
		}

		if strings.HasPrefix(row, "# group:") {
			group = strings.TrimSpace(strings.ReplaceAll(row, "# group:", ""))
			continue
		}
		if strings.HasPrefix(row, "# subgroup:") {
			subgroup = strings.TrimSpace(strings.ReplaceAll(row, "# subgroup:", ""))
			continue
		}

		if strings.HasPrefix(row, "#") {
			continue
		}

		codePoint := strings.TrimSpace(before(row, ";"))
		row = strings.Split(row, ";")[1]
		status := strings.TrimSpace(before(row, "#"))
		if status == "component" {
			continue
		}

		emoChar := strings.TrimSpace(between(row, "#", " E"))
		unicodeName := strings.TrimSpace(strings.ReplaceAll(after(row, "#"), emoChar, ""))
		slug := strings.TrimLeft(strings.TrimLeft(unicodeName, "E"), "01234567890.")
		slug = strings.ReplaceAll(strings.TrimSpace(slug), ":", "")
		slug = strings.ReplaceAll(strings.TrimSpace(slug), " ", "-")
		slug = strings.ToLower(slug)

		emojis = append(emojis, gomoji.Emoji{
			Slug:        slug,
			Character:   emoChar,
			UnicodeName: unicodeName,
			CodePoint:   codePoint,
			Group:       group,
			SubGroup:    subgroup,
		})
	}

	return emojis, nil
}

func between(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func before(value string, a string) string {
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func after(value string, a string) string {
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}
