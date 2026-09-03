<script>
    import { onMount } from 'svelte';
    import { Sparkles, BookOpen, ShieldCheck, Cpu, ArrowRight, Rss, Layers, Globe, HeartHandshake } from 'lucide-svelte';

    let articles = $state([]);
    let loading = $state(true);

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

    onMount(() => {
        fetchArticles();
    });
</script>

<div class="public-landing">
    <!-- Hero Section -->
    <header class="hero-section">
        <div class="hero-content">
            <div class="badge-tag">
                <Sparkles size={16} class="text-emerald" />
                <span>Pela Soberania Digital e Emancipação Popular</span>
            </div>
            <h1 class="hero-title">Frente Urupê</h1>
            <p class="hero-subtitle">
                Uma praça pública digital livre, ecossocialista e descentralizada. Sem algoritmos extrativistas, sem anúncios e com tecnologia soberana construída pelo povo para o povo.
            </p>
            <div class="hero-actions">
                <a href="#manifesto" class="btn btn-primary">
                    <span>Ler o Manifesto</span>
                    <ArrowRight size={18} />
                </a>
                <a href="#artigos" class="btn btn-outline">
                    <BookOpen size={18} />
                    <span>Notícias & Artigos</span>
                </a>
            </div>
        </div>
    </header>

    <!-- Pilares Section -->
    <section class="pillars-section container">
        <div class="section-header">
            <p class="eyebrow">A Trindade de Produtos</p>
            <h2>Tecnologia Autônoma & Antifrágil</h2>
        </div>

        <div class="pillars-grid">
            <div class="pillar-card">
                <div class="card-icon text-emerald">
                    <ShieldCheck size={32} />
                </div>
                <h3>1. Núcleo Urupê</h3>
                <p>Central de comando, governança comunitária 5x5, inteligência da mascot IA <b>Micélia 🍄</b> e comunicação institucional da Frente Urupê.</p>
            </div>

            <div class="pillar-card">
                <div class="card-icon text-amber">
                    <Layers size={32} />
                </div>
                <h3>2. Rizoma Engine</h3>
                <p>O protocolo P2P local-first em Rust. Garantia de soberania de dados, encriptação Ed25519 e operação offline em malha mesh durante atos públicos.</p>
            </div>

            <div class="pillar-card">
                <div class="card-icon text-teal">
                    <Globe size={32} />
                </div>
                <h3>3. App Urupê</h3>
                <p>O Super-App da classe trabalhadora: mensageria instantânea soberana, feed plural, mobilidade/entregas solidárias (Uber/iFood) e mídias integradas.</p>
            </div>
        </div>
    </section>

    <!-- Manifesto Highlight Section -->
    <section id="manifesto" class="manifesto-section container">
        <div class="manifesto-box">
            <div class="manifesto-header">
                <HeartHandshake size={32} class="text-emerald" />
                <h2>Nosso Compromisso Histórico</h2>
            </div>
            <p class="manifesto-text">
                "Não aceitamos ser meros produtos de algoritmos corporativos desenhados para extrair atenção e gerar ansiedade. A Frente Urupê retoma a tecnologia como instrumento de emancipação social, solidariedade comunitária e defesa da terra."
            </p>
        </div>
    </section>

    <!-- Articles Section -->
    <section id="artigos" class="articles-section container">
        <div class="section-header">
            <p class="eyebrow">Imprensa & Jornalismo Popular</p>
            <h2>Notícias & Publicações Oficiais</h2>
        </div>

        {#if loading}
            <div class="loading-state">Carregando artigos oficiais...</div>
        {:else if articles.length === 0}
            <div class="empty-state">Nenhum artigo publicado no momento.</div>
        {:else}
            <div class="articles-grid">
                {#each articles as item}
                    <article class="article-card">
                        <div class="card-meta">
                            <span class="category-badge">{item.category}</span>
                            <span class="date">{item.created_at}</span>
                        </div>
                        <h3>{item.title}</h3>
                        <p>{item.summary}</p>
                        <div class="card-footer">
                            <span class="author">Por <b>{item.author}</b></span>
                            {#if item.p2p_signed}
                                <span class="p2p-tag" title="Assinado criptograficamente na malha Rizoma P2P">
                                    🛡️ Assinado P2P
                                </span>
                            {/if}
                        </div>
                    </article>
                {/each}
            </div>
        {/if}
    </section>

    <!-- Footer -->
    <footer class="public-footer">
        <div class="container footer-content">
            <p>© 2026 Frente Urupê — Tecnologia Livre & Soberana</p>
            <div class="footer-links">
                <a href="#manifesto">Manifesto</a>
                <a href="#artigos">Notícias</a>
                <span class="version">Núcleo v2.0</span>
            </div>
        </div>
    </footer>
</div>

<style>
    .public-landing {
        min-height: 100vh;
        background: var(--background);
        color: var(--foreground);
        font-family: var(--font-sans);
    }

    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 4rem 1.5rem;
    }

    /* Hero */
    .hero-section {
        padding: 6rem 1.5rem 4rem;
        text-align: center;
        background: radial-gradient(circle at top, oklch(from var(--primary) l c h / 0.12), transparent 70%);
        border-bottom: 1px solid var(--border);
    }

    .hero-content {
        max-width: 800px;
        margin: 0 auto;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    .badge-tag {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.4rem 1rem;
        background: oklch(from var(--primary) l c h / 0.15);
        border: 1px solid oklch(from var(--primary) l c h / 0.3);
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
        letter-spacing: -0.02em;
        margin-bottom: 1rem;
        color: var(--foreground);
    }

    .hero-subtitle {
        font-size: 1.25rem;
        line-height: 1.6;
        color: var(--muted-foreground);
        margin-bottom: 2.5rem;
    }

    .hero-actions {
        display: flex;
        gap: 1rem;
        flex-wrap: wrap;
        justify-content: center;
    }

    .btn {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1.5rem;
        border-radius: var(--radius);
        font-weight: 600;
        text-decoration: none;
        transition: all 0.2s ease;
        cursor: pointer;
    }

    .btn-primary {
        background: var(--primary);
        color: var(--background);
    }

    .btn-primary:hover {
        filter: brightness(1.1);
        transform: translateY(-1px);
    }

    .btn-outline {
        border: 1px solid var(--border);
        background: var(--card);
        color: var(--foreground);
    }

    .btn-outline:hover {
        background: var(--muted);
    }

    /* Pillars */
    .section-header {
        text-align: center;
        margin-bottom: 3rem;
    }

    .eyebrow {
        font-size: 0.8rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--primary);
        margin-bottom: 0.5rem;
    }

    .section-header h2 {
        font-family: var(--font-serif);
        font-size: 2.25rem;
        font-weight: 700;
    }

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
        transition: transform 0.2s ease, border-color 0.2s ease;
    }

    .pillar-card:hover {
        transform: translateY(-4px);
        border-color: var(--primary);
    }

    .card-icon {
        margin-bottom: 1.25rem;
    }

    .pillar-card h3 {
        font-family: var(--font-serif);
        font-size: 1.35rem;
        font-weight: 600;
        margin-bottom: 0.75rem;
    }

    .pillar-card p {
        font-size: 0.95rem;
        line-height: 1.5;
        color: var(--muted-foreground);
    }

    /* Manifesto */
    .manifesto-box {
        background: oklch(from var(--primary) l c h / 0.08);
        border: 1px solid oklch(from var(--primary) l c h / 0.25);
        border-radius: var(--radius);
        padding: 3rem;
        text-align: center;
    }

    .manifesto-header {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.75rem;
        margin-bottom: 1.5rem;
    }

    .manifesto-header h2 {
        font-family: var(--font-serif);
        font-size: 1.85rem;
    }

    .manifesto-text {
        font-family: var(--font-serif);
        font-size: 1.35rem;
        font-style: italic;
        line-height: 1.6;
        color: var(--foreground);
        max-width: 900px;
        margin: 0 auto;
    }

    /* Articles */
    .articles-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(340px, 1fr));
        gap: 2rem;
    }

    .article-card {
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 1.75rem;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }

    .card-meta {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
        font-size: 0.8rem;
    }

    .category-badge {
        background: var(--muted);
        color: var(--primary);
        padding: 0.2rem 0.6rem;
        border-radius: 4px;
        font-weight: 600;
    }

    .date {
        color: var(--muted-foreground);
    }

    .article-card h3 {
        font-family: var(--font-serif);
        font-size: 1.25rem;
        font-weight: 600;
        margin-bottom: 0.75rem;
    }

    .article-card p {
        font-size: 0.9rem;
        line-height: 1.5;
        color: var(--muted-foreground);
        margin-bottom: 1.5rem;
        flex-grow: 1;
    }

    .card-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        border-top: 1px solid var(--border);
        padding-top: 0.75rem;
        font-size: 0.85rem;
    }

    .p2p-tag {
        font-size: 0.75rem;
        color: var(--primary);
        font-weight: 600;
    }

    .loading-state, .empty-state {
        text-align: center;
        padding: 3rem;
        color: var(--muted-foreground);
    }

    /* Footer */
    .public-footer {
        border-top: 1px solid var(--border);
        padding: 2rem 0;
        margin-top: 4rem;
        background: var(--card);
        font-size: 0.9rem;
        color: var(--muted-foreground);
    }

    .footer-content {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-top: 0;
        padding-bottom: 0;
    }

    .footer-links {
        display: flex;
        gap: 1.5rem;
        align-items: center;
    }

    .footer-links a {
        color: var(--foreground);
        text-decoration: none;
    }

    .footer-links a:hover {
        color: var(--primary);
    }

    .version {
        font-size: 0.8rem;
        background: var(--muted);
        padding: 0.2rem 0.5rem;
        border-radius: 4px;
    }
</style>
