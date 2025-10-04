# Client-Server API - Cotação do Dólar

Sistema em Go que implementa um cliente e servidor para consulta da cotação USD-BRL.

## Funcionalidades

- **Server**: Consome API de cotação, salva no SQLite e expõe endpoint REST
- **Client**: Consulta o servidor e salva cotação em arquivo texto

## Como executar

1. **Instalar dependências**:
```bash
go mod tidy
```

2. **Executar o sistema**:
```bash
go run .
```

O programa iniciará o servidor e executará automaticamente o cliente de teste.

## Endpoints

- `GET /cotacao` - Retorna a cotação atual do dólar

## Arquivos gerados

- `quotations.db` - Banco SQLite com histórico de cotações
- `cotacao.txt` - Arquivo com a última cotação consultada

## Timeouts configurados

- API externa: 200ms
- Banco de dados: 10ms
- Cliente: 300ms

## Estrutura

- `main.go` - Ponto de entrada da aplicação
- `server.go` - Servidor HTTP e lógica de negócio
- `client.go` - Cliente HTTP e salvamento em arquivo