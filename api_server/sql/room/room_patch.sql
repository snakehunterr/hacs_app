update
    room
set
    client_id = ?,
    room_people_count = ?,
    room_area = ?,
    last_edited = now()
where
    room_id = ?
