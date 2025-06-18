select
    *
from
    expense
where
    year(expense_date) = ?
    and
    month(expense_date) = ?
    and
    day(expense_date) = ?
