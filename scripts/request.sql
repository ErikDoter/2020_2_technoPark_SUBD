-- /forum/create

insert into forums(slug, title, user)
    value('dota_stories', 'Dota stories', 'ErikDoter');

select posts, slug, title, user, threads from forums
where slug = 'dota_stories';



-- /forum/{dota_stories}/create

insert into threads(author, message, title, forum, slug)
    value('ErikDoter', 'Kak igrat?', 'Kak igrat?', 'dota_stories', 'kak_igrat');

select * from threads where slug = 'kak_igrat';


-- /forum/{dota_stories}/details

select posts, slug, threads, title, user from forums
where slug = 'dota_stories';

-- /forum/{dota_stories}/threads/?limit=1&since=0&desc=0

select * from threads
where forum = 'dota_stories' and created > 0
order by created LIMIT 1;

-- /user/{ErikDoter}/create

insert into users(nickname, fullname, email, about)
value('ErikDoter', 'Erik Nabiev', 'er2000@mail.ru', 'Erik');

-- /user/{ErikDoter}/profile

select about, email, fullname, nickname from users
where nickname = 'ErikDoter';

-- /user/{ErikDoter}/profile (POST)

UPDATE users
set about = 'Erikkkkk',
    email = 'dummy12@mail.ru',
    fullname = 'Erik NeNabiev'
where nickname = 'ErikDoter';


-- /forum/{dota_stories}/users/?limit=1&since=0&desc=0

select T.about, T.email, T.fullname, T.nickname
from (
        SELECT u.about, u.email, u.fullname, u.nickname, u.id
        from forums f join threads t on f.slug = t.forum
        join users u on t.author = u.nickname
        where f.slug = 'dota_stories'
        union
        SELECT uu.about, uu.email, uu.fullname, uu.nickname, uu.id
        from forums ff join threads tt on ff.slug = tt.forum
        join posts pp on pp.thread = tt.id
        join users uu on uu.nickname = pp.author
        where ff.slug = 'dota_stories'
     ) as T
where T.id > 0
ORDER BY lower(T.nickname) LIMIT 1;

-- /thread/{1}/create

insert into posts(author, message, parent, thread, forum)
select 'ErikDoter1', 'Pobedil vhera na Enigme', 0, 1, forum from threads where id = 1;

-- /post/{1}/details

select u.about, u.email, u.fullname, u.nickname
from posts p join users u on (p.id = 1 and p.author = u.nickname);

select f.posts, f.slug, f.threads, f.title, f.user
from posts p join forums f on (p.id = 1 and p.forum = f.slug);

select t.author, t.slug, t.created, t.forum, t.id, t.message, t.title, t.votes
from posts p join threads t on (p.id = 1 and p.thread = t.id);

select author, created, forum, id, message, isEdited, parent, thread
from posts where id = 1;

-- /post/{1}/details

update posts
set message = 'ROSHAN'
where id = 1;

-- /service/clear

truncate table posts ...

-- /thread/{1}/details

select author, created, forum, id, message, slug, title, votes
from threads where id = 1;

-- /thread/{1}/details

update threads
set message = 'Kak pobedit', title = 'Kak'
where id = 1;

-- /threads/{1}/posts

select p.author, p.created, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread
from threads t join posts p on (t.id = 1 and t.id = p.id)
where p.id > 0
order by p.id limit 1;

-- /thread/{1}/vote

insert into votes(nickname, thread, vote)
value('ErikDoter', 1, -1);