CREATE TABLE events (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
    "name" TEXT,
    "date" TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "description" TEXT
);


CREATE TRIGGER update_event_updated_at_timestamp
AFTER UPDATE ON events
FOR EACH ROW
BEGIN
    UPDATE events SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
