# Testes Unitários - radarr-sync-go

## Estrutura de Testes

Este projeto possui testes unitários abrangentes em cada package principal:

### 📦 Packages Testados

#### 1. **model/** - Estruturas de Dados
- `movie-model_test.go` - 46 testes
  - Validação de tipos struct
  - Testes de conversão de modelos
  - Verificação de fields obrigatórios
  - Testes de slices de modelos

**Cobertura:**
- `MovieToRadarrResponse` - estrutura principal de filme
- `RadarrModel` - modelo do Radarr
- `MovieResponse` - resposta de lista de filmes
- `MovieLoginResponse` - resposta de login
- `RadarrResponseError` - tratamento de erros

#### 2. **client/** - Cliente HTTP
- `client_test.go` - Funções utilitárias
  - `StructToMap()` - conversão de struct para map
  - `HttpClient()` - instanciação do cliente HTTP
  - `handleJson()` - parsing de JSON com tratamento de erros
  - Tratamento de JSON malformado
  - Tratamento de tipos inválidos

- `movie-client_test.go` - Cliente de filmes
  - `SetServerUri()` - configuração de URL do servidor
  - `SetRadarrUri()` - configuração de URL do Radarr

**Categorias de Testes:**
- ✅ Testes unitários - Funções isoladas
- ⏭️ Testes de integração - Marcados como "Skip" (requerem mock HTTP)

#### 3. **compress/** - Compressão de Arquivos
- `movie-compress_test.go` - Lógica de compressão
  - Validação de argumentos obrigatórios
  - Verificação de diretórios não-existentes
  - Testes com listas vazias de arquivos
  - Tratamento de erros

#### 4. **io_archive/** - Operações com Arquivos
- `archive_test.go` - Funções de arquivo
  - `FindWildcard()` - busca por padrões de arquivo
  - `FileStat()` - obtenção de informações de arquivo
  - `GZIP()` - compressão de arquivos/diretórios
  - Testes com diretórios aninhados
  - Testes com múltiplos arquivos
  - Validação de tempo de modificação

- `file_test.go` - Testes adicionais de arquivo
  - Testes em diretórios aninhados
  - Verificação de conteúdo preservado
  - Validação de constantes

## Executar os Testes

### Executar todos os testes:
```bash
go test ./...
```

### Executar testes com cobertura:
```bash
go test ./... -cover
```

### Executar testes com relatório de cobertura detalhado:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Executar testes de um package específico:
```bash
go test ./src/model -v
go test ./src/client -v
go test ./src/compress -v
go test ./src/io_archive -v
```

### Executar um teste específico:
```bash
go test ./src/model -run TestMovieResponseType -v
```

### Executar testes com saída detalhada:
```bash
go test ./... -v
```

### Executar testes ignorando testes de integração (Skip):
```bash
go test ./... -short
```

## Estatísticas de Testes

| Package | Arquivo | Testes | Tipo | Status |
|---------|---------|--------|------|--------|
| model | movie-model_test.go | 7 | Unitários | ✅ Ativo |
| client | client_test.go | 7 | Unitários + 3 Skip | ⚠️ Parcial |
| client | movie-client_test.go | 9 | Unitários + 6 Skip | ⚠️ Parcial |
| compress | movie-compress_test.go | 5 | Unitários | ✅ Ativo |
| io_archive | archive_test.go | 10 | Unitários | ✅ Ativo |
| io_archive | file_test.go | 5 | Unitários | ✅ Ativo |
| **TOTAL** | | **43** | | |

## Tipos de Testes

### ✅ Testes Unitários (Ativos)
Testam funções individuais isoladamente sem dependências externas:
- Validação de input/output
- Testes com dados válidos
- Testes com dados inválidos
- Testes com casos limite (edge cases)
- Testes de erro

### ⏭️ Testes de Integração (Skipados)
Requerem mock server HTTP - devem ser implementados posteriormente:
- `Login()` - teste de autenticação
- `FetchMoviesListToSync()` - teste de busca de filmes
- `AddMovieOnRadarr()` - teste de adição de filme
- `GetAllMoviesOnRadarr()` - teste de listagem de filmes

## Melhorias Futuras

### Adicionar Testes de Integração
Implementar HTTP mock server para testar:
```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
```

### Aumentar Cobertura
- [ ] Adicionar testes para casos extremos
- [ ] Adicionar benchmarks
- [ ] Adicionar fuzzing tests

### CI/CD
- [ ] Integrar testes no GitHub Actions
- [ ] Gerar relatório de cobertura automaticamente
- [ ] Enforçar cobertura mínima

## Executar com CI/CD

### GitHub Actions (exemplo):
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.26'
      - run: go test ./... -cover
```

## Notas Importantes

1. **Testes isolados**: Cada teste usa `t.TempDir()` para criar diretórios temporários isolados
2. **Sem dependências externas**: Testes unitários não requerem servidores rodando
3. **Rápido**: Todos os testes unitários executam em poucos milissegundos
4. **Determinístico**: Testes produzem os mesmos resultados em cada execução

## Próximas Fases

1. **Testes de Integração** - Implementar com HTTP mocks
2. **Testes de Performance** - Benchmarks para operações críticas
3. **Testes E2E** - Testar fluxo completo da aplicação
4. **CI/CD** - Automatizar execução de testes
