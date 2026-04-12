# Cretor - Estrutura MongoDB

Database: `cretor`

---

## Collections

1. `usuarios` - Usuarios do sistema (corretores, admins)
2. `clientes` - Clientes compradores de imoveis
3. `imoveis` - Imoveis disponiveis
4. `projetos` - Empreendimentos/lancamentos
5. `pipeline` - Negocios no funil de vendas
6. `match_scores` - Pontuacoes de compatibilidade cliente-imovel/projeto
7. `conversas` - Conversas WhatsApp com historico de mensagens
8. `leads` - Leads capturados (autoatendimento)

---

## 1. Collection: `usuarios`

Usuarios que operam o CRM (corretores, gerentes, admins).

```json
{
  "_id": "ObjectID",
  "tenant_id": "ObjectID",
  "nome": "string",
  "email": "string",
  "senha_hash": "string",
  "telefone": "string",
  "foto": "string | null",
  "role": "string",               // enum: admin, gerente, corretor
  "status": "string",             // enum: ativo, inativo
  "ultimo_login": "Date | null",
  "refresh_token": "string | null",
  "created_at": "Date",
  "updated_at": "Date"
}
```

**Indices:**
```
{ tenant_id: 1, email: 1 }           -> unique
{ tenant_id: 1, role: 1, status: 1 }
```

---

## 2. Collection: `clientes`

Perfil completo do cliente comprador com todas as preferencias para matching inteligente.

