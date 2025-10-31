package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mhirii/rest-template/internal/dto"
	"github.com/mhirii/rest-template/internal/service"
	"github.com/rs/zerolog"
)

// ExampleHandlers contains all handler functions for Example endpoints.
type ExampleHandlers struct {
	Store *service.ExampleStore
}

// RegisterExampleRoutes registers all Example REST endpoints to the API.
func RegisterExampleRoutes(api huma.API, store *service.ExampleStore) {
	h := &ExampleHandlers{Store: store}

	// POST /examples - Create a new Example
	huma.Post(api, "/examples", h.Create)

	// PUT /examples/{id} - Replace an Example
	huma.Put(api, "/examples/{id}", h.Replace)

	// PATCH /examples/{id} - Update fields of an Example
	huma.Patch(api, "/examples/{id}", h.Update)

	// GET /examples/{id} - Get a single Example
	huma.Get(api, "/examples/{id}", h.Get)

	// GET /examples - List all Examples
	huma.Get(api, "/examples", h.List)

	// DELETE /examples/{id} - Delete an Example
	huma.Delete(api, "/examples/{id}", h.Delete)
}

// Create handles POST /examples
func (h *ExampleHandlers) Create(ctx context.Context, input *struct {
	Body dto.Example
}) (*struct {
	Body dto.Example
}, error) {
	e := h.Store.Create(input.Body)
	l := zerolog.Ctx(ctx)
	l.Info().Str("id", e.ID).Msg("Example created")
	return &struct{ Body dto.Example }{Body: e}, nil
}

// Replace handles PUT /examples/{id}
func (h *ExampleHandlers) Replace(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body dto.Example
}) (*struct {
	Body dto.Example
}, error) {
	e := h.Store.Replace(input.ID, input.Body)
	return &struct{ Body dto.Example }{Body: e}, nil
}

// Update handles PATCH /examples/{id}
func (h *ExampleHandlers) Update(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body struct {
		Name *string `json:"name,omitempty"`
	}
}) (*struct {
	Body dto.Example
}, error) {
	e, ok := h.Store.Update(input.ID, input.Body.Name)
	l := zerolog.Ctx(ctx)
	if !ok {
		l.Warn().Str("id", input.ID).Msg("Example not found for update")
		return nil, huma.Error404NotFound("not found")
	}
	l.Info().Str("id", e.ID).Msg("Example updated")
	return &struct{ Body dto.Example }{Body: e}, nil
}

// Get handles GET /examples/{id}
func (h *ExampleHandlers) Get(ctx context.Context, input *struct {
	ID string `path:"id"`
}) (*struct {
	Body dto.Example
}, error) {
	e, ok := h.Store.Get(input.ID)
	l := zerolog.Ctx(ctx)
	if !ok {
		l.Warn().Str("id", input.ID).Msg("Example not found for get")
		return nil, huma.Error404NotFound("not found")
	}
	l.Info().Str("id", e.ID).Msg("Example retrieved")
	return &struct{ Body dto.Example }{Body: e}, nil
}

// List handles GET /examples
func (h *ExampleHandlers) List(ctx context.Context, input *struct{}) (*struct {
	Body []dto.Example
}, error) {
	list := h.Store.List()
	l := zerolog.Ctx(ctx)
	l.Info().Int("count", len(list)).Msg("Examples listed")
	return &struct{ Body []dto.Example }{Body: list}, nil
}

// Delete handles DELETE /examples/{id}
func (h *ExampleHandlers) Delete(ctx context.Context, input *struct {
	ID string `path:"id"`
}) (*struct{}, error) {
	ok := h.Store.Delete(input.ID)
	l := zerolog.Ctx(ctx)
	if !ok {
		l.Warn().Str("id", input.ID).Msg("Example not found for delete")
		return nil, huma.Error404NotFound("not found")
	}
	l.Info().Str("id", input.ID).Msg("Example deleted")
	return &struct{}{}, nil
}
