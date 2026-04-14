package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/ponte-tech/cretor-back/modules/auth/domain"
	"github.com/ponte-tech/cretor-back/modules/auth/repository"
	"github.com/ponte-tech/cretor-back/shared/database"
)

func main() {
	_ = godotenv.Load()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := database.NewClient(ctx, uri)
	if err != nil {
		log.Fatal("failed to connect to mongodb:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("cretor")

	// Tenant ID fictício para teste (use um real do seu banco)
	tenantID := "000000000000000000000001"

	// Ensure indexes
	if err := repository.EnsureConstrutoraIndexes(ctx, db); err != nil {
		log.Fatal("indexes construtoras:", err)
	}
	if err := repository.EnsureEmpreendimentoIndexes(ctx, db); err != nil {
		log.Fatal("indexes empreendimentos:", err)
	}
	if err := repository.EnsureUnidadeIndexes(ctx, db); err != nil {
		log.Fatal("indexes unidades:", err)
	}
	fmt.Println("Indexes criados com sucesso")

	construtoraRepo := repository.NewConstrutoraRepository(db)
	empreendimentoRepo := repository.NewEmpreendimentoRepository(db)
	unidadeRepo := repository.NewUnidadeRepository(db)

	// 1. Criar Construtora: Alumbra
	construtora := &domain.Construtora{
		TenantID:   tenantID,
		Nome:       "Alumbra Empreendimentos",
		LogoS3Path: fmt.Sprintf("data/%s/construtoras/alumbra/logo.png", tenantID),
		Website:    "www.alumbraempreendimentos.com.br",
		Descricao:  "A única com propósito social",
		FonteID:    "d0c21526-919c-42e8-859e-c03b899d9839",
		Status:     "ativo",
	}

	if err := construtoraRepo.Create(ctx, construtora); err != nil {
		log.Fatal("create construtora:", err)
	}
	fmt.Printf("Construtora criada: %s (ID: %s)\n", construtora.Nome, construtora.ID)

	// 2. Criar Empreendimento: Diamond Hill Residence
	entregaDH := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	inicioDH := time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC)
	atualizacaoDH := time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC)

	diamondHill := &domain.Empreendimento{
		TenantID:      tenantID,
		ConstrutoraID: construtora.ID,
		Nome:          "Diamond Hill Residence",
		FonteID:       "cff237f8-ff65-48fc-b153-5dddb30c3cc2",
		Endereco: domain.Endereco{
			Completo:     "Capitão Ivo Manoel da Silva, 202, Perequê, Porto Belo - SC",
			Logradouro:   "Rua Capitão Ivo Manoel da Silva",
			Numero:       "202",
			Bairro:       "Perequê",
			Cidade:       "Porto Belo",
			UF:           "SC",
			Regiao:       "Litoral Norte",
			DistanciaMar: "150 metros",
		},
		Obra: domain.Obra{
			Status:        "em_construcao",
			Progresso:     75,
			Inicio:        &inicioDH,
			Entrega:       &entregaDH,
			AnoLancamento: 2022,
		},
		Incorporacao: "37.870",
		Edificacao: domain.Edificacao{
			Pavimentos:       26,
			UnidadesPorAndar: 4,
			TotalUnidades:    78,
			Coberturas:       2,
			AlturaTorre:      "90m",
		},
		Caracteristicas: domain.Caracteristicas{
			Dormitorios: domain.FaixaInt{Min: 2, Max: 4},
			Suites:      domain.FaixaInt{Min: 2, Max: 4},
			Vagas:       domain.FaixaInt{Min: 2, Max: 3},
			Metragem:    domain.FaixaFloat{Min: 90.00, Max: 216.37},
			Condominio:  500.00,
			IPTU:        1300.00,
		},
		DiferenciaisUnidade: []string{
			"Acabamento em gesso", "Cozinha Americana", "Porcelanato 70x70",
			"Portas Laqueadas", "Fechadura Biométrica", "Infra para Automação",
			"Sacada com Churrasqueira a Carvão", "Varanda Gourmet",
		},
		DiferenciaisCondominio: []string{
			"Elevador", "Playground", "Espaço Zen", "Espaço Chef Gourmet",
			"Poker Wine Bar", "Espaço Pet", "Espaço Fitness", "Prainha",
			"Spa Aquecido com Sauna Úmida", "Sistema de Segurança 24 hrs",
		},
		Secoes: []domain.Secao{
			{Nome: "Apartamentos - Tipo 1", Categoria: "apartamento", Ordem: 1},
			{Nome: "Apartamentos - Tipo 2", Categoria: "apartamento", Ordem: 2},
			{Nome: "Apartamentos - Tipo 3", Categoria: "apartamento", Ordem: 3},
			{Nome: "Apartamentos - Tipo 4", Categoria: "apartamento", Ordem: 4},
		},
		Galeria: []domain.GaleriaImagem{
			{S3Path: fmt.Sprintf("data/%s/empreendimentos/diamond-hill/galeria/001.jpg", tenantID), Ordem: 1, Tipo: "fachada"},
			{S3Path: fmt.Sprintf("data/%s/empreendimentos/diamond-hill/galeria/002.jpg", tenantID), Ordem: 2, Tipo: "interior"},
			{S3Path: fmt.Sprintf("data/%s/empreendimentos/diamond-hill/galeria/003.jpg", tenantID), Ordem: 3, Tipo: "area_lazer"},
		},
		CatalogoS3Path: fmt.Sprintf("data/%s/empreendimentos/diamond-hill/catalogo.pdf", tenantID),
		Comercializacao: domain.Comercializacao{
			UnidadesDisponiveis: 5,
			UltimaAtualizacao:   &atualizacaoDH,
		},
	}

	if err := empreendimentoRepo.Create(ctx, diamondHill); err != nil {
		log.Fatal("create empreendimento:", err)
	}
	fmt.Printf("Empreendimento criado: %s (ID: %s)\n", diamondHill.Nome, diamondHill.ID)

	// 3. Criar Unidades do Diamond Hill - Tipo 1 (amostra)
	v2052821 := 2052821.00
	v2114566 := 2114566.00
	v2176311 := 2176311.00

	unidades := []domain.Unidade{
		{TenantID: tenantID, EmpreendimentoID: diamondHill.ID, Secao: "Apartamentos - Tipo 1", Numero: "701", Andar: 7, Metragem: 123.80, ValorTexto: "Vendido", Status: "vendido"},
		{TenantID: tenantID, EmpreendimentoID: diamondHill.ID, Secao: "Apartamentos - Tipo 1", Numero: "801", Andar: 8, Metragem: 123.80, ValorTexto: "Vendido", Status: "vendido"},
		{TenantID: tenantID, EmpreendimentoID: diamondHill.ID, Secao: "Apartamentos - Tipo 1", Numero: "901", Andar: 9, Metragem: 123.80, ValorTexto: "Vendido", Status: "vendido"},
		{TenantID: tenantID, EmpreendimentoID: diamondHill.ID, Secao: "Apartamentos - Tipo 4", Numero: "804", Andar: 8, Metragem: 92.85, Valor: &v2052821, ValorTexto: "R$ 2.052.821,00", Status: "disponivel"},
		{TenantID: tenantID, EmpreendimentoID: diamondHill.ID, Secao: "Apartamentos - Tipo 4", Numero: "1604", Andar: 16, Metragem: 92.85, Valor: &v2114566, ValorTexto: "R$ 2.114.566,00", Status: "disponivel"},
		{TenantID: tenantID, EmpreendimentoID: diamondHill.ID, Secao: "Apartamentos - Tipo 4", Numero: "1704", Andar: 17, Metragem: 92.85, Valor: &v2176311, ValorTexto: "R$ 2.176.311,00", Status: "disponivel"},
		{TenantID: tenantID, EmpreendimentoID: diamondHill.ID, Secao: "Apartamentos - Tipo 2", Numero: "702", Andar: 7, Metragem: 90.00, Tipo: "Decorado", ValorTexto: "Reservado", Status: "reservado"},
		{TenantID: tenantID, EmpreendimentoID: diamondHill.ID, Secao: "Apartamentos - Tipo 4", Numero: "1404", Andar: 14, Metragem: 92.85, ValorTexto: "Reservado", Status: "reservado"},
	}

	if err := unidadeRepo.CreateMany(ctx, unidades); err != nil {
		log.Fatal("create unidades:", err)
	}
	fmt.Printf("Unidades criadas: %d\n", len(unidades))

	// 4. Verificar: listar tudo
	fmt.Println("\n--- Verificação ---")

	construtoras, _ := construtoraRepo.List(ctx, tenantID, domain.ConstrutoraFilter{})
	fmt.Printf("Construtoras: %d\n", len(construtoras))
	for _, c := range construtoras {
		fmt.Printf("  - %s (logo: %s)\n", c.Nome, c.LogoS3Path)
	}

	emps, _, _ := empreendimentoRepo.List(ctx, tenantID, domain.EmpreendimentoFilter{})
	fmt.Printf("Empreendimentos: %d\n", len(emps))
	for _, e := range emps {
		fmt.Printf("  - %s | %s | %s | Obra: %d%%\n", e.Nome, e.Endereco.Cidade, e.Endereco.UF, e.Obra.Progresso)
		fmt.Printf("    Seções: %d | Galeria: %d fotos | Catálogo: %s\n", len(e.Secoes), len(e.Galeria), e.CatalogoS3Path)
		fmt.Printf("    Diferenciais unidade: %d | Condomínio: %d\n", len(e.DiferenciaisUnidade), len(e.DiferenciaisCondominio))
	}

	unis, _ := unidadeRepo.List(ctx, tenantID, domain.UnidadeFilter{EmpreendimentoID: diamondHill.ID})
	fmt.Printf("Unidades: %d\n", len(unis))
	for _, u := range unis {
		fmt.Printf("  - %s | %s | %.2f m² | %s | %s\n", u.Secao, u.Numero, u.Metragem, u.ValorTexto, u.Status)
	}

	fmt.Println("\nSeed concluído com sucesso!")
}
