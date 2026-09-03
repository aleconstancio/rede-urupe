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

# 🛠️ 2. Os 4 Pilares da Integração Refinada

### 1. 🌱 Aprovação Manual no `#🌱│apresentacoes`
- O canal `#🌱│apresentacoes` é o espaço de recepção e apresentação dos novos militantes.
- **Aprovação 100% Manual por Administradores:** A Micélia não realiza triagem automatizada ou liberações automáticas de cargos. Os administradores e moderadores humanos leem as apresentações e atribuem os cargos manualmente.

### 2. 🍄 Atuação da Micélia 🍄 em Todos os Canais
- A mascot **Micélia 🍄** possui contexto conceitual completo das 4 categorias de canais do servidor:
  - **Canais Comuns** (`chat-comum`, `chat-memetico`, `chat-serio`, `micorriza`)
  - **Canais Culturais** (`noticias`, `fotografia`, `videos`, `jogos`, `cinema`, `musica`, `literatura`)
  - **Frentes de Estudo** (`filosofia`, `politica`, `pedagogia`, `ecologia`, `historia`, `humanidades`, `psicologia`, `tecnologia`, `engenharia`)
  - **Frentes de Atuação** (`nucleo-urupe`, `rizoma`, `app-urupe`, `spore-ops`, `jatai-ops`)
- A Micélia responde e interage reactivamente em qualquer um desses canais quando mencionada, chamada ou quando acionada por perguntas de militantes, respeitando o tom e o domínio de cada canal.

### 3. 📌 Fóruns de Informações Fixos & Sincronizados com o Manifesto
- Os 5 Fóruns Temáticos (`i-mundo`, `ii-sociedade`, `iii-cosmotecnica`, `iv-praxis`, `v-espirito`) são **fixos, oficiais e informativos**.
- **SEM Sínteses Diárias:** Não há geração de resumos automatizados diários nestes fóruns.
- **Atualização Automática via Manifesto:** Conforme o Manifesto da Frente Urupê é editado ou atualizado no CMS/Núcleo Urupê, a Micélia e o backend sincronizam automaticamente os tópicos oficiais de cada capítulo do manifesto dentro dos fóruns.

### 4. 🎛️ Painel de Governança no Admin (`nucleo-ui`)
- Pelo console Admin do Núcleo Urupê, os coordenadores podem:
  - Monitorar a saúde e latência do bot Discord em tempo real.
  - Editar os prompts de persona e diretrizes da Micélia.
  - Sincronizar atualizações do Manifesto com os 5 Fóruns do Discord.
