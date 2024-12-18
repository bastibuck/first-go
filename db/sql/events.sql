CREATE TABLE events (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
    "name" TEXT,
    "date" TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "description" TEXT
);