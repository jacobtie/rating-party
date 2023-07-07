-- Not using full migrations yet

USE ratingparty;

CREATE TABLE game (
    game_id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())),
    game_name VARCHAR(255) NOT NULL,
    game_code VARCHAR(255) NOT NULL,
    is_running BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (game_id)
);

CREATE TABLE participant (
    participant_id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())),
    game_id BINARY(16) NOT NULL,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (participant_id),
    FOREIGN KEY (game_id) REFERENCES game(game_id),
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

CREATE TABLE rating (
    rating_id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())),
    game_id BINARY(16) NOT NULL,
    participant_id BINARY(16) NOT NULL,
    wine_id BINARY(16) NOT NULL,
    sight_rating FLOAT NOT NULL,
    aroma_rating FLOAT NOT NULL,
    taste_rating FLOAT NOT NULL,
    overall_rating FLOAT NOT NULL,
    comments TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (rating_id),
    FOREIGN KEY (game_id) REFERENCES game(game_id),
    FOREIGN KEY (participant_id) REFERENCES participant(participant_id),
    FOREIGN KEY (wine_id) REFERENCES wine(wine_id),
    UNIQUE (participant_id, wine_id)
);
