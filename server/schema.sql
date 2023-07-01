-- Not using full migrations yet

USE ratingparty;

CREATE TABLE game (
    game_id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())),
    game_name VARCHAR(255) NOT NULL,
    game_code VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (game_id)
);

CREATE TABLE wine (
    wine_id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())),
    wine_name VARCHAR(255) NOT NULL,
    wine_code VARCHAR(255) NOT NULL,
    wine_year INT NOT NULL,
    game_id BINARY(16) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (wine_id),
    FOREIGN KEY (game_id) REFERENCES game(game_id)
);