```json
{
  "_id": "ObjectID",
  "tenant_id": "ObjectID",

  // === DADOS PESSOAIS ===
  "nome": "string",
  "email": "string",
  "telefone": "string",
  "cpf": "string",
  "data_nascimento": "string",
  "sexo": "string",                    // enum: masculino, feminino, outro, nao_informado
  "foto": "string | null",

  // === DEMOGRAFICOS ===
  "estado_civil": "string",            // enum: solteiro, casado, divorciado, viuvo, uniao_estavel
  "profissao": "string",
  "renda_mensal": "Decimal128",
  "patrimonio": "Decimal128",

  // === COMPOSICAO FAMILIAR ===
  "familia": {
    "tem_filhos": "boolean",
    "numero_filhos": "int",
    "idade_filhos": ["int"],
    "tem_pets": "boolean",
    "tipo_pets": ["string"]            // gato, cachorro_pequeno, cachorro_grande, outros
  },

  // === PREFERENCIAS DE IMOVEL ===
  "preferencias": {
    "tipo_imovel": ["string"],         // apartamento, casa, cobertura, terreno, comercial
    "finalidade": "string",            // morar, investir, alugar, temporada
    "metragem_min": "int",
    "metragem_max": "int",
    "quartos": "int",
    "banheiros": "int",
    "vagas": "int",
    "preferencia_andar": "string",     // baixo, medio, alto, cobertura, indiferente
    "vista_desejada": "string",        // mar, cidade, montanha, parque, indiferente
    "imovel_novo": "string",           // novo, usado, planta, indiferente
    "estado_conservacao": "string",    // excelente, bom, reforma_leve, reforma_total
    "tipo_acabamento": "string",       // alto_padrao, medio, basico, indiferente
    "orientacao_solar": "string",      // norte, sul, leste, oeste, indiferente
    "area_externa_privativa": "boolean",
    "condominio_fechado": "boolean",

    // Caracteristicas desejadas
    "caracteristicas": ["string"],     // piscina, churrasqueira, varanda_gourmet, home_office,
                                       // lareira, closet, banheira, vista_livre, pe_direito_alto,
                                       // ar_condicionado, aquecimento, lavabo, dependencia_empregada

    // Infraestrutura do condominio
    "areas_lazer": ["string"],         // academia, salao_festas, playground, quadra,
                                       // sauna, spa, coworking, brinquedoteca, piscina
    "tamanho_condominio": "string",    // pequeno, medio, grande, indiferente
    "importancia_area_lazer": "int",   // 1-5
    "perfil_vizinhanca": "string"      // familias, jovens, idosos, misto
  },

  // === ACESSIBILIDADE ===
  "acessibilidade": {
    "necessita": "boolean",
    "tipos": ["string"],               // rampa, elevador, banheiro_adaptado, porta_larga
    "mobilidade_reduzida": "boolean",
    "tem_bebe": "boolean",
    "necessita_dependencia": "boolean"
  },

  // === LOCALIZACAO ===
  "localizacao": {
    "bairros_preferidos": ["string"],
    "cidades_preferidas": ["string"],
    "zona_preferida": "string",        // norte, sul, leste, oeste, centro
    "endereco_trabalho": "string",
    "tempo_deslocamento_max": "int",   // minutos
    "modal_transporte": "string",      // carro, transporte_publico, moto, bicicleta
    "ja_mora_cidade": "boolean",
    "cidade_atual": "string",

    // Pontos de interesse proximos
    "pontos_interesse": {
      "escolas": "boolean",
      "universidades": "boolean",
      "hospitais": "boolean",
      "shoppings": "boolean",
      "parques": "boolean",
      "academias": "boolean",
      "restaurantes": "boolean",
      "supermercados": "boolean",
      "transporte_publico": "boolean",
      "praia": "boolean",
      "metro": "boolean"
    }
  },

  // === FINANCEIRO ===
  "financeiro": {
    "orcamento_min": "Decimal128",
    "orcamento_max": "Decimal128",
    "formas_pagamento": ["string"],    // vista, financiamento, permuta, consorcio
    "possui_fgts": "boolean",
    "valor_fgts": "Decimal128",
    "imovel_entrada": "boolean",       // tem imovel para dar como entrada
    "valor_imovel_entrada": "Decimal128",
    "score_credito": "string",         // excelente, bom, regular, ruim
    "preferencia_parcelas": "int",
    "pode_fiador": "boolean",
    "renda_comprometida": "int"        // percentual 0-100
  },

  // === PRIORIDADES E DEAL BREAKERS ===
  "prioridades": {
    "localizacao": "int",              // peso 1-5
    "tamanho": "int",
    "preco": "int",
    "caracteristicas": "int",
    "condominio": "int"
  },
  "deal_breakers": ["string"],         // via_expressa, aeroporto, sem_elevador,
                                       // andar_baixo, sem_vaga, barulho, sem_portaria
  "must_haves": ["string"],            // top 3 itens obrigatorios

  // === SUSTENTABILIDADE E TECNOLOGIA ===
  "sustentabilidade": {
    "carro_eletrico": "boolean",
    "certificacao_sustentavel": "boolean",
    "automacao_residencial": "boolean",
    "energia_solar": "boolean"
  },

  // === ESTILO DE VIDA ===
  "estilo_vida": {
    "trabalha_home": "boolean",
    "pratica_esportes": "boolean",
    "gosta_cozinhar": "boolean",
    "recebe_muito": "boolean",         // recebe visitas frequentemente
    "viaja_frequente": "boolean",
    "necessita_escritorio": "boolean",
    "horario_trabalho": "string"       // comercial, flexivel, noturno
  },

  // === PROCESSO DE COMPRA ===
  "processo_compra": {
    "urgencia": "string",              // alta, media, baixa
    "motivo_compra": "string",
    "ja_visitou_imoveis": "boolean",
    "tem_imovel_venda": "boolean",
    "prazo_mudanca": "string",         // imediato, 3_meses, 6_meses, 1_ano, indefinido
    "proposta_feita": "boolean",
    "motivo_busca": "string"
  },

  // === COMPORTAMENTO DE COMPRA ===
  "comportamento": {
    "imoveis_visitados": "int",
    "tempo_procurando": "string",      // menos_1_mes, 1_3_meses, 3_6_meses, mais_6_meses
    "principal_problema": "string",
    "nivel_pesquisa": "string",        // iniciante, intermediario, avancado
    "envolve_decisao": ["string"]      // conjuge, filhos, pais, socio, sozinho
  },

  // === METADATA ===
  "observacoes": "string",
  "status": "string",                  // enum: ativo, inativo, prospecto, cliente
  "origem": "string",                  // enum: site, indicacao, redes_sociais, evento, telefone, outro
  "responsavel": "string",             // nome do corretor
  "created_at": "Date",
  "updated_at": "Date"
}
```

