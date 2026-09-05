# 🌲 Ramos & Motor da Rede Urupê

> **"O App Urupê é a interface visual soberana na mão do militante. O Rizoma Engine é o motor P2P em Rust que impulsiona o App Urupê e os nós comunitários Rizoma Box."**

---

# 📐 A Relação Entre Núcleo, App Urupê e Rizoma Engine

```mermaid
graph TD
    subgraph 1. 🏛️ RAMO NÚCLEO URUPÊ (QG & GOVERNANÇA)
        NU[🏛️ Núcleo Urupê: Go 1.26 + Svelte 5 + Web2/HTTPS]
        NU --> SPORE[🧫 Spore Ops 🍄: Marketing Digital & Agitprop]
        NU --> JATAI[🐝 Jataí Ops 🐝: Automação B2B & Fundo Urupê]
    end

    subgraph 2. 📱 RAMO APP URUPÊ (SUPER-APP & INTERFACE)
        APP[📱 App Urupê: Interface Flutter GPU no Mobile / Svelte 5 na Web]
    end

    subgraph 3. 🌾 RIZOMA ENGINE (O MOTOR P2P EM RUST CORE)
        RZ[🌾 Rizoma Engine: Core Rust + Criptografia MLS + P2P iroh/libp2p + CRDTs]
        RZ -->|Embarcado via FFI C-ABI| APP
        RZ -->|Firmware Nativo| BOX[📦 Rizoma Box: Nós Comunitários Offline]
    end
```

---

## 🔑 Esclarecimento Arquitetural

1. **O Rizoma continua?**
   - **SIM! O Rizoma é a espinha dorsal de toda a nossa infraestrutura.**
   - O Rizoma **não é uma interface separada**, mas sim o **Rizoma Engine** — a biblioteca central em **Rust** que cuida de toda a criptografia, sincronização de dados CRDT e rede P2P sem servidores corporativos.

2. **Como os 3 Pilares se Conectam:**
   - **🏛️ Núcleo Urupê:** É a plataforma de controle político da **organização ecossocialista** (Go + Svelte 5), gerando a receita (*Jataí Ops*) e a comunicação de rua (*Spore Ops*).
   - **📱 App Urupê:** É o **Super-App popular** que o trabalhador instala no celular (UI Flutter/Svelte).
   - **🌾 Rizoma Engine:** É o **motor P2P em Rust** que roda **dentro do App Urupê** e nos micro-computadores comunitários (**Rizoma Box**).
