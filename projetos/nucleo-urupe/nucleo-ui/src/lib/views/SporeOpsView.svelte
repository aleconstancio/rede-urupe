<script>
    import { onMount } from 'svelte';
    import { Send, Image, Radio, Sparkles, Share2, Download, Eye, TrendingUp, AlertCircle, FileText, CheckCircle2 } from 'lucide-svelte';

    let activeTab = $state('agitprop'); // 'agitprop' | 'clipping' | 'analytics'

    const agitpropPackages = [
        {
            id: 'pkg-1',
            title: 'Pacote #01: Soberania Digital é Poder Popular',
            category: 'Redes Sociais & Cartazes',
            date: '03/Set/2026',
            downloads: 342,
            itemsCount: 6,
            desc: 'Carrossel em alta resolução para Instagram/Mastodon e arte A3 para impressão de lambes nas ruas.'
        },
        {
            id: 'pkg-2',
            title: 'Pacote #02: Contra o Extrativismo de Dados',
            category: 'Stickers & Panfletos',
            date: '01/Set/2026',
            downloads: 512,
            itemsCount: 12,
            desc: 'Figurinhas para WhatsApp/Signal/Telegram desmistificando o modelo de vigilância das Big Techs.'
        },
        {
            id: 'pkg-3',
            title: 'Pacote #03: Ecossocialismo ou Colapso',
            category: 'Fanzines & Áudio',
            date: '28/Ago/2026',
            downloads: 289,
            itemsCount: 4,
            desc: 'Zine em PDF pronto para dobra A4 e spot de rádio comunitária em MP3.'
        }
    ];

    const clippingItems = [
        {
            source: 'Jornal Popular',
            title: 'Comunidades adotam redes P2P para escapar do controle algorítmico',
            sentiment: 'Positivo',
            reach: '14.2k leitores',
            date: 'Hoje, 07:45'
        },
        {
            source: 'Blog Soberania Livre',
            title: 'Análise: Como o protocolo Rizoma desafia o monopólio das plataformas',
            sentiment: 'Positivo',
            reach: '8.7k leitores',
            date: 'Ontem, 18:20'
        },
        {
            source: 'Agência de Notícias',
            title: 'Frente Urupê lança manifesto por tecnologia livre e defesa ambiental',
            sentiment: 'Neutro',
            reach: '22.1k leitores',
            date: '01/Set/2026'
        }
    ];
</script>