**Indices:**
```
{ tenant_id: 1, email: 1 }                               -> unique
{ tenant_id: 1, cpf: 1 }                                 -> unique, sparse
{ tenant_id: 1, status: 1, created_at: -1 }
{ tenant_id: 1, "localizacao.cidades_preferidas": 1 }
{ tenant_id: 1, nome: "text", email: "text" }             -> text search
```

**Design decisions:**
- Preferencias, acessibilidade, localizacao, financeiro, estilo_vida sao **embedded** (1:1, sempre lidos juntos, sem ciclo de vida proprio)
- Arrays como `caracteristicas`, `deal_breakers`, `must_haves` sao bounded (max ~20 itens) - seguro embeddar
- `responsavel` usa Extended Reference Pattern (copia do nome do corretor, evita $lookup)

---

## 3. Collection: `imoveis`

Imoveis disponiveis para venda/locacao.

```json
{
  "_id": "ObjectID",
  "tenant_id": "ObjectID",

  // === INFORMACOES BASICAS ===
  "tipo": "string",                    // enum: apartamento, casa, cobertura, terreno, comercial
  "titulo": "string",
  "descricao": "string",

  // === LOCALIZACAO ===
  "endereco": "string",
  "bairro": "string",
  "cidade": "string",
  "estado": "string",
  "cep": "string",
  "coordenadas": {                     // GeoJSON para queries espaciais
    "type": "Point",
    "coordinates": ["double", "double"] // [longitude, latitude]
  },

  // === CARACTERISTICAS ===
  "quartos": "int",
  "banheiros": "int",
  "vagas": "int",
  "area_total": "double",             // m2
  "area_util": "double",              // m2
  "andar": "int",
  "ano_construcao": "int",
  "mobiliado": "boolean",
  "aceita_pets": "boolean",

  // === FINANCEIRO ===
  "preco": "Decimal128",
  "condominio": "Decimal128",          // taxa mensal
  "iptu": "Decimal128",               // valor anual

  // === FEATURES ===
  "caracteristicas": ["string"],       // varanda_gourmet, home_office, lareira, closet,
                                       // vista_livre, pe_direito_alto, piscina_privativa, etc.

  // === STATUS ===
  "status": "string",                  // enum: disponivel, reservado, vendido, em_construcao
  "disponibilidade": "string",         // imediata ou data futura

  // === MIDIA ===
  "fotos": ["string"],                 // URLs
  "video_url": "string | null",
  "tour_virtual_url": "string | null",

  // === METADATA ===
  "corretor_responsavel": "string",
  "created_at": "Date",
  "updated_at": "Date"
}
```

**Indices:**
```
{ tenant_id: 1, status: 1, preco: 1 }
{ tenant_id: 1, tipo: 1, cidade: 1, bairro: 1 }
{ tenant_id: 1, quartos: 1, area_util: 1 }
{ tenant_id: 1, titulo: "text", descricao: "text", bairro: "text" }
{ "coordenadas": "2dsphere" }                              -> geoespacial
```

**Design decisions:**
- `coordenadas` usa formato GeoJSON para permitir queries `$near` e `$geoWithin`
- `fotos` e array bounded (max ~30 fotos) - seguro embeddar
- `preco`, `condominio`, `iptu` sao `Decimal128` para precisao monetaria

---

## 4. Collection: `projetos`

Empreendimentos imobiliarios (lancamentos, construcoes).

