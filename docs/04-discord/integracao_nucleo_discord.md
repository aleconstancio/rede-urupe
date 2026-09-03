# 🍄 Arquitetura de Integração: Núcleo Urupê ↔ Discord

> **"A ponte cibernética entre o servidor do Discord da Frente Urupê, a inteligência da mascot Micélia 🍄 (Arandu Engine 7 Estágios) e o Console de Administração Admin."**

---

# 📐 1. Fluxo de Dados Bidirecional

```mermaid
graph TD
    subgraph DISCORD SERVER (Frente Urupê)
        CH_PRESENT[🌱│apresentacoes: Boas-Vindas & Recepção]
        CH_MICORRIZA[🍄│micorriza: Chat Direto com a IA]
        FORUMS[📌 INFORMAÇÕES: 5 Fóruns do Manifesto]
        CH_ANUNCIOS[📢│anuncios: Comunicados Oficiais]
    end

    subgraph NÚCLEO ENGINE (Go Backend)
        DISCORD_HANDLER[handler.go: Escuta Eventos & Slash Cmds]
        ARANDU_PIPELINE[agent_stages.go: 7 Estágios da Arandu Engine]
        SQLITE_DB[(data/nucleo.db: FTS5 & Perfis)]
        SSE_SERVER[server.go: Server-Sent Events / REST API]
    end

    subgraph NÚCLEO UI (Svelte 5 Admin)
        ADMIN_DASH[Centro de Comando: Status do Bot]
        MICELIA_STUDIO[Estúdio da Micélia: Prompts & Memórias]
        CMS_PANEL[CMS: Notícias & Publicações]
    end

    DISCORD_HANDLER <-->|Capta Mensagens & Responde| CH_MICORRIZA
    DISCORD_HANDLER <-->|Registra Perfil & Boas-Vindas| CH_PRESENT
    DISCORD_HANDLER <-->|Indexa Tópicos para Síntese| FORUMS

    DISCORD_HANDLER <-->|Processa Raciocínio| ARANDU_PIPELINE
    ARANDU_PIPELINE <-->|Persiste & Consulta FTS5| SQLITE_DB

    SQLITE_DB <-->|Stream SSE ao vivo| SSE_SERVER
    SSE_SERVER <-->|Visualização & Controle| ADMIN_DASH
    CMS_PANEL -->|Publica Anúncios Diretos| CH_ANUNCIOS
```

---

# 🛠️ 2. Os 4 Pilares da Integração

### 1. 🌾 Recepção Automática (`🌱│apresentacoes`)
- Quando um novo militante entra no servidor e posta uma mensagem no `#🌱│apresentacoes`, o `OnGuildMemberAdd` e `OnMessageCreate`:
  1. Cria ou atualiza o perfil do membro no SQLite (`member_profiles`).
  2. Envia uma mensagem acolhedora de boas-vindas com as diretrizes do manifesto.
  3. Atribui o cargo inicial de militante.

### 2. 🍄 Receptáculo da Micélia (`#🍄│micorriza`)
- Todas as mensagens enviadas no `#🍄│micorriza` (ou marcando a IA) ativam automaticamente a resposta reativa da Micélia.
- A mensagem é processada pelos **7 Estágios da Arandu Engine** (Gater ➔ Memory ➔ Persona Overlay ➔ Planner ➔ Executor ➔ Evaluator ➔ Learning).
- As sínteses e aprendizados são salvos nas cápsulas de memória episódica no SQLite.

### 3. 💬 Síntese Automatizada dos 5 Fóruns do Manifesto (`📌 INFORMAÇÕES`)
- O worker `ForumWorker` analisa os tópicos criados nos 5 fóruns do manifesto (`i-mundo`, `ii-sociedade`, `iii-cosmotecnica`, `iv-praxis`, `v-espirito`).
- Gera um resumo diário dos principais debates e disponibiliza no `#📢│anuncios`, no CMS do Site Público e no Admin (`nucleo-ui`).

### 4. 🎛️ Painel de Governança no Admin (`nucleo-ui`)
- Pelo console Admin do Núcleo Urupê, os coordenadores podem:
  - Monitorar o status do bot Discord em tempo real.
  - Editar os prompts de persona e diretrizes da Micélia.
  - Inspecionar e aprovar propostas de aprendizado contínuo geradas pela IA.
  - Disparar comunicados oficiais do CMS direto no canal `#📢│anuncios`.
