CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	discord_id TEXT UNIQUE,
	timestamp DATETIME DEFAULT (datetime('now')),
	author TEXT NOT NULL,
	author_id TEXT NOT NULL DEFAULT '',
	content TEXT NOT NULL,
	channel_id TEXT NOT NULL,
	category TEXT DEFAULT 'uncategorized',
	mentions_json TEXT NOT NULL DEFAULT '[]',
	attachments_json TEXT NOT NULL DEFAULT '[]',
	reactions_json TEXT NOT NULL DEFAULT '[]',
	reply_to_id TEXT DEFAULT '',
	episode_id TEXT NOT NULL DEFAULT '',
	is_bot BOOLEAN NOT NULL DEFAULT 0,
	is_deleted BOOLEAN NOT NULL DEFAULT 0,
	internal_monologue TEXT NOT NULL DEFAULT '',
	grounding_ledger TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_messages_channel_timestamp ON messages(channel_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_messages_channel_episode_row ON messages(channel_id, episode_id, id);

CREATE TABLE IF NOT EXISTS memory_capsules (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	day_date TEXT NOT NULL,
	time_span TEXT NOT NULL,
	episode_id TEXT NOT NULL DEFAULT '',
	kind TEXT NOT NULL DEFAULT 'episode_snapshot',
	source_start_row_id INTEGER NOT NULL DEFAULT 0,
	source_end_row_id INTEGER NOT NULL DEFAULT 0,
	source_message_count INTEGER NOT NULL DEFAULT 0,
	participants JSON NOT NULL,
	main_topic TEXT NOT NULL,
	mood TEXT NOT NULL,
	key_facts JSON NOT NULL,
	typed_facts_json TEXT NOT NULL DEFAULT '{}',
	unresolved_questions JSON NOT NULL,
	open_loops_json TEXT NOT NULL DEFAULT '[]',
	category TEXT NOT NULL DEFAULT 'general',
	is_merged BOOLEAN DEFAULT FALSE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_memory_capsules_date ON memory_capsules(day_date);

CREATE TABLE IF NOT EXISTS budget_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	turn_id TEXT,
	cost_tokens INTEGER NOT NULL,
	input_tokens INTEGER NOT NULL DEFAULT 0,
	output_tokens INTEGER NOT NULL DEFAULT 0,
	model TEXT NOT NULL DEFAULT '',
	reason TEXT NOT NULL,
	trigger_type TEXT NOT NULL DEFAULT '',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gateway_state (
	channel_id TEXT PRIMARY KEY,
	last_pulse_row_id INTEGER NOT NULL DEFAULT 0,
	last_pulse_msg_id TEXT NOT NULL DEFAULT '',
	last_pulse_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE VIRTUAL TABLE IF NOT EXISTS memory_search USING fts5(
	participants,
	main_topic,
	key_facts,
	category,
	content='',
	contentless_delete=1,
	tokenize='unicode61 remove_diacritics 2'
);

CREATE TABLE IF NOT EXISTS message_artifacts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	message_id TEXT NOT NULL,
	type TEXT NOT NULL,
	content TEXT NOT NULL,
	metadata_json TEXT NOT NULL DEFAULT '{}',
	created_at DATETIME DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_message_artifacts_msg ON message_artifacts(message_id);

CREATE TABLE IF NOT EXISTS stance_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	channel_id TEXT NOT NULL,
	author_id TEXT NOT NULL,
	topic TEXT NOT NULL,
	position TEXT NOT NULL,
	action TEXT NOT NULL,
	evidence_type TEXT NOT NULL,
	confidence REAL NOT NULL,
	source_msg_id TEXT NOT NULL,
	episode_id TEXT NOT NULL DEFAULT '',
	validated BOOLEAN NOT NULL DEFAULT 0,
	created_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS behavior_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	channel_id TEXT NOT NULL,
	author_id TEXT NOT NULL,
	episode_id TEXT NOT NULL DEFAULT '',
	archetype TEXT NOT NULL,
	rhetorical_style TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'pending',
	created_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS agent_profiles (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	system_prompt TEXT NOT NULL,
	is_active BOOLEAN DEFAULT 0
);

CREATE TABLE IF NOT EXISTS core_identity_profiles (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	display_name TEXT NOT NULL DEFAULT '',
	avatar_url TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',
	identity_prompt TEXT NOT NULL,
	core_values_json TEXT NOT NULL DEFAULT '[]',
	is_enabled BOOLEAN DEFAULT 1,
	is_default BOOLEAN DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS persona_overlays (
	id TEXT PRIMARY KEY,
	identity_id TEXT NOT NULL,
	name TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	style_prompt TEXT NOT NULL,
	traits_json TEXT NOT NULL DEFAULT '{}',
	allowed_intents_json TEXT NOT NULL DEFAULT '[]',
	is_enabled BOOLEAN DEFAULT 1,
	is_default BOOLEAN DEFAULT 0,
	sort_order INTEGER DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(identity_id) REFERENCES core_identity_profiles(id)
);

CREATE TABLE IF NOT EXISTS adaptive_persona_memory (
	channel_id TEXT NOT NULL,
	identity_id TEXT NOT NULL,
	persona_id TEXT NOT NULL,
	adaptive_style_json TEXT NOT NULL DEFAULT '{}',
	confidence REAL NOT NULL DEFAULT 0,
	source TEXT NOT NULL DEFAULT 'manual',
	updated_at DATETIME NOT NULL,
	expires_at DATETIME,
	PRIMARY KEY(channel_id, identity_id, persona_id),
	FOREIGN KEY(identity_id) REFERENCES core_identity_profiles(id),
	FOREIGN KEY(persona_id) REFERENCES persona_overlays(id)
);

CREATE TABLE IF NOT EXISTS persona_delta_proposals (
	id TEXT PRIMARY KEY,
	channel_id TEXT NOT NULL,
	identity_id TEXT NOT NULL,
	persona_id TEXT NOT NULL DEFAULT '',
	target TEXT NOT NULL,
	proposed_changes_json TEXT NOT NULL DEFAULT '{}',
	reason TEXT NOT NULL DEFAULT '',
	evidence_message_ids_json TEXT NOT NULL DEFAULT '[]',
	confidence REAL NOT NULL DEFAULT 0,
	status TEXT NOT NULL DEFAULT 'pending',
	snapshot_id TEXT NOT NULL DEFAULT '',
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	reviewed_at DATETIME,
	expires_at DATETIME
);

CREATE TABLE IF NOT EXISTS message_annotations (
	message_id TEXT PRIMARY KEY,
	channel_id TEXT NOT NULL,
	episode_id TEXT NOT NULL DEFAULT '',
	author_id TEXT NOT NULL,
	topic_tags_json TEXT NOT NULL DEFAULT '[]',
	stance_tags_json TEXT NOT NULL DEFAULT '[]',
	style_tags_json TEXT NOT NULL DEFAULT '[]',
	evidence_type TEXT NOT NULL DEFAULT '',
	durability_score REAL NOT NULL DEFAULT 0,
	retrieval_score REAL NOT NULL DEFAULT 0,
	humor_score REAL NOT NULL DEFAULT 0,
	sarcasm_score REAL NOT NULL DEFAULT 0,
	snapshot_id TEXT NOT NULL DEFAULT '',
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS persona_policy (
	channel_id TEXT PRIMARY KEY,
	default_identity_id TEXT NOT NULL,
	default_persona_id TEXT NOT NULL,
	selection_mode TEXT NOT NULL DEFAULT 'fixed',
	allowed_identity_ids_json TEXT NOT NULL DEFAULT '[]',
	allowed_persona_ids_json TEXT NOT NULL DEFAULT '[]',
	intent_persona_map_json TEXT NOT NULL DEFAULT '{}',
	mode_persona_map_json TEXT NOT NULL DEFAULT '{}',
	manual_override_identity_id TEXT NOT NULL DEFAULT '',
	manual_override_persona_id TEXT NOT NULL DEFAULT '',
	updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS sync_state (
	channel_id TEXT PRIMARY KEY,
	latest_synced_msg_id TEXT NOT NULL DEFAULT '',
	oldest_synced_msg_id TEXT NOT NULL DEFAULT '',
	updated_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS system_metrics (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	metric_type TEXT NOT NULL,
	sub_type TEXT NOT NULL DEFAULT '',
	value REAL NOT NULL DEFAULT 1,
	latency_ms INTEGER NOT NULL DEFAULT 0,
	metadata_json TEXT NOT NULL DEFAULT '{}',
	timestamp DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS system_config (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL,
	updated_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS capsule_hwm (
	channel_id TEXT PRIMARY KEY,
	last_message_id INTEGER NOT NULL DEFAULT 0,
	updated_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS pending_turns (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	channel_id TEXT NOT NULL,
	message_id TEXT NOT NULL,
	message_row_id INTEGER NOT NULL,
	timestamp DATETIME NOT NULL,
	is_reactive BOOLEAN NOT NULL,
	retry_count INTEGER NOT NULL DEFAULT 0,
	last_error TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'pending',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(channel_id, message_id)
);
CREATE INDEX IF NOT EXISTS idx_pending_turns_status ON pending_turns(status);

CREATE TABLE IF NOT EXISTS member_profiles (
	discord_id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	roles_json TEXT NOT NULL DEFAULT '[]',
	age INTEGER DEFAULT 0,
	interests TEXT NOT NULL DEFAULT '',
	religion TEXT NOT NULL DEFAULT '',
	notes TEXT NOT NULL DEFAULT '',
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS guild_config (
	guild_id TEXT PRIMARY KEY,
	prefix TEXT NOT NULL DEFAULT '/',
	default_channel_id TEXT NOT NULL DEFAULT '',
	mod_log_channel_id TEXT NOT NULL DEFAULT '',
	gate_model TEXT NOT NULL DEFAULT '',
	reply_model TEXT NOT NULL DEFAULT '',
	memory_model TEXT NOT NULL DEFAULT '',
	mod_enabled BOOLEAN NOT NULL DEFAULT 0,
	max_warnings INTEGER NOT NULL DEFAULT 3,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS guild_warnings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	channel_id TEXT NOT NULL,
	reason TEXT NOT NULL,
	message_id TEXT NOT NULL DEFAULT '',
	severity INTEGER NOT NULL DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_guild_warnings_user ON guild_warnings(guild_id, user_id);

CREATE TABLE IF NOT EXISTS forum_templates (
	id TEXT PRIMARY KEY,
	guild_id TEXT NOT NULL,
	channel_id TEXT NOT NULL,
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	tags_json TEXT NOT NULL DEFAULT '[]',
	variables_json TEXT NOT NULL DEFAULT '[]',
	schedule TEXT NOT NULL DEFAULT 'manual',
	schedule_config_json TEXT NOT NULL DEFAULT '{}',
	is_enabled BOOLEAN NOT NULL DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS forum_posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	template_id TEXT NOT NULL,
	guild_id TEXT NOT NULL,
	channel_id TEXT NOT NULL,
	discord_message_id TEXT NOT NULL DEFAULT '',
	discord_thread_id TEXT NOT NULL DEFAULT '',
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	tags_json TEXT NOT NULL DEFAULT '[]',
	status TEXT NOT NULL DEFAULT 'draft',
	error TEXT NOT NULL DEFAULT '',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(template_id) REFERENCES forum_templates(id)
);

CREATE TABLE IF NOT EXISTS member_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	user_name TEXT NOT NULL DEFAULT '',
	event_type TEXT NOT NULL,
	reason TEXT NOT NULL DEFAULT '',
	executor_id TEXT NOT NULL DEFAULT '',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_member_events_guild ON member_events(guild_id, created_at);
CREATE INDEX IF NOT EXISTS idx_member_events_user ON member_events(guild_id, user_id);

CREATE TABLE IF NOT EXISTS audit_log (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	guild_id TEXT NOT NULL,
	action_type TEXT NOT NULL,
	action TEXT NOT NULL,
	actor_id TEXT NOT NULL DEFAULT '',
	actor_name TEXT NOT NULL DEFAULT '',
	target_id TEXT NOT NULL DEFAULT '',
	target_name TEXT NOT NULL DEFAULT '',
	details TEXT NOT NULL DEFAULT '',
	channel_id TEXT NOT NULL DEFAULT '',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_audit_log_guild ON audit_log(guild_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_log_action ON audit_log(guild_id, action_type);

CREATE TABLE IF NOT EXISTS welcome_config (
	guild_id TEXT PRIMARY KEY,
	enabled BOOLEAN NOT NULL DEFAULT 0,
	channel_id TEXT NOT NULL DEFAULT '',
	welcome_message TEXT NOT NULL DEFAULT '',
	goodbye_message TEXT NOT NULL DEFAULT '',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS role_policies (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	guild_id TEXT NOT NULL,
	role_id TEXT NOT NULL,
	role_name TEXT NOT NULL DEFAULT '',
	trigger_type TEXT NOT NULL,
	trigger_value TEXT NOT NULL DEFAULT '',
	enabled BOOLEAN NOT NULL DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_role_policies_guild ON role_policies(guild_id);

-- FTS5 triggers for memory_search auto-sync
DROP TRIGGER IF EXISTS memory_capsules_ai;
CREATE TRIGGER memory_capsules_ai AFTER INSERT ON memory_capsules BEGIN
	INSERT OR REPLACE INTO memory_search(rowid, participants, main_topic, key_facts, category)
	VALUES (
		NEW.id,
		COALESCE((SELECT group_concat(value, ' ') FROM json_each(NEW.participants)), ''),
		COALESCE(NEW.main_topic, ''),
		TRIM(
			COALESCE((SELECT group_concat(value, ' ') FROM json_each(NEW.key_facts)), '') || ' ' ||
			COALESCE((SELECT group_concat(value, ' ') FROM json_each(NEW.unresolved_questions)), '')
		),
		COALESCE(NEW.category, 'general')
	);
END;

DROP TRIGGER IF EXISTS memory_capsules_au;
CREATE TRIGGER memory_capsules_au AFTER UPDATE ON memory_capsules BEGIN
	DELETE FROM memory_search WHERE rowid = OLD.id;
	INSERT INTO memory_search(rowid, participants, main_topic, key_facts, category)
	VALUES (
		NEW.id,
		COALESCE((SELECT group_concat(value, ' ') FROM json_each(NEW.participants)), ''),
		COALESCE(NEW.main_topic, ''),
		TRIM(
			COALESCE((SELECT group_concat(value, ' ') FROM json_each(NEW.key_facts)), '') || ' ' ||
			COALESCE((SELECT group_concat(value, ' ') FROM json_each(NEW.unresolved_questions)), '')
		),
		COALESCE(NEW.category, 'general')
	);
END;

DROP TRIGGER IF EXISTS memory_capsules_ad;
CREATE TRIGGER memory_capsules_ad AFTER DELETE ON memory_capsules BEGIN
	DELETE FROM memory_search WHERE rowid = OLD.id;
END;

-- Rebuild FTS index from existing data
DELETE FROM memory_search;
INSERT OR REPLACE INTO memory_search(rowid, participants, main_topic, key_facts, category)
SELECT
	id,
	COALESCE((SELECT group_concat(value, ' ') FROM json_each(participants)), ''),
	COALESCE(main_topic, ''),
	TRIM(
		COALESCE((SELECT group_concat(value, ' ') FROM json_each(key_facts)), '') || ' ' ||
		COALESCE((SELECT group_concat(value, ' ') FROM json_each(unresolved_questions)), '')
	),
	COALESCE(category, 'general')
FROM memory_capsules;
