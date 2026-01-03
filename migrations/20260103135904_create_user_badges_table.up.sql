CREATE TABLE user_badges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    habit_id UUID NOT NULL,
    badge_id UUID NOT NULL,
    earned_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_user_badges_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_badges_habit
        FOREIGN KEY (habit_id)
        REFERENCES habits(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_badges_badge
        FOREIGN KEY (badge_id)
        REFERENCES badges(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_user_habit_badge
        UNIQUE (user_id, habit_id, badge_id)
);
