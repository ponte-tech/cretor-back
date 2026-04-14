package service

import (
	"context"

	"github.com/ponte-tech/cretor-back/modules/auth/domain"
	"github.com/ponte-tech/cretor-back/modules/auth/dto"
	"go.uber.org/zap"
)

type EmpreendimentoService struct {
	empreendimentoRepo domain.EmpreendimentoRepository
	construtoraRepo    domain.ConstrutoraRepository
	unidadeRepo        domain.UnidadeRepository
	logger             *zap.Logger
}

func NewEmpreendimentoService(
	empreendimentoRepo domain.EmpreendimentoRepository,
	construtoraRepo domain.ConstrutoraRepository,
	unidadeRepo domain.UnidadeRepository,
	logger *zap.Logger,
) *EmpreendimentoService {
	return &EmpreendimentoService{
		empreendimentoRepo: empreendimentoRepo,
		construtoraRepo:    construtoraRepo,
		unidadeRepo:        unidadeRepo,
		logger:             logger,
	}
}

func buildPaginatedResponse(cards []dto.EmpreendimentoCardResponse, total int64, page, pageSize int) *dto.PaginatedEmpreendimentosResponse {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	return &dto.PaginatedEmpreendimentosResponse{
		Data: cards, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages,
	}
}

// ListEmpreendimentos returns paginated card responses.
// construtora_nome is denormalized in the document — zero additional queries.
func (s *EmpreendimentoService) ListEmpreendimentos(ctx context.Context, tenantID string, filter domain.EmpreendimentoFilter) (*dto.PaginatedEmpreendimentosResponse, error) {
	empreendimentos, total, err := s.empreendimentoRepo.List(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}

	cards := make([]dto.EmpreendimentoCardResponse, len(empreendimentos))
	for i, e := range empreendimentos {
		cards[i] = dto.ToEmpreendimentoCardResponse(&e, e.ConstrutoraNome)
	}

	return buildPaginatedResponse(cards, total, filter.Pagination.Page, filter.Pagination.PageSize), nil
}

// SearchEmpreendimentos uses MongoDB text search for relevance-ranked results.
func (s *EmpreendimentoService) SearchEmpreendimentos(ctx context.Context, tenantID, query string, page, pageSize int) (*dto.PaginatedEmpreendimentosResponse, error) {
	empreendimentos, total, err := s.empreendimentoRepo.Search(ctx, tenantID, query, page, pageSize)
	if err != nil {
		return nil, err
	}

	cards := make([]dto.EmpreendimentoCardResponse, len(empreendimentos))
	for i, e := range empreendimentos {
		cards[i] = dto.ToEmpreendimentoCardResponse(&e, e.ConstrutoraNome)
	}

	return buildPaginatedResponse(cards, total, page, pageSize), nil
}

// GetEmpreendimentoDetail returns the full detail with construtora.
func (s *EmpreendimentoService) GetEmpreendimentoDetail(ctx context.Context, tenantID string, id string) (*dto.EmpreendimentoResponse, error) {
	emp, err := s.empreendimentoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := dto.ToEmpreendimentoResponse(emp)

	c, err := s.construtoraRepo.FindByID(ctx, emp.ConstrutoraID)
	if err != nil {
		s.logger.Warn("construtora not found", zap.String("id", emp.ConstrutoraID), zap.Error(err))
	} else {
		cr := dto.ToConstrutoraResponse(c)
		resp.Construtora = &cr
	}

	return &resp, nil
}

// ListUnidades returns unidades for an empreendimento.
func (s *EmpreendimentoService) ListUnidades(ctx context.Context, tenantID string, filter domain.UnidadeFilter) ([]dto.UnidadeResponse, error) {
	unidades, err := s.unidadeRepo.List(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}

	result := make([]dto.UnidadeResponse, len(unidades))
	for i, u := range unidades {
		result[i] = dto.ToUnidadeResponse(&u)
	}

	return result, nil
}

// GetFilters returns distinct values for all filterable fields.
func (s *EmpreendimentoService) GetFilters(ctx context.Context, tenantID string) (*domain.EmpreendimentoFilters, error) {
	return s.empreendimentoRepo.GetDistinctFilters(ctx, tenantID)
}

// ListConstrutoras returns construtoras for a tenant.
func (s *EmpreendimentoService) ListConstrutoras(ctx context.Context, tenantID string, filter domain.ConstrutoraFilter) ([]dto.ConstrutoraResponse, error) {
	construtoras, err := s.construtoraRepo.List(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ConstrutoraResponse, len(construtoras))
	for i, c := range construtoras {
		result[i] = dto.ToConstrutoraResponse(&c)
	}

	return result, nil
}
