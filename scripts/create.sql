CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users
(
    id SERIAL primary key NOT NULL,
    nickname CITEXT UNIQUE NOT NULL,
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
    userf CITEXT  NOT NULL,
    posts    INTEGER DEFAULT 0                                    NOT NULL,
    threads  INTEGER DEFAULT 0                                    NOT NULL
);

create table threads
(
    id SERIAL primary key NOT NULL,
    author CITEXT not null,
    created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum CITEXT not null,
    message text not null,
    slug CITEXT UNIQUE not null,
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
    forum CITEXT not null,
    isEdited bool DEFAULT false not null,
    message text not null,
    parent int not null,
    thread int not null
);

CREATE UNIQUE INDEX ON posts(id, thread);
CREATE UNIQUE INDEX ON posts(id, author);
CREATE INDEX ON posts(thread, id DESC);
CREATE INDEX ON posts(thread, id ASC);

create table votes
(
    nickname CITEXT  not null,
    thread int not null,
    vote int,
    PRIMARY KEY (thread, nickname)
);

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



