-- Create pokedex summary table
CREATE TABLE IF NOT EXISTS player_pokedex_summary (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    total_caught INTEGER DEFAULT 0,
    total_seen INTEGER DEFAULT 0,
    regions_completed INTEGER DEFAULT 0,
    national_completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

-- Create regional pokedex tables
CREATE TABLE IF NOT EXISTS player_pokedex_kanto (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    caught_flags BYTEA DEFAULT '\x000000000000000000000000000000000000000',
    seen_flags BYTEA DEFAULT '\x000000000000000000000000000000000000000',
    completion_date TIMESTAMP WITH TIME ZONE,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

CREATE TABLE IF NOT EXISTS player_pokedex_johto (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    caught_flags BYTEA DEFAULT '\x0000000000000000000000000000',
    seen_flags BYTEA DEFAULT '\x0000000000000000000000000000',
    completion_date TIMESTAMP WITH TIME ZONE,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

CREATE TABLE IF NOT EXISTS player_pokedex_hoenn (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    caught_flags BYTEA DEFAULT '\x00000000000000000000000000000000000',
    seen_flags BYTEA DEFAULT '\x00000000000000000000000000000000000',
    completion_date TIMESTAMP WITH TIME ZONE,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

CREATE TABLE IF NOT EXISTS player_pokedex_sinnoh (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    caught_flags BYTEA DEFAULT '\x0000000000000000000000000000000000',
    seen_flags BYTEA DEFAULT '\x0000000000000000000000000000000000',
    completion_date TIMESTAMP WITH TIME ZONE,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

CREATE TABLE IF NOT EXISTS player_pokedex_unova (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    caught_flags BYTEA DEFAULT '\x000000000000000000000000000000000000000000',
    seen_flags BYTEA DEFAULT '\x000000000000000000000000000000000000000000',
    completion_date TIMESTAMP WITH TIME ZONE,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

CREATE TABLE IF NOT EXISTS player_pokedex_alola (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    caught_flags BYTEA DEFAULT '\x0000000000000000000000000000000000000',
    seen_flags BYTEA DEFAULT '\x0000000000000000000000000000000000000',
    completion_date TIMESTAMP WITH TIME ZONE,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

CREATE TABLE IF NOT EXISTS player_pokedex_galar (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    caught_flags BYTEA DEFAULT '\x00000000000000000000000000000000000000000',
    seen_flags BYTEA DEFAULT '\x00000000000000000000000000000000000000000',
    completion_date TIMESTAMP WITH TIME ZONE,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

CREATE TABLE IF NOT EXISTS player_pokedex_hisui (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    caught_flags BYTEA DEFAULT '\x000000000000000000000000000000000000000',
    seen_flags BYTEA DEFAULT '\x000000000000000000000000000000000000000',
    completion_date TIMESTAMP WITH TIME ZONE,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

CREATE TABLE IF NOT EXISTS player_pokedex_paldea (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    caught_flags BYTEA DEFAULT '\x0000000000000000000000000000000000000000000',
    seen_flags BYTEA DEFAULT '\x0000000000000000000000000000000000000000000',
    completion_date TIMESTAMP WITH TIME ZONE,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_pokedex_summary_player_id ON player_pokedex_summary(player_id);
CREATE INDEX IF NOT EXISTS idx_pokedex_summary_completion ON player_pokedex_summary(national_completion_percentage DESC);
CREATE INDEX IF NOT EXISTS idx_pokedex_kanto_player_id ON player_pokedex_kanto(player_id);
CREATE INDEX IF NOT EXISTS idx_pokedex_johto_player_id ON player_pokedex_johto(player_id);
CREATE INDEX IF NOT EXISTS idx_pokedex_hoenn_player_id ON player_pokedex_hoenn(player_id);
CREATE INDEX IF NOT EXISTS idx_pokedex_sinnoh_player_id ON player_pokedex_sinnoh(player_id);
CREATE INDEX IF NOT EXISTS idx_pokedex_unova_player_id ON player_pokedex_unova(player_id);
CREATE INDEX IF NOT EXISTS idx_pokedex_alola_player_id ON player_pokedex_alola(player_id);
CREATE INDEX IF NOT EXISTS idx_pokedex_galar_player_id ON player_pokedex_galar(player_id);
CREATE INDEX IF NOT EXISTS idx_pokedex_hisui_player_id ON player_pokedex_hisui(player_id);
CREATE INDEX IF NOT EXISTS idx_pokedex_paldea_player_id ON player_pokedex_paldea(player_id);