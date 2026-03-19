-- 1. Users table (Auth + game stats)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    image TEXT,
    games_played INT NOT NULL DEFAULT 0,
    games_won INT NOT NULL DEFAULT 0,
    average_wpm NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    highest_wpm INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- 2. Texts table (typing prompts)
CREATE TABLE IF NOT EXISTS texts (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    difficulty VARCHAR(50) NOT NULL, -- 'easy', 'medium', 'hard'
    word_count INT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- 3. Matches table (race sessions)
CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_code VARCHAR(20) UNIQUE NOT NULL, -- shareable link code
    text_id INT NOT NULL REFERENCES texts(id),
    status VARCHAR(50) NOT NULL DEFAULT 'completed', -- 'completed', 'aborted'
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- 4. Match Results table (performance per player per match)
CREATE TABLE IF NOT EXISTS match_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wpm INT NOT NULL,
    accuracy NUMERIC(5,2) NOT NULL,
    is_winner BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (match_id, user_id)
);

-- Helpful index for user history
CREATE INDEX IF NOT EXISTS idx_match_results_user_id ON match_results(user_id);
