CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users
(
    nickname CITEXT PRIMARY KEY NOT NULL,
    email    CITEXT UNIQUE      NOT NULL,
    about    TEXT               NOT NULL,
    fullname TEXT               NOT NULL
);

CREATE UNIQUE INDEX ON users (nickname, email);
CREATE UNIQUE INDEX ON users (nickname, email, about, fullname);
CREATE UNIQUE INDEX ON users (nickname DESC);

CREATE UNLOGGED TABLE forums
(
    slug     CITEXT PRIMARY KEY                                   NOT NULL,
    title    TEXT                                                 NOT NULL,
    userf CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    posts    INTEGER DEFAULT 0                                    NOT NULL,
    threads  INTEGER DEFAULT 0                                    NOT NULL
);

CREATE UNLOGGED TABLE forum_users
(
    author CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    slug   CITEXT REFERENCES forums (slug) ON DELETE CASCADE    NOT NULL,
    PRIMARY KEY (slug, author)
);

CREATE INDEX ON forum_users (slug);
CREATE INDEX ON forum_users (author);

create UNLOGGED table threads
(
    id SERIAL primary key NOT NULL,
    author CITEXT REFERENCES users (nickname) ON DELETE CASCADE not null,
    created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum CITEXT REFERENCES forums (slug) ON DELETE CASCADE not null,
    message text not null,
    slug CITEXT,
    title text not null,
    votes int default 0
);

CREATE INDEX ON threads(slug, author);
CREATE INDEX ON threads(forum, created ASC);
CREATE INDEX ON threads(forum, created DESC);
CREATE INDEX ON threads(slug, id);
CREATE INDEX ON threads(id, forum);
CREATE INDEX ON threads(slug, id, forum);

create UNLOGGED table posts
(
    id SERIAL primary key not null,
    author CITEXT REFERENCES users (nickname) ON DELETE CASCADE not null,
    created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum CITEXT REFERENCES forums (slug) ON DELETE CASCADE not null,
    isEdited bool DEFAULT false not null,
    message text not null,
    parent int not null,
    thread int REFERENCES threads (id) ON DELETE CASCADE not null,
    path       INTEGER ARRAY               DEFAULT '{}'              NOT NULL
);

CREATE UNIQUE INDEX ON posts(id, thread);
CREATE UNIQUE INDEX ON posts(id, author);
CREATE INDEX ON posts(thread, path DESC);
CREATE INDEX ON posts(thread, path ASC);
CREATE INDEX ON posts(thread, id DESC);
CREATE INDEX ON posts(thread, id ASC);

create UNLOGGED table votes
(
    nickname CITEXT REFERENCES users (nickname) ON DELETE CASCADE not null,
    thread int REFERENCES threads (id) ON DELETE CASCADE not null,
    vote int,
    PRIMARY KEY (thread, nickname)
);
CREATE UNIQUE INDEX ON votes(thread, nickname);


CREATE FUNCTION update_path_check_parent() RETURNS TRIGGER AS
    $$
DECLARE
temp INT ARRAY;
t integer;
BEGIN
    IF new.parent ISNULL OR new.parent = 0 THEN
        new.path = ARRAY [new.id];
ELSE

SELECT thread
INTO t
FROM posts
WHERE id = new.parent;
IF t ISNULL OR t <> new.thread THEN
            RAISE EXCEPTION 'Not in this thread ID ' USING HINT = 'Please check your parent ID';
END IF;

SELECT path
INTO temp
FROM posts
WHERE id = new.parent;
new.path = array_append(temp, new.id);

END IF;
RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_posts_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_path_check_parent();

CREATE FUNCTION  trigger_posts() RETURNS TRIGGER AS
    $$
BEGIN
UPDATE forums
SET posts = (posts + 1)
WHERE slug = new.forum;
RETURN new;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER triggerPosts
    AFTER INSERT
    ON posts
    for each row
    EXECUTE PROCEDURE trigger_posts();

CREATE FUNCTION  trigger_threads() RETURNS TRIGGER AS
    $$
BEGIN
UPDATE forums
SET threads = (threads + 1)
WHERE slug = new.forum;
RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER triggerThreads
    AFTER INSERT
    ON threads
    for each row
    EXECUTE PROCEDURE trigger_threads();

CREATE FUNCTION  trigger_vote() RETURNS TRIGGER AS
    $$
BEGIN
UPDATE threads
SET votes = (votes + new.vote)
WHERE id = new.thread;
RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER triggerVote
    AFTER INSERT
    ON votes
    for each row
    EXECUTE PROCEDURE trigger_vote();

CREATE FUNCTION  trigger_voteup() RETURNS TRIGGER AS
    $$
BEGIN
UPDATE threads
SET votes = (votes - old.vote + new.vote)
WHERE id = new.thread;
RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER triggerVoteUp
    AFTER UPDATE
    ON votes
    for each row
    EXECUTE PROCEDURE trigger_voteup();

CREATE FUNCTION  insert_forum_users() RETURNS TRIGGER AS
    $$
BEGIN
INSERT INTO forum_users
VALUES (new.author, new.forum)
    ON CONFLICT DO NOTHING;
RETURN NULL;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER update_forum_users_thread
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE insert_forum_users();

CREATE TRIGGER update_forum_users_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE insert_forum_users();



