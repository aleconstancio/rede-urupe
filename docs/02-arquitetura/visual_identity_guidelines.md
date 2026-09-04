# 🎨 Manual Mestre de Identidade Visual: Gruvbox Solarpunk

> **"A estética da Frente Urupê une a ancestralidade da terra, a sobriedade terrosa e a soberania tecnológica popular."**

---

## 🏛️ 1. Princípios da Identidade Visual

1. **Orgânica & Terrosa:** Evitamos o preto estéril e o roxo/azul corporativo de Big Techs. A nossa paleta é ancorada em tons de solo, folhagem, semente e sol.
2. **Soberania sem Ruído:** Interfaces e pôsteres limpos, sem caixas pesadas, sem poluição visual e sem emojis informais em materiais executivos.
3. **Tipografia Trina:**
   - **Cormorant Garamond:** Títulos e chamadas históricas (elegância sacra e editorial).
   - **Source Sans 3:** Corpo de texto, artigos e leitura fluida.
   - **IBM Plex Mono:** Dados técnicos, badges, métricas e soberania digital.

---

## 🎨 2. Paleta de Cores Oficial (Espaço de Cor OKLCH)

Utilizamos o espaço de cor **OKLCH** para garantir transições de cor naturais, alto contraste e fidelidade de impressão e tela:

```css
/* Palette Gruvbox Solarpunk Oficial */
:root {
    /* Fundos & Superfícies Terrosas */
    --color-bg-base:        oklch(0.18 0.02 65);   /* #1d2021 - Carvão Terroso Profundo */
    --color-bg-subtle:      oklch(0.22 0.03 130);  /* Verde Escuro de Clorofila Subterrânea */
    --color-bg-darker:      oklch(0.14 0.02 50);   /* Carvão de Solo Fechado */

    /* Acentos & Elementos Ativos */
    --color-primary:        oklch(0.75 0.18 135);  /* #b8bb26 - Verde Clorofila Vivo */
    --color-solar:          oklch(0.82 0.16 85);   /* #fabd2f - Ocre Solar */
    --color-terracotta:     oklch(0.70 0.18 45);   /* #fe8019 - Terracota Argila */
    --color-teal:           oklch(0.76 0.12 160);  /* #8ec07c - Verde Mineral */

    /* Textos & Leitura */
    --color-text-main:      oklch(0.92 0.03 85);   /* #ebdbb2 - Pergaminho Creme */
    --color-text-muted:     oklch(0.70 0.03 75);   /* #a89984 - Rocha Mútua */
    --color-border:         oklch(0.30 0.02 65);   /* #3c3836 - Linha Terrosa Sutil */
}
```

---

## 📐 3. Regras de Diagramação & Pôsteres Executivos

- **Gradiente de Fundo:** Fundo sempre em `radial-gradient` sutil em OKLCH:
  `radial-gradient(circle at 15% 20%, oklch(0.22 0.04 135) 0%, oklch(0.18 0.02 65) 45%, oklch(0.14 0.02 50) 100%)`.
- **Diagramas de Branching:** Árvores com bifurcação da esquerda para a direita (L-to-R), interligadas por conectores vetoriais em SVG curvos sutilmente coloridos.
- **Iconografia:** Zero emojis em documentos oficiais. Apenas ícones vetoriais SVG de linha limpa (`stroke-width: 2`).
- **Segunda Pessoa:** Textos de convocação escritos em segunda pessoa (*tu/teu/tua*) com tom inspirador, sem jargões computacionais pesados.
