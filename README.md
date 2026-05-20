# app-openfinance-fake-no-database

Mock server do Open Finance Brasil rodando em pod, sem banco de dados. Toda a configuracao vive em um unico arquivo JSON, e novas rotas podem ser adicionadas em runtime via PUT no endpoint administrativo.

## Arquitetura

Hexagonal, com tres camadas:

```
cmd/server/                                  # entrypoint, wiring
internal/
  domain/                                    # entidades puras (Route, Config, Consentimento)
  port/
    input/                                   # interfaces consumidas pelos controllers
    output/                                  # interfaces de persistencia
  application/service/                       # MockService implementa as duas portas de entrada
  adapter/
    input/
      controller/                            # handlers Gin (mock e admin)
      routes/                                # registro dos endpoints fixos + NoRoute
      matcher/                               # match de path com parametros (":id")
    output/
      filestore/                             # JSONStore (escrita atomica via rename)
  config/logger/                             # zap singleton
```

O `MockService` e um struct unico que implementa tanto `input.Matcher` (consumida pelo handler generico) quanto `input.Admin` (consumida pelo handler administrativo). Mantem a Config em memoria com `sync.RWMutex`, persiste mutacoes no `output.ConfigStorage`.

## Como funciona o roteamento

Apenas `/healthz` e `/_admin/*` sao rotas explicitas no Gin. Tudo o mais cai no `gin.NoRoute`, que delega ao `MockController.Handle`. Ele:

1. Procura uma entrada em `routes` cuja chave seja `METODO PATH`, com matcher proprio para `:params`.
2. Se a rota tiver `skipConsent: false` (padrao), valida `x-consent-id` contra `consentimentos`.
3. Devolve `401` se o consentimento nao existe ou nao esta `AUTHORISED`.
4. Devolve `403` se a permissao da rota estiver na lista `negar` do consentimento.
5. Devolve o `body` configurado com o `status` definido.

Adicionar uma rota nova so requer um PUT, sem restart.

## Rodar localmente

```bash
go mod tidy
make run
```

Server sobe em `http://localhost:8080`.

## Endpoints administrativos

| Metodo | Path | Descricao |
|--------|------|-----------|
| GET    | `/healthz` | Health check |
| GET    | `/_admin/config` | Dump completo da configuracao |
| PUT    | `/_admin/routes` | Insere ou atualiza uma rota |
| DELETE | `/_admin/routes?key=GET%20/accounts` | Remove uma rota |
| PUT    | `/_admin/consentimentos` | Insere ou atualiza um consentimento |
| DELETE | `/_admin/consentimentos?consentId=urn:...` | Remove um consentimento |

### Exemplo: adicionar uma rota nova em runtime

```bash
curl -X PUT http://localhost:8080/_admin/routes \
  -H "Content-Type: application/json" \
  -d '{
    "key": "GET /credit-cards-accounts",
    "route": {
      "permission": "CREDIT_CARDS_ACCOUNTS_READ",
      "status": 200,
      "body": {
        "data": [{ "creditCardAccountId": "fake-cc-001" }],
        "meta": { "totalRecords": 1 }
      }
    }
  }'
```

A partir desse momento, `GET /credit-cards-accounts` ja responde, e o `mocks.json` em disco foi atualizado.

## Cenarios de teste pre-configurados

| Consent ID | Status | Permissoes negadas | Resultado |
|------------|--------|--------------------|-----------|
| `urn:obc-fake:consent-sucesso-001` | AUTHORISED | nenhuma | tudo 200 |
| `urn:obc-fake:consent-rejeitado-002` | REJECTED | n/a | 401 |
| `urn:obc-fake:consent-sem-balances-003` | AUTHORISED | ACCOUNTS_BALANCES_READ | balances retorna 403 |
| `urn:obc-fake:consent-sem-transactions-004` | AUTHORISED | ACCOUNTS_TRANSACTIONS_READ | transactions retorna 403 |
| `urn:obc-fake:consent-sem-overdraft-005` | AUTHORISED | ACCOUNTS_OVERDRAFT_LIMITS_READ | overdraft retorna 403 |

### Teste rapido

```bash
# token (skipConsent, retorna 200)
curl -X POST http://localhost:8080/auth/realms/fake/protocol/openid-connect/token

# accounts com sucesso
curl http://localhost:8080/accounts \
  -H "x-consent-id: urn:obc-fake:consent-sucesso-001"

# balances negado
curl -i http://localhost:8080/accounts/123/balances \
  -H "x-consent-id: urn:obc-fake:consent-sem-balances-003"

# consentimento rejeitado
curl -i http://localhost:8080/accounts \
  -H "x-consent-id: urn:obc-fake:consent-rejeitado-002"
```

## Deploy

```bash
make docker
kubectl apply -f deployments/k8s/pvc.yaml
kubectl apply -f deployments/k8s/deployment.yaml
kubectl apply -f deployments/k8s/service.yaml
```

O PVC garante que as alteracoes feitas via PUT sobrevivam a restart do pod. Se preferir verdade-no-Git, troque o PVC por um ConfigMap montado em readonly e perca as edicoes em runtime no proximo deploy.

## Variaveis de ambiente

| Variavel | Default | Descricao |
|----------|---------|-----------|
| `MOCKS_FILE` | `/data/mocks.json` | Caminho do arquivo operacional |
| `SEED_FILE` | `/seed/mocks.json` | Seed copiado para `MOCKS_FILE` na primeira execucao se este nao existir |
| `PORT` | `8080` | Porta de escuta |
