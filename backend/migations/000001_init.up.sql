//User part
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)

CREATE TABLE refresh_tokens (
    id UUID NOT NULL UNIQUE PRIMARY KEY,
    user_id BIGINT references users(id) NOT NULL,
    token_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT false
)

CREATE TABLE user_contacts (
    contact_id BIGSERIAL PRIMARY KEY;
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Constraints
    CONSTRAINT user_contact_pair_uniqueness
        UNIQUE (user_id, contact_user_id),
    CONSTRAINT user_contact_ids_not_equal
        CHECK (user_id <> contact_user_id)
)

-- Chat part

CREATE TYPE conversation_type as ENUM('dm', 'group')

CREATE TABLE conversations (
    id BIGSERIAL PRIMARY KEY,
    guid UUID NOT NULL UNIQUE,
    title conversation_type,
    type TEXT NOT NULL CHECK(type IN ('dm', 'group')),
    dm_user1_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    dm_user2_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    creator_id BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT dm_unique_pair_1
        UNIQUE (dm_user1_id, dm_user2_id),
    CONSTRAINT dm_unique_pair_2
        CHECK (type <> 'dm' OR (dm_user1_id < dm_user2_id))
    CONSTRAINT dm_ids_is_null_for_group
        CHECK (type <> 'dm' OR (dm_user1_id IS NOT NULL AND dm_user2_id IS NOT NULL)),
    CONSTRAINT dm_ids_not_null_for_dm
        CHECK (type <> 'group' OR (dm_user1_id IS NULL AND dm_user2_id IS NULL))
)

CREATE TABLE participants (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Constraints
    CONSTRAINT user_conversation_pair_uniqueness
        UNIQUE (conversation_id, user_id)
)

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id),
    sender_id BIGINT NOT NULL REFERENCES users(id),
    message_body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)