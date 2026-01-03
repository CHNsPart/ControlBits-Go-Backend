CREATE TABLE habit_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    habit_id UUID NOT NULL,
    entry_date DATE NOT NULL,
    status VARCHAR(10) NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_entries_habit
        FOREIGN KEY (habit_id)
        REFERENCES habits(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_habit_date
        UNIQUE (habit_id, entry_date)
);
