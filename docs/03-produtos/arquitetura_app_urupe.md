# 📱 Arquitetura Técnica Mestre: App Urupê (Super-App Descentralizado)

> **"Engenharia soberana local-first: Rust Core nativo, renderização GPU a 120 FPS, rede P2P sem blockchains corporativas e isolamento por WebAssembly."**

---

## 🏛️ 1. Diagnóstico & Consenso Arquitetural

Após avaliação por primeiros princípios, validamos a crítica ao uso do Tauri v2 no ecossistema **Mobile**:

```mermaid
graph TD
    subgraph CAMADA CORE RUST (O CÉREBRO INDEPENDENTE)
        RC[🧠 Rust Core: Criptografia MLS/Signal + SQLite + CRDT Loro/Automerge + Network P2P iroh/libp2p]
    end

    subgraph ADAPTADORES DE INTERFACE (UI LAYER)
        RC -->|FFI Direta C-ABI / Zero Copy| MOB[📱 Mobile App: Flutter Impeller / Native GPU 120 FPS]
        RC -->|Compilação WebAssembly / Wasm| WEB[🌐 Web & Portal Público: Svelte 5 / DOM Ultra-Leve]
    end

    subgraph ECOSSISTEMA DE MINI-APPS
        MOB -->|Sandbox Wasmtime/Wasmer + WASI| MINI[📦 Mini-Apps Comunitários Assinados]
    end
```

---

## 📐 2. Os 6 Pilares do Consenso Técnico

### 1. 🧠 Core Desacoplado em Rust (O Cérebro)
- **Onde roda:** No mesmo processo nativo no Mobile (via C-ABI FFI) e via WebAssembly no browser.
- **Responsabilidades:**
  - Criptografia ponta a ponta (**MLS / Signal Protocol**).
  - Banco de dados local-first soberano em **SQLite + FTS5**.
  - Sincronização incremental de dados em tempo real via **CRDTs (Loro / Automerge)**.
  - Daemon de rede **P2P (`libp2p` / `iroh` / `Bao` BLAKE3)**.

### 2. 📱 Mobile Front-End (Flutter Impeller / Kotlin Compose)
- **Por que não Tauri no Mobile:** No Android/iOS, o isolamento rigoroso da WebView do browser congela sockets P2P quando o aplicativo vai para segundo plano. Além disso, o DOM da WebView sofre com *layout thrashing* ao rolar feeds pesados de mídias/vídeos.
- **A Solução:** Adotar **Flutter com engine Impeller** no Mobile. A comunicação com o *Rust Core* ocorre via FFI direta no mesmo processo em nanossegundos (sem overhead de serialização JSON). A UI renderiza direto na GPU a 120 FPS em celulares populares.

### 3. 🌐 Web & Portal Público (Svelte 5 + Wasm)
- **Por que Svelte na Web:** O Flutter Web possui um *bundle* inicial pesado (CanvasKit). No browser, onde o carregamento rápido de notícias e o SEO do Núcleo Urupê são vitais, o **Svelte 5** domina absoluto: páginas carregadas em milissegundos com consumo mínimo de dados.
- O Rust Core roda compilado para **WebAssembly (`wasm-bindgen`)**, mantendo a mesma lógica de negócios do mobile no navegador.

### 4. 📦 Sandbox de Mini-Apps Soberanos (WASI)
- Permite que a comunidade desenvolva mini-aplicativos (ex: mapas de feiras agroecológicas, logística de entregas solidárias, rádios comunitárias) rodando em uma sandbox ultra-segura baseada em **WebAssembly (Wasmtime / Wasmer)** com permissões estritas estilo WASI e assinaturas criptográficas comunitárias.

### 5. 🌾 Protocolo P2P Sem Blockchain Predatória
- Zero *gas fees*, zero especulação de criptomoedas, zero desperdício de prova de trabalho.
- Transporte via **QUIC / WebTransport** e redes mesh locais (**Bluetooth LE / Wi-Fi Direct**) para funcionamento em atos públicos e territórios sem internet.
- Verificação de arquivos sob demanda por árvores **BLAKE3 (Iroh/Bao)**.
