ALTER TABLE repo_analyzes
  ADD COLUMN attempt_number INTEGER NOT NULL DEFAULT 1;