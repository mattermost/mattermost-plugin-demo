package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

//go:embed fixtures/mm_blocks_examples.json
var mmBlocksExamplesJSON []byte

type mmBlocksExample struct {
	Message string                `json:"message"`
	Props   model.StringInterface `json:"props"`
}

var mmBlocksExamples map[string]mmBlocksExample

func init() {
	if err := json.Unmarshal(mmBlocksExamplesJSON, &mmBlocksExamples); err != nil {
		panic("failed to parse mm_blocks examples: " + err.Error())
	}
}

func mmBlocksExampleNames() []string {
	names := make([]string, 0, len(mmBlocksExamples))
	for name := range mmBlocksExamples {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (p *Plugin) executeMmBlocksExample(args *model.CommandArgs, name string) *model.CommandResponse {
	if name == "" {
		var b strings.Builder
		b.WriteString("###### mm_blocks examples\n")
		for _, n := range mmBlocksExampleNames() {
			fmt.Fprintf(&b, "- `%s` — %s\n", n, mmBlocksExamples[n].Message)
		}
		b.WriteString("\nUse `/mm_blocks example <name>` to post one.")
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         b.String(),
		}
	}

	example, ok := mmBlocksExamples[name]
	if !ok {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown mm_blocks example: `%s`. Available: `%s`", name, strings.Join(mmBlocksExampleNames(), "`, `")),
		}
	}

	props := model.StringInterface{}
	for k, v := range example.Props {
		props[k] = v
	}
	ensureMmBlocksExampleActions(props)

	post := &model.Post{
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		UserId:    p.botID,
		Message:   example.Message,
		Props:     props,
	}

	if _, err := p.API.CreatePost(post); err != nil {
		const errorMessage = "Failed to create mm_blocks example post"
		p.API.LogError(errorMessage, "err", err.Error(), "example", name)
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         errorMessage,
		}
	}

	return &model.CommandResponse{}
}

// ensureMmBlocksExampleActions fills mm_blocks_actions so every referenced
// action_id (except disabled controls) has an entry. E2E fixtures often omit
// this map; the server rejects those posts.
func ensureMmBlocksExampleActions(props model.StringInterface) {
	blocks, ok := props["mm_blocks"].([]any)
	if !ok {
		return
	}

	referenced := collectMmBlocksActionIDs(blocks)
	existing, _ := props["mm_blocks_actions"].(map[string]any)

	actions := make(map[string]any, len(referenced))
	for id := range referenced {
		if existing != nil {
			if entry, ok := existing[id]; ok {
				actions[id] = entry
				continue
			}
		}
		actions[id] = mmBlocksPostAction("/mm_blocks_integration", map[string]any{"fixture_action": id})
	}

	if len(actions) == 0 {
		delete(props, "mm_blocks_actions")
		return
	}
	props["mm_blocks_actions"] = actions
}

func collectMmBlocksActionIDs(blocks []any) map[string]struct{} {
	ids := map[string]struct{}{}
	for _, block := range blocks {
		collectMmBlocksActionIDsFromValue(ids, block)
	}
	return ids
}

func collectMmBlocksActionIDsFromValue(ids map[string]struct{}, value any) {
	switch v := value.(type) {
	case []any:
		for _, item := range v {
			collectMmBlocksActionIDsFromValue(ids, item)
		}
	case map[string]any:
		disabled, _ := v["disabled"].(bool)

		switch v["type"] {
		case "button", "static_select":
			if !disabled {
				addMmBlocksActionID(ids, v["action_id"])
			}
		case "select":
			if !disabled {
				addMmBlocksActionID(ids, v["data_source_action"])
				addMmBlocksActionID(ids, v["onChange"])
			}
		case "text_input", "bool_input", "date_input", "datetime_input", "file_input":
			if !disabled {
				addMmBlocksActionID(ids, v["onChange"])
			}
		case "container", "collapsible":
			collectMmBlocksActionIDsFromValue(ids, v["header"])
			collectMmBlocksActionIDsFromValue(ids, v["content"])
		case "column_set":
			collectMmBlocksActionIDsFromValue(ids, v["columns"])
		case "column":
			collectMmBlocksActionIDsFromValue(ids, v["items"])
		}
	}
}

func addMmBlocksActionID(ids map[string]struct{}, raw any) {
	id, ok := raw.(string)
	if !ok || id == "" {
		return
	}
	ids[id] = struct{}{}
}
