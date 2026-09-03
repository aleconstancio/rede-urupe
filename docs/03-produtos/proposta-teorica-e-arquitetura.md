# Proposta Teórica e Arquitetura do App Urupê

> **Plataforma Descentralizada, Local-First e P2P para a Construção do Poder Popular**

---

## 1. Imperativo Político e Visão

As Big Techs corporativas (Meta, Google, Amazon, Uber, iFood) impõem um regime de colonialismo digital, racismo algorítmico e extrativismo de dados sobre a classe trabalhadora periférica. Elas utilizam algoritmos opacos para capturar a ansiedade da precarização, disseminar o fetiche do "empreendedorismo de si" e desmobilizar a organização coletiva.

O **App Urupê** surge como uma contra-infraestrutura tecnológica soberana. Ele é a encarnação em software da **Matriz Ontológica 5x5**, projetado para desintermediar as relações sociais, desprivatizar os meios digitais de comunicação, conectar as lutas territoriais e dar suporte material à construção do **Poder Popular**.

---

## 2. Especificação da Arquitetura Técnica

O aplicativo adota o paradigma **Local-First** (local primeiro) com sincronização peer-to-peer (P2P), garantindo custódia única dos dados pelo usuário e funcionamento pleno sem dependência de internet comercial ou datacenters corporativos.

```
┌────────────────────────────────────────────────────────────────────────┐
│                          CAMADAS DO APP URUPÊ                          │
├────────────────────────────────────────────────────────────────────────┤
│ 1. NÍVEL DE NAVEGAÇÃO INTERFAZ│ Decidim (Deliberação) • PWA/Native UI  │
├──────────────────────────────┼─────────────────────────────────────────┤
│ 2. MÓDULOS DE PRÁXIS (5 EIXOS)│ Agroecologia, Cuidado, Mesh, OpSec, Econs│
├──────────────────────────────┼─────────────────────────────────────────┤
│ 3. INTELIGÊNCIA ARTIFICIAL   │ Edge AI Frugal (SLMs) • Previsão Agroec │
├──────────────────────────────┼─────────────────────────────────────────┤
│ 4. CONSISTÊNCIA & ARMAZÉM    │ CRDTs (Loro) • SQLite Criptografado     │
├──────────────────────────────┼─────────────────────────────────────────┤
│ 5. REDE & TRANSPORTE P2P     │ Libp2p • Matrix P2P • Relays Nostr      │
├──────────────────────────────┼─────────────────────────────────────────┤
│ 6. ENLACE FÍSICO / MESH      │ LibreMesh • Wi-Fi Direct • BLE • Rádio  │
└──────────────────────────────┴─────────────────────────────────────────┘
```

---

## 3. Módulos Operacionais Espelhados nos 5 Eixos Ontológicos

O aplicativo é estruturado em **cinco módulos funcionais** correspondentes aos 5 Eixos da Matriz Ontológica:

```
                          [ APP REDE URUPÊ ]
                                   │
      ┌───────────────┬────────────┼────────────┬───────────────┐
      ▼               ▼            ▼            ▼               ▼
 [ Módulo 1:     [ Módulo 2:   [ Módulo 3:  [ Módulo 4:   [ Módulo 5: 
   TERRA          ÁGUA         AR           FOGO          MICÉLIO 
   Solo & Hosp. ] Cuidado &    Cosmotécnica Autodefesa    Economia dos
                  Saúde ]      & Mesh ]     & Decidim ]   Comuns ]
```

### 3.1. Módulo 1 (Eixo TERRA): Gestão de Solo, Habitação e Agroecologia
* **Mapeamento de Veto Comunal:** Mapeamento offline de imóveis ociosos, áreas de risco e terrenos para destinação popular.
* **OpenStreetMap & Banco de Sementes:** Geolocalização de hortas comunitárias, fontes de água potável e inventário CRDT de sementes crioulas.

### 3.2. Módulo 2 (Eixo ÁGUA): Cozinhas Comunitárias e Saúde Popular
* **Gestão de Comuns Reprodutivos:** Organização de escalas, insumos e logística de cozinhas, lavanderias e creches comunitárias.
* **Rede de Saúde Popular:** Agendamento autônomo de círculos de escuta, fitoterapia e suporte psicossocial sem a indústria farmacêutica.

### 3.3. Módulo 3 (Eixo AR): Cosmotécnica, Comunicação Mesh e IA Frugal
* **Comunicação Off-Grid Mesh:** Troca P2P de mensagens, arquivos e alertas via BLE, Wi-Fi Direct e Relays Nostr locais sem operadora.
* **IA Frugal Ecossocialista (*Edge AI*):** Execução local de modelos SLM quantizados no dispositivo para apoio agroecológico e diagnósticos de saúde comunitária.

### 3.4. Módulo 4 (Eixo FOGO): Autodefesa, OpSec e Deliberação Directa
* **Alerta de Emergência & Autodefesa:** Disparo de notificações anonimizadas contra despejos, violência policial ou ameaças aos territórios.
* **Governança Decidim Embed:** Orçamentos participativos, votação criptográfica e deliberação autônoma de assembleias de bairro.

### 3.5. Módulo 5 (Eixo MICÉLIO): Banco do Tempo e Circuito dos Comuns
* **Bancos de Tempo & Crédito Mútuo:** Registro de trocas não-capitalistas baseadas em horas dedicadas ao Comum.
* **Moedas Sociais (Plataforma e-Dinheiro):** Liquidação offline via QR Code para retenção da riqueza nos territórios sem juros bancários.
* **Biblioteca Teórica Offline (Camada 3):** Armazenamento SQLite/CRDT dos 25 Cadernos de Aprofundamento Epistemológico em Markdown.

---

## 4. Governança Micelar e Federação de Células

* **Autonomia do Micro-Micélio:** Autonomia decisória total de cada célula territorial sobre seus dados e ações.
* **Rede de Confiança (*Web of Trust*):** Atestados de confiança criptográficos trocados presencialmente via leitura de QR Codes.
* **Reputação baseada na Práxis:** Permissão para coordenações regionais fundamentada no histórico verificado de trabalho comunitário.
