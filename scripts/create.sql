drop database forum;
create database forum;

use forum

create table users
(
    id       int auto_increment primary key,
    nickname varchar(80) not null,
    fullname varchar(80) not null,
    email    varchar(80) not null,
    about    text        not null,
    unique(nickname)
);

create table forums
(
    id int auto_increment primary key,
    slug varchar(80) not null,
    title varchar(80) not null,
    posts int unsigned default 0,
    threads int unsigned default 0,
    user varchar(80) not null,
    unique(slug)
);

create table threads
(
    id int auto_increment primary key,
    author varchar(80) not null,
    created DATETIME DEFAULT NOW() ON UPDATE NOW(),
    forum varchar(80) not null,
    message varchar(255),
    slug varchar(80) not null,
    title varchar(80) not null,
    votes int unsigned default 0,
    unique(slug)
);

create table posts
(
    id int auto_increment primary key,
    author varchar(80) not null,
    created DATETIME DEFAULT NOW() ON UPDATE NOW(),
    forum varchar(80) not null,
    isEdited bool DEFAULT false,
    message text not null,
    parent int not null,
    thread int not null
);

create table votes
(
    id int auto_increment primary key,
    nickname varchar(80) not null,
    thread int not null,
    vote int;
);

CREATE TRIGGER triggerPosts
    AFTER INSERT
    ON posts
    for each row
    update forums f
    set f.posts = f.posts + 1
    where f.slug = new.forum;

CREATE TRIGGER triggerThreads
    AFTER INSERT
    ON threads
    for each row
    update forums f
    set f.threads = f.threads + 1
    where f.slug = new.forum;

CREATE TRIGGER triggerVote
    AFTER INSERT
    ON votes
    for each row
    update threads t
    set t.votes = t.votes + new.vote
    where t.id = new.thread;


