select
    *
from
    client
where
    client_name like concat('%', ?, '%');
