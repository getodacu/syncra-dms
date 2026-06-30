ALTER TABLE "user" ADD COLUMN IF NOT EXISTS "preferred_language" varchar(5) NOT NULL DEFAULT 'en';

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'chk_user_preferred_language'
  ) THEN
    ALTER TABLE "user" ADD CONSTRAINT "chk_user_preferred_language" CHECK ("preferred_language" IN ('en','ro'));
  END IF;
END $$;
