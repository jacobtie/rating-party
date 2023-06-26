-- Not using full migrations yet

USE ratingparty;

CREATE TABLE Game (
    game_id BINARY DEFAULT (UUID_TO_BIN(UUID())),
    game_name VARCHAR(255) NOT NULL,
    game_code VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (game_id)
);
