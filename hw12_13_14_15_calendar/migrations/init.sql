CREATE TABLE "events" (
		"id" TEXT NOT NULL,
		"title" TEXT NOT NULL,
		"date" BIGINT NOT NULL,
		"duration_until" BIGINT NOT NULL,
		"description" TEXT NOT NULL,
		"owner_id" TEXT NOT NULL,
		"notice_before" BIGINT NOT NULL,
		PRIMARY KEY ("id")
	);