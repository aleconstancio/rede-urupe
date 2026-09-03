# 🌿 Frente Urupê: Ecossistema Monorepo

> **"Pela Soberania Digital, Emancipação Popular e Organização Ecossocialista."**

O repositório **`rede-urupe`** é o monorepo unificado da **Frente Urupê**. Ele engloba a documentação teórica do movimento (PKM) e o código-fonte da nossa **Trindade de Produtos**: o **Núcleo Urupê** (QG & Plataforma de Controle), o **Rizoma** (Protocolo P2P Local-First & Hardware Open-Source) e o **App Urupê** (Super-App Popular Soberano).

---

## 🏛️ A Trindade de Produtos

```mermaid
graph TD
    A[Monorepo rede-urupe] --> B[1. Núcleo Urupê: QG & Plataforma de Controle]
    A --> C[2. Rizoma: Infraestrutura P2P & Hardware Open-Source]
    A --> D[3. App Urupê: Super-App Popular Soberano]

    B --> B1[Go 1.26 + Svelte 5 + Mascote IA Micélia 🍄 + Estúdio 1.2 + Gruvbox Solarpunk]
    C --> C1[Rust Core + iroh/libp2p + Roteamento Cebola + Mesh Offline]
    D --> D1[Svelte 5 + Tauri v2: Mensageria, Feed Plural, Uber/iFood Solidários & Urupê Mídia]
```

### 1. 🏛️ Núcleo Urupê (`projetos/nucleo-urupe`) — Version 1.2
Central de inteligência, governança e controle da Frente Urupê.
- **Site Público (3 Abas):** `Quem Somos` (História & Manifesto Ativo), `Notícias` (Urupê News no estilo VicoNews) e `Projetos` (Vitrine do ecossistema).
- **Estúdio Micélium 🍄 2.0:** Versionador de edições do Manifesto (`v1.0`, `v1.1`, `v1.2`), Playground dos 7 estágios da **Arandu Engine** 🏹 e inspetor de memórias FTS5.
- **Spore Ops 🍄 (Guerrilha & Agitprop):** Central de clipping de imprensa e disparos de agitação popular.
- **Guará Geo 🪶 (Inteligência Geoespacial):** Sensoriamento por satélites (*Sentinel-2A, Landsat 9*), alertas climáticos e apoio a brigadas agroecológicas.
- **Jataí Ops 🐝 (Frota B2B & Autofinanciamento):** Dashboard da frota de agentes industriais (Vico, RH, etc.) e transparência financeira do **Fundo Urupê**.

### 2. 🌿 Rizoma (`projetos/rizoma`)
Infraestrutura descentralizada e anti-predatória.
- **Web3 Real (Sem Criptomoedas/NFTs):** Dados local-first, encriptados Ed25519 e soberanos.
- **Resiliência a Censura:** Operação online com capacidade de malha P2P offline (Wi-Fi Direct/Bluetooth) para atos públicos.
- **Hardware Open-Source (Rizoma Box):** Firmware para micro-computadores comunitários.

### 3. 📱 App Urupê (`projetos/app-urupe`)
Super-App de uso diário da classe trabalhadora.
- **Mensageria Soberana:** Chat P2P instantâneo com assistente IA Micélia e cofre de pânico (ZKP).
- **Feed Plural:** Cultura, humor, memes, luta de classes, formação política e vlogs.
- **Mapa Territorial & Logística:** Mobilidade soberana ("Uber Popular") e entregas solidárias sem taxas extrativistas.

---

## 🎨 Identidade Visual: Gruvbox Solarpunk
A interface do Núcleo Urupê adota a estética **Gruvbox Solarpunk**:
- **Palette Terrosa:** Verde Clorofila (`#b8bb26`), Ocre Solar (`#fabd2f`), Terracota Argila (`#fe8019`), Carvão Terroso (`#1d2021`) e Pergaminho Creme (`#fbf1c7`).
- **Tipografia:** Títulos em *Cormorant Garamond*, texto em *Source Sans 3* e código/métricas em *IBM Plex Mono*.

---

## 🧠 Base de Conhecimento & PKM (`docs/`)
Toda a documentação teórica e arquitetural está organizada no diretório `docs/`.
Consulte o mapa central de navegação em [docs/README.md](file:///home/ale/Projects/rede-urupe/docs/README.md):

1. **[01-manifesto](file:///home/ale/Projects/rede-urupe/docs/01-manifesto/manifesto.md):** O Manifesto Ecossocialista por uma Soberania Digital Popular.
2. **[02-arquitetura](file:///home/ale/Projects/rede-urupe/docs/02-arquitetura/nucleo_urupe_1_2_refinamento.md):** Engenharia descentralizada e especificações do Núcleo Urupê 1.2.
3. **[03-produtos](file:///home/ale/Projects/rede-urupe/docs/03-produtos/deep_dive_product_pillars.md):** Detalhamento dos 3 pilares do ecossistema.
4. **[04-discord](file:///home/ale/Projects/rede-urupe/docs/04-discord/matriz_discord_urupe.md):** Matriz enxuta oficial do servidor do Discord (`#emoji│nome-do-canal`).

---

## 💻 Guia Rápido de Desenvolvimento (`projetos/nucleo-urupe`)

```bash
# 1. Navegar até o projeto
cd projetos/nucleo-urupe

# 2. Executar o backend Go (nucleo-engine)
cd nucleo-engine
go run ./cmd/bot

# 3. Executar o frontend Svelte 5 (nucleo-ui)
cd ../nucleo-ui
bun run dev
```
