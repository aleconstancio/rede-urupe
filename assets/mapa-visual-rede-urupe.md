# Mapa Visual de Arquitetura Sistêmica — Rede Urupê

> **Visualização Diagramática do Ecossistema, da Matriz Ontológica 5x5 e da Pilha Técnica**

---

## 1. A Matriz Ontológica 5x5 e as 3 Camadas Pedagógicas

O diagrama abaixo ilustra a estrutura fractal de formação política da Rede Urupê: a apreensão holística dos **5 Eixos Ontológicos** (Camada 1), o desdobramento nas **25 Subcategorias Universais** (Camada 2) e o aprofundamento nos **Cadernos Teóricos** (Camada 3).

```mermaid
graph TD
    subgraph C1 ["CAMADA 1: SÍNTESE HOLÍSTICA (Os 5 Eixos)"]
        E1["🪵 Eixo I: TERRA (Matéria & Economia)"]
        E2["💧 Eixo II: ÁGUA (Corpo & Reprodução)"]
        E3["🌬️ Eixo III: AR (Intelecto & Técnica)"]
        E4["🔥 Eixo IV: FOGO (Poder & Práxis)"]
        E5["🍄 Eixo V: MICÉLIO (Espírito & Comum)"]
    end

    subgraph C2 ["CAMADA 2: 25 SUBCATEGORIAS ONTOLÓGICAS (Formato Dual)"]
        S1["1.1 Metabolismo | 1.2 Espacialidade | 1.3 Infraestrutura | 1.4 Forma-Valor | 1.5 Capital Financeiro"]
        S2["2.1 Corporeidade | 2.2 Racialidade | 2.3 Reprodução | 2.4 Subjetividade | 2.5 Ciclos de Vida"]
        S3["3.1 Linguagem | 3.2 Intelecto Geral | 3.3 Cosmotécnica | 3.4 Consciência | 3.5 Memória"]
        S4["4.1 Forma-Estado | 4.2 Forma-Jurídica | 4.3 Coerção | 4.4 Autodefesa | 4.5 Poder Popular"]
        S5["5.1 Tempo Livre | 5.2 Alma Humana | 5.3 Reciprocidade | 5.4 O Riso | 5.5 Utopia Concreta"]
    end

    subgraph C3 ["CAMADA 3: NÚCLEO TEÓRICO DALSO (Cadernos de Aprofundamento)"]
        T1["Marx, Kurz, Postone, Saito, Foster, Marini, Lélia, Pasquinelli, Yuk Hui, Benjamin, Saviani"]
    end

    E1 --> S1
    E2 --> S2
    E3 --> S3
    E4 --> S4
    E5 --> S5

    S1 --> T1
    S2 --> T1
    S3 --> T1
    S4 --> T1
    S5 --> T1

    style C1 fill:#1e293b,stroke:#38bdf8,stroke-width:2px,color:#fff
    style C2 fill:#14532d,stroke:#22c55e,stroke-width:2px,color:#fff
    style C3 fill:#4c1d95,stroke:#a855f7,stroke-width:2px,color:#fff
```

---

## 2. O Ecossistema Micelar da Rede Urupê

Separacão de papéis entre a **Rede Urupê** (teoria, software e infraestrutura P2P) e as **Organizações Parceiras** (movimentos sociais e populares).

```mermaid
graph TD
    subgraph COLETIVO_URUPÊ ["REDE URUPÊ (100% Digital)"]
        A[Comuna Cibernética & Comunidade Discord] -->|Fase 1| B[Engenharia do App Urupê Local-First/P2P]
        B -->|Fase 2| C[Rede Descentralizada de Infraestrutura]
    end

    subgraph MOVIMENTOS_PARCEIROS ["ORGANIZAÇÕES PARCEIRAS ('Os Organizados pelo App')"]
        D[MST - Reforma Agrária & Sementes Crioulas]
        E[MTST - Cozinhas Solidárias & Hortas Urbanas]
        F[CONAQ - Territórios Quilombolas & Alerta Climático]
        G[APIB - Povos Indígenas & Etnoagroecologia]
        H[Cooperativas de Plataforma & Economia Solidária]
    end

    C -->|Provimento de Software & Protocolos| D
    C -->|Provimento de Software & Protocolos| E
    C -->|Provimento de Software & Protocolos| F
    C -->|Provimento de Software & Protocolos| G
    C -->|Provimento de Software & Protocolos| H

    style COLETIVO_URUPÊ fill:#1e293b,stroke:#38bdf8,stroke-width:2px,color:#fff
    style MOVIMENTOS_PARCEIROS fill:#14532d,stroke:#22c55e,stroke-width:2px,color:#fff
```

---

## 3. A Pilha de Cosmotécnica & Arquitetura do App Urupê

```mermaid
graph BT
    L6["Camada 6: Enlace Físico & Mesh (LibreMesh, Wi-Fi Direct, BLE, Rádio)"]
    L5["Camada 5: Transporte P2P (Libp2p, Matrix P2P, Relays Nostr Locais)"]
    L4["Camada 4: Armazenamento & Consistência (CRDTs Loro, SQLite Criptografado)"]
    L3["Camada 3: Inteligência Artificial Frugal (Edge AI, SLMs Locais Offline)"]
    L2["Camada 2: Regras & Economia Solidária (Banco do Tempo, e-Dinheiro/Mumbuca)"]
    L1["Camada 1: Interface & Governança (Framework Decidim Embed, PWA UI)"]

    L6 --> L5
    L5 --> L4
    L4 --> L3
    L3 --> L2
    L2 --> L1

    style L1 fill:#0f172a,stroke:#38bdf8,color:#fff
    style L2 fill:#0f172a,stroke:#38bdf8,color:#fff
    style L3 fill:#065f46,stroke:#10b981,color:#fff
    style L4 fill:#0f172a,stroke:#38bdf8,color:#fff
    style L5 fill:#0f172a,stroke:#38bdf8,color:#fff
    style L6 fill:#0f172a,stroke:#38bdf8,color:#fff
```

---

## 4. O Ciclo da Desintoxicação Algorítmica e Práxis Freiriana

```mermaid
flowchart LR
    subgraph REDES_CORPORATIVAS ["Plataformas das Big Techs"]
        X1[Vício em Dopamina & Telas] --> X2[Algoritmo de Ódio & Pânico Moral]
        X2 --> X3[Ideologia do 'Empreendedorismo de Si']
        X3 --> X4[Isolamento, Ansiedade & Passividade]
    end

    subgraph TRANSIÇÃO_URUPÊ ["Práxis com App Urupê"]
        Y1[Custódia Local-First & Foco] --> Y2[Círculos Digitais de Escuta Freiriana]
        Y2 --> Y3[Conscientização Anti-Ressentimento]
        Y3 --> Y4[Coordenação P2P de Ação Direta Comunitária]
    end

    X4 -.->|Desintoxicação Digital| Y1

    style REDES_CORPORATIVAS fill:#7f1d1d,stroke:#ef4444,color:#fff
    style TRANSIÇÃO_URUPÊ fill:#14532d,stroke:#22c55e,color:#fff
```