```json
{
  "_id": "ObjectID",
  "tenant_id": "ObjectID",

  // === INFORMACOES BASICAS ===
  "nome": "string",
  "construtora": "string",
  "status": "string",                  // enum: lancamento, em_construcao, pronto, entregue
  "tipo_empreendimento": "string",     // enum: residencial, comercial, misto
  "descricao": "string",

  // === LOCALIZACAO ===
  "endereco": "string",
  "bairro": "string",
  "cidade": "string",
  "estado": "string",
  "cep": "string",
  "coordenadas": {
    "type": "Point",
    "coordinates": ["double", "double"]
  },

  // === UNIDADES ===
  "total_unidades": "int",
  "unidades_disponiveis": "int",
  "tipos_unidades": ["string"],        // 1_dorm, 2_dorm, 3_dorm, cobertura, sala_comercial
  "area_privativa_min": "double",
  "area_privativa_max": "double",
  "vagas_min": "int",
  "vagas_max": "int",

  // === FINANCEIRO ===
  "preco_min": "Decimal128",
  "preco_max": "Decimal128",
  "entrada_minima": "Decimal128",
  "aceita_financiamento": "boolean",

  // === TIMELINE ===
  "data_lancamento": "Date",
  "previsao_entrega": "Date",
  "fase_obra": "string",              // enum: fundacao, estrutura, acabamento, pronto
  "percentual_concluido": "int",       // 0-100

  // === INFRAESTRUTURA ===
  "areas_lazer": ["string"],           // piscina, academia, salao_festas, churrasqueira,
                                       // playground, quadra, sauna, spa, coworking
  "seguranca": ["string"],             // portaria_24h, cameras, cerca_eletrica, controle_acesso
  "sustentabilidade": ["string"],      // energia_solar, coleta_seletiva, reuso_agua, areas_verdes

  // === MARKETING ===
  "diferenciais": ["string"],
  "logo": "string | null",
  "foto_destaque": "string | null",
  "fotos": ["string"],
  "plantas": ["string"],               // URLs das plantas

  // === GESTAO ===
  "corretor_responsavel": "string",
  "vendedores": ["string"],

  // === METADATA ===
  "created_at": "Date",
  "updated_at": "Date"
}
```

**Indices:**
```
{ tenant_id: 1, status: 1, tipo_empreendimento: 1 }
{ tenant_id: 1, cidade: 1, preco_min: 1 }
{ tenant_id: 1, nome: "text", construtora: "text", descricao: "text" }
{ "coordenadas": "2dsphere" }
```

---

## 5. Collection: `pipeline`

Negocios no funil de vendas. Usa **Extended Reference Pattern** para evitar $lookup em cliente e imovel.

```json
{
  "_id": "ObjectID",
  "tenant_id": "ObjectID",

  // === REFERENCIAS (Extended Reference Pattern) ===
  "cliente_id": "ObjectID",
  "cliente_nome": "string",            // copia - evita $lookup
  "cliente_foto": "string | null",     // copia
  "cliente_email": "string",           // copia
  "cliente_telefone": "string",        // copia

  "imovel_id": "ObjectID | null",      // pode ser projeto tambem
  "imovel_titulo": "string",           // copia
  "imovel_foto": "string | null",      // copia
  "imovel_endereco": "string",         // copia

  "projeto_id": "ObjectID | null",     // quando o negocio e sobre um projeto

  // === STATUS DO NEGOCIO ===
  "etapa": "string",                   // enum: primeiro_contato, qualificado, visita_agendada,
                                       //       proposta_enviada, negociacao, fechado, perdido
  "prioridade": "string",             // enum: baixa, media, alta, urgente
  "valor_negocio": "Decimal128",
  "probabilidade_fechamento": "int",   // 0-100

  // === TIMELINE ===
  "data_criacao": "Date",
  "data_ultima_interacao": "Date",
  "data_movimentacao": "Date",         // ultima troca de etapa
  "dias_na_etapa": "int",

  // === PROXIMOS PASSOS ===
  "proxima_acao": "string",
  "data_proxima_acao": "Date | null",

  // === TRACKING ===
  "ultima_anotacao": "string",
  "corretor_responsavel": "string",
  "tags": ["string"],
  "motivo_perda": "string | null",     // preenchido quando etapa == perdido

  // === METADATA ===
  "created_at": "Date",
  "updated_at": "Date"
}
```

**Indices:**
```
{ tenant_id: 1, etapa: 1, data_criacao: -1 }
{ tenant_id: 1, cliente_id: 1 }
{ tenant_id: 1, imovel_id: 1 }
{ tenant_id: 1, corretor_responsavel: 1, etapa: 1 }
{ tenant_id: 1, prioridade: 1, data_proxima_acao: 1 }
```

**Design decisions:**
- Extended Reference Pattern: `cliente_nome`, `cliente_foto`, `imovel_titulo` etc. sao copias desnormalizadas
- Atualizar copias quando o original mudar (evento ou batch job)
- `dias_na_etapa` calculado e atualizado quando muda de etapa

---

## 6. Collection: `match_scores`

Pontuacoes de compatibilidade entre clientes e imoveis/projetos.

