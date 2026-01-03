CREATE TABLE badges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(255) NOT NULL,
    streak_required INTEGER UNIQUE NOT NULL,
    icon VARCHAR(50) NOT NULL
);
