-- User and access
-- -------------------------

CREATE TABLE IF NOT EXISTS lego (
    name    TEXT NOT NULL, 
    model   INTEGER PRIMARY KEY NOT NULL,
    catalog TEXT NOT NULL,       
);