```json
{
  "_id": "ObjectID",
  "tenant_id": "ObjectID",

  // === REFERENCIAS ===
  "cliente_id": "ObjectID",
  "cliente_nome": "string",            // Extended Reference

  // Um ou outro preenchido
  "imovel_id": "ObjectID | null",
  "imovel_titulo": "string | null",
  "projeto_id": "ObjectID | null",
  "projeto_nome": "string | null",

  // === SCORE ===
  "score": "int",                      // 0-100
  "probabilidade_fechamento": "string", // enum: baixa, media, alta, muito_alta

  // === DETALHAMENTO DO SCORE ===
  "criterios": [
    {
      "criterio": "string",            // localizacao, preco, tamanho, quartos, must_haves, etc.
      "match": "boolean",
      "peso": "int",                   // pontos atribuidos a este criterio
      "peso_max": "int",               // pontuacao maxima possivel
      "detalhe": "string"              // descricao do match/mismatch
    }
  ],

  // === DEAL BREAKERS DETECTADOS ===
  "deal_breakers_ativos": ["string"],   // lista de deal breakers que impedem o match
  "eliminado": "boolean",              // true se algum deal breaker ativo

  // === METADATA ===
  "calculado_em": "Date",
  "updated_at": "Date"
}
```

**Indices:**
```
{ tenant_id: 1, cliente_id: 1, score: -1 }
{ tenant_id: 1, imovel_id: 1, score: -1 }
{ tenant_id: 1, projeto_id: 1, score: -1 }
{ tenant_id: 1, eliminado: 1, score: -1 }
```

**Design decisions:**
- `criterios` e array bounded (max ~15 criterios) - seguro embeddar
- Recalculado quando cliente ou imovel/projeto e atualizado
- `eliminado` permite filtrar rapidamente matches invalidos

---

## 7. Collection: `conversas`

Conversas WhatsApp. Mensagens sao **embedded** (Subset Pattern - manter ultimas N).

```json
{
  "_id": "ObjectID",
  "tenant_id": "ObjectID",

  // === CONTATO ===
  "contato_nome": "string",
  "contato_telefone": "string",
  "contato_avatar": "string | null",
  "cliente_id": "ObjectID | null",     // vinculo com cliente (se existir)

  // === STATUS ===
  "online": "boolean",
  "nao_lidas": "int",
  "ultima_mensagem": "string",
  "horario_ultima_mensagem": "Date",

  // === MENSAGENS (Subset - ultimas 50) ===
  "mensagens": [
    {
      "_id": "ObjectID",
      "tipo": "string",                // enum: texto, audio, imagem, documento
      "conteudo": "string",
      "remetente": "string",           // enum: eu, contato
      "horario": "Date",
      "lida": "boolean",

      // Audio
      "transcricao_audio": "string | null",
      "duracao_audio": "string | null",

      // Traducoes
      "traducao": {
        "en": "string | null",
        "es": "string | null"
      }
    }
  ],

  // === METADATA ===
  "created_at": "Date",
  "updated_at": "Date"
}
```

**Collection auxiliar: `mensagens_historico`** (historico completo)
```json
{
  "_id": "ObjectID",
  "conversa_id": "ObjectID",
  "tenant_id": "ObjectID",
  "tipo": "string",
  "conteudo": "string",
  "remetente": "string",
  "horario": "Date",
  "lida": "boolean",
  "transcricao_audio": "string | null",
  "duracao_audio": "string | null",
  "traducao": {
    "en": "string | null",
    "es": "string | null"
  }
}
```

**Indices - conversas:**
```
{ tenant_id: 1, horario_ultima_mensagem: -1 }
{ tenant_id: 1, cliente_id: 1 }
{ tenant_id: 1, contato_telefone: 1 }
```

**Indices - mensagens_historico:**
```
{ conversa_id: 1, horario: -1 }
{ tenant_id: 1, horario: 1 }           -> TTL opcional para limpeza
```

**Design decisions:**
- **Subset Pattern**: `conversas.mensagens` mantem apenas as ultimas 50 mensagens
- Mensagens antigas vao para `mensagens_historico`
- Nova mensagem: `$push` + `$slice: -50` na conversa, e `InsertOne` no historico
- Evita documentos gigantes em conversas longas

