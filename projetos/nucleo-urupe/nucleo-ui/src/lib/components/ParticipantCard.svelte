<script>
    import { User, Activity, Heart, Zap, Quote } from 'lucide-svelte';
    let { participant } = $props();

    function getInitials(name) {
        return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
    }
</script>

<div class="participant-card panel">
    <div class="card-header">
        <div class="avatar-glow">
            <div class="avatar">{getInitials(participant.AuthorName)}</div>
        </div>
        <div class="header-info">
            <h3>{participant.AuthorName}</h3>
            <span class="updated">Ativo em {participant.UpdatedAt}</span>
        </div>
    </div>

    <div class="card-content">
        <div class="profile-section">
            <div class="section-title">
                <Activity size={14} />
                <span>Estilo</span>
            </div>
            <p class="style-text">{participant.Style || 'Aguardando mapeamento...'}</p>
        </div>

        <div class="profile-section">
            <div class="section-title">
                <Heart size={14} />
                <span>Interesses</span>
            </div>
            <div class="interests-tags">
                {#if participant.Interests}
                    {#each participant.Interests.split(',') as interest}
                        <span class="tag">{interest.trim()}</span>
                    {/each}
                {:else}
                    <span class="tag empty">sem interesses mapeados</span>
                {/if}
            </div>
        </div>

        <div class="profile-section">
            <div class="section-title">
                <Zap size={14} />
                <span>Tendencias</span>
            </div>
            <p class="tendencies-text">{participant.Tendencies || 'Aguardando observação...'}</p>
        </div>

        {#if participant.SocialNote}
            <div class="social-note">
                <Quote size={12} class="quote-icon" />
                <p>{participant.SocialNote}</p>
            </div>
        {/if}
    </div>
</div>

<style>
    .participant-card {
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
        background: var(--card);
        position: relative;
        overflow: hidden;
    }

    .participant-card::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 2px;
        background: linear-gradient(90deg, transparent, var(--primary), transparent);
        opacity: 0.3;
    }

    .card-header {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .avatar-glow {
        position: relative;
        padding: 2px;
        border-radius: 50%;
        background: linear-gradient(135deg, var(--primary), var(--primary));
    }

    .avatar {
        width: 48px;
        height: 48px;
        background: var(--background);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-family: var(--font-mono);
        font-weight: 700;
        color: var(--primary);
        font-size: 1.1rem;
    }

    .header-info h3 {
        font-size: 1.2rem;
        margin-bottom: 0.2rem;
    }

    .updated {
        font-size: 0.75rem;
        color: var(--muted-foreground);
        font-family: var(--font-mono);
    }

    .card-content {
        display: flex;
        flex-direction: column;
        gap: 1.2rem;
    }

    .section-title {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: var(--muted-foreground);
        text-transform: uppercase;
        font-size: 0.7rem;
        font-weight: 700;
        letter-spacing: 0.05em;
        margin-bottom: 0.5rem;
    }

    .style-text, .tendencies-text {
        font-size: 0.9rem;
        line-height: 1.5;
        color: var(--muted-foreground);
    }

    .interests-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
    }

    .tag {
        font-size: 0.75rem;
        padding: 0.2rem 0.6rem;
        background: var(--muted);
        color: var(--primary);
        border-radius: 6px;
        border: 1px solid var(--border);
    }

    .tag.empty {
        background: transparent;
        color: var(--muted-foreground);
        border-style: dashed;
    }

    .social-note {
        margin-top: 0.5rem;
        padding: 1rem;
        background: rgba(255, 255, 255, 0.03);
        border-radius: 12px;
        border-left: 3px solid var(--primary);
        position: relative;
    }

    .social-note p {
        font-size: 0.85rem;
        font-style: italic;
        color: var(--muted-foreground);
        line-height: 1.4;
    }

    :global(.quote-icon) {
        position: absolute;
        top: -6px;
        left: 6px;
        color: var(--primary);
        background: var(--card);
        padding: 2px;
    }

    .participant-card:hover {
        transform: translateY(-4px);
        box-shadow: 0 12px 30px rgba(0, 0, 0, 0.4);
        border-color: var(--border);
    }
</style>
