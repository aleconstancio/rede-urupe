<script>
    import { onMount } from 'svelte';
    import { Sparkles, BookOpen, ShieldCheck, Cpu, ArrowRight, Rss, Layers, Globe, HeartHandshake, Search, X, Tag, Calendar, User, CheckCircle2 } from 'lucide-svelte';

    let articles = $state([]);
    let loading = $state(true);
    let searchQuery = $state('');
    let activeCategory = $state(null);

    const categories = ['Todos', 'Manifesto', 'Notícia', 'Formação Política', 'Editorial', 'Soberania'];

    async function fetchArticles() {
        try {
            const res = await fetch('/api/cms/articles?public=true');
            if (res.ok) {
                articles = await res.json();
            }
        } catch (e) {
            console.error('Error fetching public articles:', e);
        } finally {
            loading = false;
        }
    }

    const filteredArticles = $derived(
        articles.filter(item => {
            const matchesCat = !activeCategory || activeCategory === 'Todos' || item.category === activeCategory;
            const matchesSearch = !searchQuery ||
                item.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                (item.summary && item.summary.toLowerCase().includes(searchQuery.toLowerCase()));
            return matchesCat && matchesSearch;
        })
    );

    const bigFeatured = $derived(filteredArticles.length > 0 ? filteredArticles[0] : null);
    const regularArticles = $derived(filteredArticles.length > 1 ? filteredArticles.slice(1) : []);

    onMount(() => {
        fetchArticles();
    });
</script>