<div class="spore-ops-view">
    <div class="view-header">
        <div>
            <p class="eyebrow text-pink">Operações Meméticas & Clipping</p>
            <h2 class="font-serif text-2xl font-bold flex items-center gap-2">
                <span>Spore Ops 🍄</span>
                <span class="badge-tag">Agitprop & Imprensa</span>
            </h2>
        </div>

        <div class="tab-switcher">
            <button class="tab-btn {activeTab === 'agitprop' ? 'active' : ''}" onclick={() => activeTab = 'agitprop'}>
                <Send size={16} />
                <span>Agitprop & Materiais</span>
            </button>
            <button class="tab-btn {activeTab === 'clipping' ? 'active' : ''}" onclick={() => activeTab = 'clipping'}>
                <Radio size={16} />
                <span>Clipping de Imprensa</span>
            </button>
        </div>
    </div>

    <!-- TAB 1: Agitprop -->
    {#if activeTab === 'agitprop'}
        <div class="ops-grid">
            <div class="panel main-panel">
                <div class="panel-header">
                    <h3>Pacotes de Agitação & Propaganda Popular</h3>
                    <span class="active-tag">3 Pacotes Prontos</span>
                </div>

                <div class="packages-list">
                    {#each agitpropPackages as pkg}
                        <div class="package-card">
                            <div class="pkg-header">
                                <span class="pkg-cat">{pkg.category}</span>
                                <span class="pkg-date">{pkg.date}</span>
                            </div>
                            <h4 class="pkg-title">{pkg.title}</h4>
                            <p class="pkg-desc">{pkg.desc}</p>

                            <div class="pkg-footer">
                                <span class="pkg-stat"><Download size={14} /> {pkg.downloads} downloads</span>
                                <span class="pkg-stat"><Image size={14} /> {pkg.itemsCount} arquivos</span>
                                <button class="btn btn-primary btn-sm ml-auto">
                                    <Download size={14} /> Baixar Kit
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>

            <!-- Quick Action Box -->
            <div class="panel side-panel">
                <div class="panel-header">
                    <h3>Criar Nova Campanha</h3>
                </div>
                <div class="action-form">
                    <div class="form-group">
                        <label for="spore-camp-title">Nome da Campanha</label>
                        <input id="spore-camp-title" type="text" placeholder="Ex: Jornada de Formação Soberana" class="input" />
                    </div>
                    <div class="form-group">
                        <label for="spore-camp-type">Tipo de Material</label>
                        <select id="spore-camp-type" class="input">
                            <option>Carrossel de Redes Sociais</option>
                            <option>Stickers & Figurinhas</option>
                            <option>Cartazes & Lambes A3</option>
                            <option>Spot de Áudio para Rádio</option>
                        </select>
                    </div>
                    <button class="btn btn-pink w-full">
                        <Sparkles size={16} /> Gerar Pacote com a Micélia 🍄
                    </button>
                </div>
            </div>
        </div>
    {/if}

    <!-- TAB 2: Clipping -->
    {#if activeTab === 'clipping'}
        <div class="panel">
            <div class="panel-header">
                <h3>Clipping Automatizado de Imprensa & Redes</h3>
            </div>
            <div class="clipping-list">
                {#each clippingItems as item}
                    <div class="clip-card">
                        <div class="clip-header">
                            <span class="clip-source">{item.source}</span>
                            <span class="clip-date">{item.date}</span>
                        </div>
                        <h4 class="clip-title">{item.title}</h4>
                        <div class="clip-footer">
                            <span class="clip-badge positive">{item.sentiment}</span>
                            <span class="clip-reach"><Eye size={14} /> Alcance: {item.reach}</span>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</div>

<style>
    .spore-ops-view { padding: 2rem; max-width: 1600px; margin: 0 auto; }
    .view-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
    .eyebrow { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; }
    .text-pink { color: #ec4899; }
    .badge-tag { font-size: 0.75rem; background: rgba(236, 72, 153, 0.15); color: #ec4899; padding: 0.2rem 0.6rem; border-radius: 9999px; }

    .tab-switcher { display: flex; gap: 0.5rem; background: var(--card); border: 1px solid var(--border); padding: 0.3rem; border-radius: var(--radius); }
    .tab-btn { display: flex; align-items: center; gap: 0.4rem; padding: 0.5rem 1rem; border-radius: calc(var(--radius) - 2px); font-size: 0.85rem; font-weight: 600; background: transparent; border: none; color: var(--muted-foreground); cursor: pointer; }
    .tab-btn.active { background: #ec4899; color: white; }

    .ops-grid { display: grid; grid-template-columns: 1fr 360px; gap: 2rem; }
    .panel { background: var(--card); border: 1px solid var(--border); border-radius: var(--radius); overflow: hidden; }
    .panel-header { padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
    .panel-header h3 { font-weight: 600; font-size: 1rem; }

    .packages-list { padding: 1.5rem; display: flex; flex-direction: column; gap: 1.25rem; }
    .package-card { background: var(--background); border: 1px solid var(--border); border-radius: var(--radius); padding: 1.5rem; }
    .pkg-header { display: flex; justify-content: space-between; font-size: 0.8rem; margin-bottom: 0.5rem; }
    .pkg-cat { font-weight: 700; color: #ec4899; }
    .pkg-date { color: var(--muted-foreground); }
    .pkg-title { font-family: var(--font-serif); font-size: 1.2rem; font-weight: 700; margin-bottom: 0.5rem; }
    .pkg-desc { font-size: 0.9rem; color: var(--muted-foreground); line-height: 1.5; margin-bottom: 1.25rem; }
    .pkg-footer { display: flex; align-items: center; gap: 1.25rem; font-size: 0.8rem; color: var(--muted-foreground); }
    .pkg-stat { display: flex; align-items: center; gap: 0.4rem; }

    .action-form { padding: 1.5rem; display: flex; flex-direction: column; gap: 1.25rem; }
    .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
    .input { background: var(--background); border: 1px solid var(--border); border-radius: var(--radius); padding: 0.6rem 0.8rem; color: var(--foreground); font-size: 0.9rem; }
    .btn { display: inline-flex; align-items: center; justify-content: center; gap: 0.5rem; padding: 0.6rem 1.2rem; border-radius: var(--radius); font-weight: 600; cursor: pointer; border: none; }
    .btn-pink { background: #ec4899; color: white; }
    .btn-primary { background: var(--primary); color: var(--background); }
    .btn-sm { padding: 0.4rem 0.8rem; font-size: 0.8rem; }
    .ml-auto { margin-left: auto; }
    .w-full { width: 100%; }

    .clipping-list { padding: 1.5rem; display: flex; flex-direction: column; gap: 1rem; }
    .clip-card { background: var(--background); border: 1px solid var(--border); border-radius: var(--radius); padding: 1.25rem; }
    .clip-header { display: flex; justify-content: space-between; font-size: 0.8rem; color: var(--muted-foreground); margin-bottom: 0.4rem; }
    .clip-source { font-weight: 700; color: var(--primary); }
    .clip-title { font-family: var(--font-serif); font-size: 1.1rem; font-weight: 600; margin-bottom: 0.75rem; }
    .clip-footer { display: flex; gap: 1rem; align-items: center; font-size: 0.8rem; }
    .clip-badge { padding: 0.15rem 0.5rem; border-radius: 4px; font-weight: 600; }
    .clip-badge.positive { background: rgba(16, 185, 129, 0.15); color: #10b981; }

    @media (max-width: 1000px) { .ops-grid { grid-template-columns: 1fr; } }
</style>
