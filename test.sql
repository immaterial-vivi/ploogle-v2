select book_url,
    chapter_url,
    title,
    chapter,
    author,
    ts_headline(
        title || ' ' || chapter_title || ' ' || chapter_text,
        query,
        'StartSel=**,StopSel=**, MaxFragments=3, MinWords=5, MaxWords=100'
    ) as excerpt,
    rank,
    query,
    count
from (
        select *,
            count(*) over () as count
        from (
                select distinct on (book_id) *
                from (
                        select *,
                            books.url as book_url,
                            ts_rank_cd(textsearchable_index_col, query) as rank
                        from (
                                select *,
                                    chapters.url as chapter_url
                                from chapters,
                                    websearch_to_tsquery($1) as query
                                where textsearchable_index_col @@ query
                            )
                            join books on book_id = books.id
                        order by rank desc
                    )
            )
        limit $2 offset $3
    );