# Plano de Execução: Fase 2 - Rizoma Engine P2P (Rust)

Este plano detalha as especificações técnicas para o desenvolvimento da **Rizoma Engine em Rust** (`projetos/rizoma`), o motor P2P de rede subterrânea, criptografia de chaves soberanas Ed25519, sincronização local-first (CRDT) e o daemon `rizomad` para transformar a máquina local na primeira **Rizoma Box Alpha**.

---

## 🎯 Arquitetura da Rizoma Engine

1. **Linguagem & Performance:** A engine é construída em **Rust 2024** para obter o máximo de desempenho, segurança de memória sem garbage collector e facilidade de compilação nativa para celulares (Android/iOS) e hardware livre (Raspberry Pi).
2. **Daemon `rizomad`:** Roda como um serviço em background escutando chamadas IPC/HTTP locais (porta `9090`) do `nucleo-engine` e do `App Urupê`.
3. **Criptografia Soberana:** Implementa pares de chaves criptográficas **Ed25519** para assinatura soberana de mensagens e artigos sem dependência de autoridade central.

---

## 📁 Estrutura do Repositório `projetos/rizoma`

```
projetos/rizoma/
├── Cargo.toml               # Dependências: tokio, ed25519-dalek, serde, axum, tracing
└── src/
    ├── main.rs              # Ponto de entrada do daemon rizomad (Nó Alpha)
    ├── identity.rs          # Gestão de chaves públicas/privadas Ed25519
    ├── crdt.rs              # Engine de sincronização Local-First (Automerge CRDTs)
    ├── mesh.rs              # Transporte P2P e transmissão de pacotes assinados
    └── api.rs               # API REST/IPC local na porta 9090 para o App e Núcleo
```

---

## 🔌 Conector no Backend Go (`projetos/nucleo-urupe/nucleo-engine`)

- **`internal/data/p2p/rizoma_client.go`**: Cliente HTTP em Go para o `nucleo-engine` comunicar-se com o daemon local `rizomad` (localhost:9090) e sincronizar artigos do CMS automaticamente na malha P2P.

---

## 🧪 Plano de Verificação

### Testes Automatizados
- Compilar a Rizoma Engine em Rust: `cargo build` no diretório `projetos/rizoma`.
- Executar os testes unitários do módulo de identidade e criptografia em Rust: `cargo test` no diretório `projetos/rizoma`.
- Compilar o backend Go do Núcleo Urupê: `CGO_ENABLED=0 go build -o /tmp/nucleo_test ./cmd/bot` em `projetos/nucleo-urupe/nucleo-engine`.

### Verificação Manual
- Iniciar o daemon `rizomad` na porta 9090 e verificar a resposta de status com chave de nó Ed25519.
- Publicar um artigo no CMS do Núcleo Urupê e confirmar o envio assinado para o daemon da Rizoma Box local.