---

## 8. Collection: `leads`

Leads capturados via autoatendimento e formularios.

```json
{
  "_id": "ObjectID",
  "tenant_id": "ObjectID",

  // === DADOS DO LEAD ===
  "nome": "string",
  "email": "string",
  "telefone": "string",
  "mensagem": "string | null",

  // === ORIGEM ===
  "origem": "string",                  // enum: site, landing_page, indicacao, redes_sociais, evento
  "token": "string",                   // token do link de autoatendimento
  "url_origem": "string | null",

  // === STATUS ===
  "status": "string",                  // enum: novo, contatado, qualificado, convertido, descartado
  "convertido_cliente_id": "ObjectID | null",  // quando vira cliente

  // === PREFERENCIAS BASICAS (preenchidas no autoatendimento) ===
  "preferencias_resumo": {
    "tipo_imovel": "string | null",
    "cidade": "string | null",
    "orcamento_max": "Decimal128 | null",
    "quartos": "int | null"
  },

  // === ATRIBUICAO ===
  "corretor_responsavel": "string | null",

  // === METADATA ===
  "created_at": "Date",
  "updated_at": "Date"
}
```

**Indices:**
```
{ tenant_id: 1, status: 1, created_at: -1 }
{ tenant_id: 1, email: 1 }
{ token: 1 }                           -> unique, sparse
```

---

## Resumo de Relacionamentos

```
usuarios (1) --[gerencia]--> (N) clientes
usuarios (1) --[gerencia]--> (N) imoveis
usuarios (1) --[gerencia]--> (N) projetos

clientes (1) --[participa]--> (N) pipeline
imoveis  (1) --[participa]--> (N) pipeline
projetos (1) --[participa]--> (N) pipeline

clientes (1) --[tem score]--> (N) match_scores
imoveis  (1) --[tem score]--> (N) match_scores
projetos (1) --[tem score]--> (N) match_scores

clientes (1) --[conversa]--> (N) conversas
conversas (1) --[contem]--> (N) mensagens (embedded, subset)
conversas (1) --[historico]--> (N) mensagens_historico

leads (1) --[converte]--> (0..1) clientes
```

---

## Regras de Multi-Tenancy

1. **Todo documento tem `tenant_id`** como primeiro campo apos `_id`
2. **Todo indice comeca com `tenant_id`** para isolamento de dados
3. **Toda query filtra por `tenant_id`** - enforced na camada de middleware/repository
4. Indices unique sao compostos: `{ tenant_id: 1, email: 1 }` (email unico POR tenant)

---

## Regras de Denormalizacao (Extended Reference Pattern)

Campos copiados que precisam ser atualizados quando o original muda:

| Collection | Campo copiado | Origem |
|---|---|---|
| `pipeline` | `cliente_nome`, `cliente_foto`, `cliente_email`, `cliente_telefone` | `clientes` |
| `pipeline` | `imovel_titulo`, `imovel_foto`, `imovel_endereco` | `imoveis` |
| `match_scores` | `cliente_nome` | `clientes` |
| `match_scores` | `imovel_titulo`, `projeto_nome` | `imoveis`, `projetos` |

**Estrategia de atualizacao:** Quando atualizar um cliente/imovel/projeto, executar `UpdateMany` nas collections que referenciam para manter copias sincronizadas.

---

## Tamanho Estimado (M0 Free Tier - 512MB)

| Collection | Docs estimados | Tamanho medio/doc | Total estimado |
|---|---|---|---|
| usuarios | 50 | 1 KB | 50 KB |
| clientes | 5.000 | 5 KB | 25 MB |
| imoveis | 2.000 | 3 KB | 6 MB |
| projetos | 200 | 4 KB | 800 KB |
| pipeline | 10.000 | 2 KB | 20 MB |
| match_scores | 50.000 | 1.5 KB | 75 MB |
| conversas | 5.000 | 15 KB | 75 MB |
| mensagens_historico | 200.000 | 0.5 KB | 100 MB |
| leads | 10.000 | 1 KB | 10 MB |
| **Total** | | | **~312 MB** |

Cabe confortavelmente no M0 Free Tier (512 MB) para fase inicial.
