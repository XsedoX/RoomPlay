package response

import "github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"

type JsonPatch struct {
	Op    string `json:"op" example:"replace"`
	Path  string `json:"path" example:"/name"`
	Value any    `json:"value,omitempty" example:"New name"`
}

type JsonPatchResponseBuilder struct {
	patches []JsonPatch
}

func NewJsonPatchResponseBuilder() *JsonPatchResponseBuilder {
	return &JsonPatchResponseBuilder{}
}

func (b *JsonPatchResponseBuilder) Add(path string, value any) *JsonPatchResponseBuilder {
	b.patches = append(b.patches, JsonPatch{
		Op:    "add",
		Path:  path,
		Value: value,
	})
	return b
}

func (b *JsonPatchResponseBuilder) Replace(path string, value any) *JsonPatchResponseBuilder {
	b.patches = append(b.patches, JsonPatch{
		Op:    "replace",
		Path:  path,
		Value: value,
	})
	return b
}

func (b *JsonPatchResponseBuilder) Remove(path string) *JsonPatchResponseBuilder {
	b.patches = append(b.patches, JsonPatch{
		Op:   "remove",
		Path: path,
	})
	return b
}

func (b *JsonPatchResponseBuilder) Test(path string, value any) *JsonPatchResponseBuilder {
	b.patches = append(b.patches, JsonPatch{
		Op:    "test",
		Path:  path,
		Value: value,
	})
	return b
}

func (b *JsonPatchResponseBuilder) Build(actionName websocket_action.WebSocketOutgoingAction) WebSocketSuccess[[]JsonPatch] {
	copyJustInCase := make([]JsonPatch, len(b.patches))
	copy(copyJustInCase, b.patches)
	return WebSocketSuccess[[]JsonPatch]{
		ActionName: actionName,
		Success: Success[[]JsonPatch]{
			Data: copyJustInCase,
		},
	}
}
