update
    client
set
    client_name = ?,
    is_admin = ?,
    last_edited = now()
where
    client_id = ?