<div class="vico-style-landing">
    <!-- Parallax Glow Background -->
    <div class="parallax-glow" aria-hidden="true">
        <div class="glow-circle glow-1"></div>
        <div class="glow-circle glow-2"></div>
    </div>

    <!-- Navigation Header -->
    <nav class="vico-nav">
        <div class="nav-container">
            <div class="brand flex items-center gap-3">
                <div class="brand-logo">🍄</div>
                <span class="brand-title">Frente Urupê</span>
            </div>
            <div class="nav-links">
                <a href="#manifesto" class="nav-link">Manifesto</a>
                <a href="#noticias" class="nav-link">Notícias & Formação</a>
                <a href="#pilares" class="nav-link">A Trindade</a>
                <a href="#discord" class="nav-link">Rede Discord</a>
            </div>
        </div>
    </nav>

    <!-- Hero Section Inspired by Vico -->
    <header class="hero-banner">
        <div class="hero-inner">
            <div class="badge-pill">
                <Sparkles size={14} class="text-emerald" />
                <span>Soberania Digital & Ecossocialismo</span>
            </div>
            <h1 class="hero-title">
                Comunicação Soberana.<br />
                <span class="gradient-text">Sem Big Techs, Sem Algoritmos Extrativistas.</span>
            </h1>
            <p class="hero-subtitle">
                A Frente Urupê é uma praça pública digital descentralizada e ecossocialista. Tecnologia livre construída pela classe trabalhadora para a emancipação popular e defesa da terra.
            </p>
            <div class="hero-ctas">
                <a href="#noticias" class="btn btn-primary">
                    <span>Explorar Notícias & Manifesto</span>
                    <ArrowRight size={18} />
                </a>
                <a href="#pilares" class="btn btn-glass">
                    <BookOpen size={18} />
                    <span>Conhecer a Trindade Urupê</span>
                </a>
            </div>

            <div class="trust-strip">
                <div class="trust-item"><CheckCircle2 size={16} class="text-emerald" /> <span>100% Código Livre & Open-Hardware</span></div>
                <div class="trust-item"><CheckCircle2 size={16} class="text-emerald" /> <span>Zero Rastreadores ou Anúncios</span></div>
                <div class="trust-item"><CheckCircle2 size={16} class="text-emerald" /> <span>Rede Mesh P2P Offline</span></div>
            </div>
        </div>
    </header>

    <!-- News & Editorial Section Inspired by VicoNews -->
    <section id="noticias" class="news-section container">
        <div class="section-title-box">
            <p class="eyebrow">Imprensa Popular & Teoria</p>
            <h2>Urupê News & Formação Política</h2>
            <p class="section-desc">Publicações oficiais assinadas criptograficamente pela Frente Urupê e replicadas na malha Rizoma P2P.</p>
        </div>

        <!-- Search and FilterBar (VicoNews Style) -->
        <div class="news-controls">
            <div class="search-box">
                <Search size={18} class="search-icon" />
                <input type="text" placeholder="Buscar no acervo de artigos e manifestos..." bind:value={searchQuery} class="search-input" />
                {#if searchQuery}
                    <button onclick={() => searchQuery = ''} class="clear-btn"><X size={16} /></button>
                {/if}
            </div>

            <!-- Category Filter Chips -->
            <div class="filter-chips">
                {#each categories as cat}
                    <button class="chip {activeCategory === cat || (!activeCategory && cat === 'Todos') ? 'active' : ''}"
                        onclick={() => activeCategory = cat === 'Todos' ? null : cat}>
                        {cat}
                    </button>
                {/each}
            </div>
        </div>

        {#if loading}
            <div class="skeleton-grid">
                <div class="skeleton-featured"></div>
                <div class="skeleton-card"></div>
                <div class="skeleton-card"></div>
            </div>
        {:else if filteredArticles.length === 0}
            <div class="empty-news">
                <Search size={32} class="text-muted" />
                <p>Nenhuma publicação encontrada para esses critérios.</p>
            </div>
        {:else}
            <!-- Big Featured Card (VicoNews Spotlight) -->
            {#if bigFeatured}
                <div class="big-featured-card">
                    <div class="featured-badge">
                        <Sparkles size={14} /> Destaque Principal
                    </div>
                    <div class="featured-meta">
                        <span class="cat-tag">{bigFeatured.category}</span>
                        <span class="meta-item"><Calendar size={14} /> {bigFeatured.created_at}</span>
                        <span class="meta-item"><User size={14} /> {bigFeatured.author}</span>
                    </div>
                    <h3 class="featured-title">{bigFeatured.title}</h3>
                    <p class="featured-summary">{bigFeatured.summary || bigFeatured.content.slice(0, 200) + '...'}</p>
                    <div class="featured-footer">
                        {#if bigFeatured.p2p_signed}
                            <span class="p2p-seal">
                                <ShieldCheck size={16} /> Assinado P2P Ed25519
                            </span>
                        {/if}
                    </div>
                </div>
            {/if}

            <!-- Secondary Articles Grid -->
            {#if regularArticles.length > 0}
                <div class="articles-grid">
                    {#each regularArticles as item}
                        <article class="vico-news-card">
                            <div class="card-top">
                                <span class="cat-tag">{item.category}</span>
                                <span class="date">{item.created_at}</span>
                            </div>
                            <h4 class="card-title">{item.title}</h4>
                            <p class="card-summary">{item.summary || item.content.slice(0, 120) + '...'}</p>
                            <div class="card-bottom">
                                <span class="author">Por <b>{item.author}</b></span>
                                {#if item.p2p_signed}
                                    <span class="p2p-icon" title="Verificado P2P">🛡️ P2P</span>
                                {/if}
                            </div>
                        </article>
                    {/each}
                </div>
            {/if}
        {/if}
    </section>

    <!-- Pillars Section -->
    <section id="pilares" class="pillars-section container">
        <div class="section-title-box">
            <p class="eyebrow">A Arquitetura de 3 Camadas</p>
            <h2>A Trindade Urupê</h2>
        </div>

        <div class="pillars-grid">
            <div class="pillar-card">
                <div class="icon-wrap text-emerald"><ShieldCheck size={28} /></div>
                <h3>1. Núcleo Urupê</h3>
                <p>Plataforma de governança, bot do Discord, IA <b>Micélia 🍄</b> e gestão do portal público de comunicação.</p>
            </div>
            <div class="pillar-card">
                <div class="icon-wrap text-amber"><Layers size={28} /></div>
                <h3>2. Rizoma Engine</h3>
                <p>O motor P2P local-first em Rust. Criptografia Ed25519, rede mesh offline e nós sementes comunitários (*Rizoma Box*).</p>
            </div>
            <div class="pillar-card">
                <div class="icon-wrap text-teal"><Globe size={28} /></div>
                <h3>3. App Urupê</h3>
                <p>O Super-App popular da classe trabalhadora: mensageria soberana, feed plural, entregas/Uber solidários e mídias.</p>
            </div>
        </div>
    </section>

    <!-- Manifesto Highlight Section -->
    <section id="manifesto" class="manifesto-section container">
        <div class="manifesto-card">
            <HeartHandshake size={36} class="text-emerald mb-4" />
            <h2 class="font-serif text-3xl font-bold mb-4">Soberania Tecnológica Popular</h2>
            <p class="manifesto-text">
                "Não aceitamos que o futuro da comunicação humana pertença a corporações que visam o lucro e a vigilância. A Frente Urupê ergue a tecnologia como instrumento de emancipação social, solidariedade comunitária e defesa da terra."
            </p>
        </div>
    </section>

    <!-- Footer -->
    <footer class="vico-footer">
        <div class="container footer-inner">
            <p>© 2026 Frente Urupê — Tecnologia Livre & Soberana</p>
            <div class="footer-links">
                <a href="#noticias">Urupê News</a>
                <a href="#manifesto">Manifesto</a>
                <span class="v-tag">Núcleo v2.0</span>
            </div>
        </div>
    </footer>
</div>

<style>
    .vico-style-landing {
        min-height: 100vh;
        background: var(--background);
        color: var(--foreground);
        font-family: var(--font-sans);
        position: relative;
        overflow-x: hidden;
    }

    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 4rem 1.5rem;
    }

    /* Parallax Glow */
    .parallax-glow {
        position: absolute;
        top: 0;
        left: 50%;
        transform: translateX(-50%);
        width: 100%;
        max-width: 1400px;
        height: 600px;
        pointer-events: none;
        z-index: 0;
    }

    .glow-circle {
        position: absolute;
        border-radius: 9999px;
        filter: blur(120px);
        opacity: 0.15;
    }

    .glow-1 {
        width: 500px;
        height: 500px;
        background: var(--primary);
        top: -100px;
        left: 10%;
    }

    .glow-2 {
        width: 400px;
        height: 400px;
        background: #10b981;
        top: 100px;
        right: 15%;
    }

    /* Nav */
    .vico-nav {
        position: sticky;
        top: 0;
        z-index: 40;
        background: oklch(from var(--background) l c h / 0.8);
        backdrop-filter: blur(16px);
        border-bottom: 1px solid var(--border);
    }

    .nav-container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 1rem 1.5rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .brand-logo { font-size: 1.5rem; }
    .brand-title { font-family: var(--font-serif); font-size: 1.35rem; font-weight: 700; }

    .nav-links { display: flex; gap: 2rem; }
    .nav-link { color: var(--muted-foreground); text-decoration: none; font-size: 0.9rem; font-weight: 500; transition: color 0.2s; }
    .nav-link:hover { color: var(--primary); }

    /* Hero */
    .hero-banner {
        padding: 6rem 1.5rem 4rem;
        text-align: center;
        position: relative;
        z-index: 1;
    }

    .hero-inner {
        max-width: 850px;
        margin: 0 auto;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    .badge-pill {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.4rem 1rem;
        background: oklch(from var(--primary) l c h / 0.12);
        border: 1px solid oklch(from var(--primary) l c h / 0.25);
        border-radius: 9999px;
        font-size: 0.85rem;
        font-weight: 600;
        color: var(--primary);
        margin-bottom: 1.5rem;
    }

    .hero-title {
        font-family: var(--font-serif);
        font-size: 3.5rem;
        font-weight: 700;
        line-height: 1.15;
        letter-spacing: -0.02em;
        margin-bottom: 1.25rem;
    }

    .gradient-text {
        background: linear-gradient(135deg, var(--primary), #10b981);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
    }

    .hero-subtitle {
        font-size: 1.2rem;
        line-height: 1.6;
        color: var(--muted-foreground);
        margin-bottom: 2.5rem;
    }

    .hero-ctas {
        display: flex;
        gap: 1rem;
        flex-wrap: wrap;
        justify-content: center;
        margin-bottom: 3rem;
    }

    .btn {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.8rem 1.75rem;
        border-radius: var(--radius);
        font-weight: 600;
        text-decoration: none;
        transition: all 0.2s ease;
        cursor: pointer;
    }

    .btn-primary { background: var(--primary); color: var(--background); }
    .btn-primary:hover { filter: brightness(1.1); transform: translateY(-2px); }

    .btn-glass {
        background: oklch(from var(--card) l c h / 0.6);
        backdrop-filter: blur(10px);
        border: 1px solid var(--border);
        color: var(--foreground);
    }
    .btn-glass:hover { background: var(--muted); }

    .trust-strip {
        display: flex;
        gap: 2rem;
        flex-wrap: wrap;
        justify-content: center;
        font-size: 0.85rem;
        color: var(--muted-foreground);
    }
    .trust-item { display: flex; align-items: center; gap: 0.4rem; }

    /* News Section */
    .section-title-box { text-align: center; margin-bottom: 3rem; }
    .eyebrow { font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: var(--primary); margin-bottom: 0.5rem; }
    .section-title-box h2 { font-family: var(--font-serif); font-size: 2.5rem; font-weight: 700; }
    .section-desc { color: var(--muted-foreground); max-width: 600px; margin: 0.5rem auto 0; font-size: 1rem; }

    /* Controls (VicoNews Search & Filter) */
    .news-controls {
        max-width: 900px;
        margin: 0 auto 3rem;
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
        align-items: center;
    }

    .search-box {
        position: relative;
        width: 100%;
        max-width: 600px;
    }

    .search-input {
        width: 100%;
        padding: 0.85rem 1rem 0.85rem 2.8rem;
        background: oklch(from var(--card) l c h / 0.6);
        backdrop-filter: blur(12px);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        color: var(--foreground);
        font-size: 0.95rem;
    }

    .search-icon {
        position: absolute;
        left: 1rem;
        top: 50%;
        transform: translateY(-50%);
        color: var(--muted-foreground);
    }

    .clear-btn {
        position: absolute;
        right: 1rem;
        top: 50%;
        transform: translateY(-50%);
        background: none;
        border: none;
        color: var(--muted-foreground);
        cursor: pointer;
    }

    .filter-chips {
        display: flex;
        gap: 0.5rem;
        flex-wrap: wrap;
        justify-content: center;
    }

    .chip {
        padding: 0.4rem 1rem;
        border-radius: 9999px;
        font-size: 0.85rem;
        font-weight: 600;
        background: var(--card);
        border: 1px solid var(--border);
        color: var(--muted-foreground);
        cursor: pointer;
        transition: all 0.2s;
    }

    .chip.active, .chip:hover {
        background: var(--primary);
        color: var(--background);
        border-color: var(--primary);
    }

    /* Big Featured Card */
    .big-featured-card {
        background: linear-gradient(135deg, oklch(from var(--card) l c h / 0.9), oklch(from var(--primary) l c h / 0.08));
        border: 1px solid oklch(from var(--primary) l c h / 0.3);
        border-radius: var(--radius);
        padding: 3rem;
        margin-bottom: 2.5rem;
        position: relative;
    }

    .featured-badge {
        display: inline-flex;
        align-items: center;
        gap: 0.4rem;
        background: var(--primary);
        color: var(--background);
        font-size: 0.75rem;
        font-weight: 700;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        margin-bottom: 1.25rem;
    }

    .featured-meta {
        display: flex;
        gap: 1.25rem;
        align-items: center;
        font-size: 0.85rem;
        color: var(--muted-foreground);
        margin-bottom: 1rem;
    }

    .cat-tag { font-weight: 700; color: var(--primary); }
    .meta-item { display: flex; align-items: center; gap: 0.4rem; }

    .featured-title {
        font-family: var(--font-serif);
        font-size: 2.25rem;
        font-weight: 700;
        margin-bottom: 1rem;
        line-height: 1.25;
    }

    .featured-summary {
        font-size: 1.1rem;
        line-height: 1.6;
        color: var(--muted-foreground);
        margin-bottom: 1.75rem;
    }

    .p2p-seal {
        display: inline-flex;
        align-items: center;
        gap: 0.4rem;
        font-size: 0.8rem;
        font-weight: 700;
        color: var(--primary);
    }

    /* Articles Grid */
    .articles-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
        gap: 2rem;
    }

    .vico-news-card {
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 1.75rem;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        transition: transform 0.2s, border-color 0.2s;
    }

    .vico-news-card:hover {
        transform: translateY(-4px);
        border-color: var(--primary);
    }

    .card-top {
        display: flex;
        justify-content: space-between;
        font-size: 0.8rem;
        margin-bottom: 0.75rem;
    }

    .card-title {
        font-family: var(--font-serif);
        font-size: 1.25rem;
        font-weight: 600;
        margin-bottom: 0.75rem;
    }

    .card-summary {
        font-size: 0.9rem;
        color: var(--muted-foreground);
        line-height: 1.5;
        margin-bottom: 1.5rem;
        flex-grow: 1;
    }

    .card-bottom {
        display: flex;
        justify-content: space-between;
        align-items: center;
        border-top: 1px solid var(--border);
        padding-top: 0.75rem;
        font-size: 0.85rem;
    }

    .p2p-icon { font-size: 0.75rem; font-weight: 700; color: var(--primary); }

    /* Pillars */
    .pillars-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 2rem;
    }

    .pillar-card {
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 2rem;
    }

    .icon-wrap { margin-bottom: 1rem; }
    .pillar-card h3 { font-family: var(--font-serif); font-size: 1.35rem; font-weight: 600; margin-bottom: 0.5rem; }
    .pillar-card p { font-size: 0.95rem; color: var(--muted-foreground); line-height: 1.5; }

    /* Manifesto */
    .manifesto-card {
        background: oklch(from var(--primary) l c h / 0.08);
        border: 1px solid oklch(from var(--primary) l c h / 0.25);
        border-radius: var(--radius);
        padding: 3.5rem;
        text-align: center;
    }

    .manifesto-text {
        font-family: var(--font-serif);
        font-size: 1.45rem;
        font-style: italic;
        line-height: 1.6;
        max-width: 900px;
        margin: 0 auto;
    }

    /* Footer */
    .vico-footer { border-top: 1px solid var(--border); padding: 2rem 0; background: var(--card); font-size: 0.9rem; color: var(--muted-foreground); }
    .footer-inner { display: flex; justify-content: space-between; align-items: center; }
    .footer-links { display: flex; gap: 1.5rem; align-items: center; }
    .footer-links a { color: var(--foreground); text-decoration: none; }
    .v-tag { font-size: 0.8rem; background: var(--muted); padding: 0.2rem 0.5rem; border-radius: 4px; }
</style>
