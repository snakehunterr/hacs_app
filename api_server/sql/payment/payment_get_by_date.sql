select
    *
from
    payment
where
    year(payment_date) = ?
    and
    month(payment_date) = ?
    and 
    day(payment_date) = ?
