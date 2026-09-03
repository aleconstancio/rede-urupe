<script>
    import { onMount } from 'svelte';
    import { Sparkles, BookOpen, ShieldCheck, Cpu, ArrowRight, Rss, Layers, Globe, HeartHandshake, Search, X, Tag, Calendar, User, CheckCircle2, Compass, Send } from 'lucide-svelte';

    let activePublicTab = $state('noticias'); // 'quem-somos' | 'noticias' | 'projetos'

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

    const projetosList = [
        {
            id: 'app-urupe',
            title: 'App Urupê',
            badge: 'Super-App Popular',
            desc: 'A plataforma soberana da classe trabalhadora. Mensageria criptografada Ed25519, feed autônomo e apoio a mobilizações.',
            icon: Globe,
            color: '#10b981'
        },
        {
            id: 'rizoma',
            title: 'Rizoma Engine',
            badge: 'Motor P2P Local-First',
            desc: 'Rede mesh descentralizada escrita em Rust. Funciona sem internet através de nós comunitários (Rizoma Box).',
            icon: Layers,
            color: '#f59e0b'
        },
        {
            id: 'spore-ops',
            title: 'Spore Ops 🍄',
            badge: 'Guerrilha Memética',
            desc: 'Central de clipping de mídias, agitprop e distribuição de peças de comunicação popular para militantes.',
            icon: Send,
            color: '#ec4899'
        },
        {
            id: 'guara-geo',
            title: 'Guará Geo 🪶',
            badge: 'Inteligência Territorial',
            desc: 'Módulo de análise de dados de satélite para monitorar o uso da terra, microclimas e apoio a brigadas agroecológicas.',
            icon: Compass,
            color: '#06b6d4'
        },
        {
            id: 'jatai-ops',
            title: 'Jataí Ops 🐝',
            badge: 'Autofinanciamento B2B',
            desc: 'Fábrica de agentes robóticos comerciais (Vico, RH) cujo faturamento financia a infraestrutura livre da Frente Urupê.',
            icon: Cpu,
            color: '#8b5cf6'
        }
    ];

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
                <button class="nav-btn {activePublicTab === 'quem-somos' ? 'active' : ''}" onclick={() => activePublicTab = 'quem-somos'}>
                    Quem Somos
                </button>
                <button class="nav-btn {activePublicTab === 'noticias' ? 'active' : ''}" onclick={() => activePublicTab = 'noticias'}>
                    Notícias
                </button>
                <button class="nav-btn {activePublicTab === 'projetos' ? 'active' : ''}" onclick={() => activePublicTab = 'projetos'}>
                    Projetos
                </button>
            </div>
        </div>
    </nav>

    <!-- ABA 1: QUEM SOMOS -->
    {#if activePublicTab === 'quem-somos'}
        <section class="tab-content container">
            <header class="hero-banner">
                <div class="hero-inner">
                    <div class="badge-pill">
                        <Sparkles size={14} class="text-emerald" />
                        <span>Soberania Digital & Ecossocialismo</span>
                    </div>
                    <h1 class="hero-title">
                        Quem Somos.<br />
                        <span class="gradient-text">Frente Popular por uma Comunicação Soberana.</span>
                    </h1>
                    <p class="hero-subtitle">
                        A Frente Urupê é um movimento comunitário e tecnológico soberano. Erguemos ferramentas de código livre e redes descentralizadas para garantir a emancipação da classe trabalhadora, a solidariedade de classe e a defesa da terra.
                    </p>
                </div>
            </header>

            <div class="manifesto-card mt-8">
                <HeartHandshake size={36} class="text-emerald mb-4" />
                <h2 class="font-serif text-3xl font-bold mb-4">O Manifesto Urupê</h2>
                <p class="manifesto-text">
                    "Não aceitamos que o futuro da comunicação humana pertença a corporações que visam o lucro e a vigilância. A Frente Urupê ergue a tecnologia como instrumento de emancipação social, solidariedade comunitária e defesa da terra."
                </p>
            </div>
        </section>
    {/if}

    <!-- ABA 2: NOTÍCIAS (Urupê News & Formação Política) -->
    {#if activePublicTab === 'noticias'}
        <section class="tab-content container">
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
                <div class="state-msg">Carregando mídias e artigos...</div>
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
                                </div>
                            </article>
                        {/each}
                    </div>
                {/if}
            {/if}
        </section>
    {/if}

    <!-- ABA 3: PROJETOS -->
    {#if activePublicTab === 'projetos'}
        <section class="tab-content container">
            <div class="section-title-box">
                <p class="eyebrow">Ecossistema Digital & Terrestre</p>
                <h2>Projetos da Frente Urupê</h2>
                <p class="section-desc">Conheça o conjunto de ferramentas livres, redes P2P e centrais operacionais desenvolvidas para a autonomia tecnológica.</p>
            </div>

            <div class="projects-grid">
                {#each projetosList as prj}
                    <div class="project-card" style="border-top: 4px solid {prj.color}">
                        <div class="prj-header">
                            <span class="prj-badge">{prj.badge}</span>
                        </div>
                        <h3 class="prj-title">{prj.title}</h3>
                        <p class="prj-desc">{prj.desc}</p>
                    </div>
                {/each}
            </div>
        </section>
    {/if}

    <!-- Footer -->
    <footer class="vico-footer">
        <div class="container footer-inner">
            <p>© 2026 Frente Urupê — Tecnologia Livre & Soberana</p>
            <div class="footer-links">
                <button class="ft-btn" onclick={() => activePublicTab = 'quem-somos'}>Quem Somos</button>
                <button class="ft-btn" onclick={() => activePublicTab = 'noticias'}>Notícias</button>
                <button class="ft-btn" onclick={() => activePublicTab = 'projetos'}>Projetos</button>
                <span class="v-tag">Núcleo v1.2</span>
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
        padding: 3rem 1.5rem;
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

    .nav-links { display: flex; gap: 1rem; }
    .nav-btn {
        background: transparent;
        border: none;
        color: var(--muted-foreground);
        font-size: 0.95rem;
        font-weight: 600;
        padding: 0.5rem 1rem;
        border-radius: var(--radius);
        cursor: pointer;
        transition: all 0.2s;
    }

    .nav-btn.active, .nav-btn:hover {
        color: var(--primary);
        background: var(--muted);
    }

    /* Hero */
    .hero-banner {
        padding: 3rem 1.5rem 2rem;
        text-align: center;
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
        font-size: 3.25rem;
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
        font-size: 1.15rem;
        line-height: 1.6;
        color: var(--muted-foreground);
    }

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
        transition: transform 0.2s;
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
    }

    .card-bottom {
        display: flex;
        justify-content: space-between;
        align-items: center;
        border-top: 1px solid var(--border);
        padding-top: 0.75rem;
        font-size: 0.85rem;
    }

    /* Projects Grid */
    .projects-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
        gap: 2rem;
    }

    .project-card {
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 2rem;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .prj-badge {
        font-size: 0.75rem;
        font-weight: 700;
        color: var(--primary);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .prj-title { font-family: var(--font-serif); font-size: 1.5rem; font-weight: 700; }
    .prj-desc { font-size: 0.95rem; color: var(--muted-foreground); line-height: 1.5; }

    /* Manifesto Card */
    .manifesto-card {
        background: oklch(from var(--primary) l c h / 0.08);
        border: 1px solid oklch(from var(--primary) l c h / 0.25);
        border-radius: var(--radius);
        padding: 3.5rem;
        text-align: center;
    }

    .manifesto-text {
        font-family: var(--font-serif);
        font-size: 1.35rem;
        font-style: italic;
        line-height: 1.6;
        max-width: 900px;
        margin: 0 auto;
    }

    /* Footer */
    .vico-footer { border-top: 1px solid var(--border); padding: 2rem 0; background: var(--card); font-size: 0.9rem; color: var(--muted-foreground); }
    .footer-inner { display: flex; justify-content: space-between; align-items: center; }
    .footer-links { display: flex; gap: 1rem; align-items: center; }
    .ft-btn { background: none; border: none; color: var(--foreground); font-weight: 600; cursor: pointer; font-size: 0.9rem; }
    .v-tag { font-size: 0.8rem; background: var(--muted); padding: 0.2rem 0.5rem; border-radius: 4px; }

    .state-msg { padding: 4rem; text-align: center; color: var(--muted-foreground); }
</style>
