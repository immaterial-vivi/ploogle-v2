prepare ploogle_websearch_query (tsquery, int, int) as
select $1 as parsed_query,
    title,
    url,
    excerpt,
    author,
    chapter_title,
    chapter,
    total_count
from (
        select distinct on (book_id) books.title as title,
            chapters.url as url,
            ts_headline(
                title || ' ' || chapter_title || ' ' || chapter_text,
                $1
            ) as excerpt,
            books.author as author,
            chapters.chapter_title as chapter_title,
            chapters.chapter as chapter,
            ts_rank_cd(textsearchable_index_col, $1) as rank,
            count(*) over () as total_count
        from chapters
            join books on book_id = books.id
        where textsearchable_index_col @@ $1
    )
order by rank desc
limit $2 offset $3;
"select * from (\n\t\t\tselect distinct on (book_id) query, books.title as title, chapters.url as url, ts_headline(title || ' ' || chapter_title || ' ' || chapter_text, query) as excerpt, books.author as author, chapters.chapter_title as chapter_title, chapters.chapter as chapter, ts_rank_cd(textsearchable_index_col, query) as rank, COUNT(*) OVER () AS total_count\n\t\t\tfrom chapters join books on book_id=books.id, websearch_to_tsquery($1::text) as query\n\t\twhere textsearchable_index_col @@ query)\n\t\torder by rank DESC\n\t\tlimit $2\n\t\toffset $3;"
select books.url,
    title,
    chapter,
    query,
    author,
    ts_headline(
        title || ' ' || chapter_title || ' ' || chapter_text,
        query
    ) as excerpt,
    ts_rank_cd(textsearchable_index_col, query) as rank
from (
        select * (
                select *
                from chapters,
                    websearch_to_tsquery($1) as query
                where textsearchable_index_col @@ query
            )
            join books on book_id = books.id
        order by rank
        limit $2 offset $3
    )