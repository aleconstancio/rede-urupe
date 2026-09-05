package sqlite

import (
	_ "embed"
	"log"
)

//go:embed migrations/000001_initial_schema.up.sql
var initialSchemaSQL string

//go:embed migrations/000002_showcase_projects.up.sql
var showcaseProjectsSQL string

//go:embed migrations/000003_active_channels.up.sql
var activeChannelsSQL string

func (r *Repository) initSchema() error {
	if _, err := r.db.Exec(initialSchemaSQL); err != nil {
		log.Printf("[Migrations] Warning: failed to apply initialSchemaSQL: %v", err)
	}
	if _, err := r.db.Exec(showcaseProjectsSQL); err != nil {
		log.Printf("[Migrations] Warning: failed to apply showcaseProjectsSQL: %v", err)
	}
	if _, err := r.db.Exec(activeChannelsSQL); err != nil {
		log.Printf("[Migrations] Warning: failed to apply activeChannelsSQL: %v", err)
	}
	return r.seedInitialData()
}

func (r *Repository) seedInitialData() error {
	if err := r.ensureCMSTables(); err != nil {
		return err
	}
	if err := r.ensureManifestoTables(); err != nil {
		return err
	}
	if err := r.ensureAIreliusSeed(); err != nil {
		return err
	}
	if err := r.ensureV33PersonaSeed(); err != nil {
		return err
	}
	if err := r.ensureCMSArticlesSeed(); err != nil {
		return err
	}
	if err := r.ensureManifestoSeed(); err != nil {
		return err
	}
	return r.ensureShowcaseProjects()
}

func (r *Repository) ensureShowcaseProjects() error {
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM showcase_projects").Scan(&count)
	if count > 0 {
		return nil
	}

	projects := []struct {
		name, slug, desc, icon, url, github, category string
		featured                                      bool
		order                                         int
	}{
		{"Vico", "vico", "Automação jurídica com copilot IA, cálculos, prazos e petições para advogados brasileiros.", "⚖️", "", "https://github.com/aleconstancio/vico", "product", true, 1},
		{"Orb", "orb", "Prática mística local-first com journaling, astrologia, grimório e soberania de dados.", "🔮", "", "https://github.com/aleconstancio/orb", "product", true, 2},
		{"Núcleo Urupê", "nucleo-urupe", "Hub/QG central integrado ao Discord com a mascot IA Micélia 🍄 e governança.", "🍄", "", "https://github.com/aleconstancio/rede-urupe", "product", true, 3},
		{"Talos", "talos", "Motor de agentes IA com pipeline de 7 estágios, engenharia de contexto e aprendizado.", "🧠", "", "", "library", false, 4},
		{"Bindrunes", "bindrunes", "220+ componentes Svelte 5 com 6 temas, 4 estéticas e design system ortogonal.", "🎨", "", "https://github.com/aleconstancio/bindrunes", "library", false, 5},
		{"NuitOS", "nuitos", "Framework NixOS dendrítico com CLI Rust, TUI Dart e GUI Flutter.", "💻", "", "", "framework", false, 6},
	}

	for _, p := range projects {
		tagsJSON := "[]"
		_, err := r.db.Exec(`INSERT INTO showcase_projects (name, slug, description, icon, url, github_url, category, status, sort_order, tags_json, is_featured)
		         VALUES (?, ?, ?, ?, ?, ?, ?, 'active', ?, ?, ?)`,
			p.name, p.slug, p.desc, p.icon, p.url, p.github, p.category, p.order, tagsJSON, p.featured)
		if err != nil {
			log.Printf("[Migrations] Warning: failed to seed showcase project %s: %v", p.name, err)
		}
	}
	log.Printf("[Migrations] Seeded %d showcase projects", len(projects))
	return nil
}

func (r *Repository) ensureAIreliusSeed() error {
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM core_identity_profiles WHERE id = ?", "airelius_core").Scan(&count)
	if count > 0 {
		return nil
	}

	_, err := r.db.Exec(`
		INSERT OR IGNORE INTO core_identity_profiles (id, name, display_name, avatar_url, description, identity_prompt, core_values_json, is_enabled, is_default)
		VALUES (?, ?, ?, '', ?, ?, ?, 1, 0)
	`,
		"airelius_core", "AIrelius", "AIrelius",
		"Philosopher-bot. Socratic method, dialectical reasoning, explores first principles.",
		"", // IdentityPrompt loaded from Talos minotaur module at runtime
		`["intellectual_honesty","first_principles","dialectical_reasoning","stoic_equanimity","radical_questioning"]`,
	)
	if err != nil {
		log.Printf("[Migrations] Warning: failed to seed AIrelius identity: %v", err)
	}

	log.Printf("[Migrations] Seeded AIrelius identity")
	return nil
}
