## Change log

1. Modify the returned ID type, it will affect the ID types of GetByID and List for `⓵Create web service based on sql`, which are consistent with the ID types in the database.

> If you are using code for `⓵Create web service based on sql` before v1.8.6, do not modify the sponge version under go.mod and upgrade to v1.8.6 or above. Otherwise, the List interface will return empty data because the original `size` field has become invalid (replaced by `limit` field).
