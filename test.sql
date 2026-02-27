select books.url             as book_url,
       chapter_ranks.url     as chapter_url,
       title,
       chapter_ranks.chapter as chapter,
       author,
       books.summary         as summary,
       ts_headline(
               title || ' ' || chapter_title || ' ' || chapter_text,
               query,
               'StartSel=**,StopSel=**, MaxFragments=1, MinWords=5, MaxWords=100'
       )                     as excerpt,
       rank                  as rank,
       count(*) over () as total_count,
       reason is not NULL    as blacklisted,
       coalesce(reason, '')  as blacklist_reason
from (select distinct on (book_id) book_id,
                                   chapter_title,
                                   chapter_text,
                                   url,
                                   chapter,
                                   ts_rank_cd('{0.014925373, 0.2, 0.4, 1.0}', textsearchable_index_col, query) as rank,
                                   query
      from (select *, chapters.url as chapter_url
            from chapters,
                 websearch_to_tsquery('"water park"') as query
            where textsearchable_index_col @@ query)
     ) as chapter_ranks
         join books on chapter_ranks.book_id = books.id
         left join blacklist on chapter_ranks.book_id = blacklist.book_id
order by rank desc limit 20 offset 0;