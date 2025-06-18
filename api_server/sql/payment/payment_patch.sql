update
    payment
set
    client_id = ?,
    room_id = ?,
    payment_date = ?,
    payment_amount = ?,
    last_edited = now()
where
    payment_id = ?
