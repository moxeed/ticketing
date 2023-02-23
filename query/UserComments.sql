WITH S
         AS
         (
             SELECT Id, ReplyToId, 1 AS Depth From Comment WHERE UserId = ?
             UNION ALL
             SELECT C.Id, C.ReplyToId, Depth + 1 FROM Comment C
             JOIN S ON S.Id = C.ReplyToId
             WHERE Depth <= 5
         )
SELECT C.ID,
       CreatedAt,
       UpdatedAt,
       DeletedAt,
       C.UserId,
       C.UserName,
       C.Content,
       C.Origin,
       C.IsConfirmed,
       SUM(CASE WHEN Type= 1 THEN 1 ELSE 0 END) AS LikeCount,
       SUM(CASE WHEN Type= 2 THEN 1 ELSE 0 END) AS DislikeCount,
       C.ConfirmUserId,
       C.ConfirmDateTime,
       C.ReplyToId
FROM S
         JOIN Comment C ON C.ID = S.ID OR C.ID = S.ReplyToId
         LEFT JOIN React R ON R.CommentId = C.ID
WHERE C.IsConfirmed = 1
GROUP BY
    C.ID,
    CreatedAt,
    UpdatedAt,
    DeletedAt,
    C.UserId,
    C.UserName,
    C.Content,
    C.Origin,
    C.IsConfirmed,
    C.ConfirmUserId,
    C.ConfirmDateTime,
    C.ReplyToId
