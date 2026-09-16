# Shoebox

Shoebox is a media management tool built to be executed as a single binary. The backoffice is built using stdlib Go and the interface is built using templ + htmx.

I built this project to learn Go as coming from a TypeScript background.

# What is in the repo

- Manaul database migration script
- File indexer
- Checksum dedup logic
- Thumbnail generator using libvips

# Some useful things I learned

1. Upon introducing concurrency in the indexer, the checksum dedup check broke. Reason: one worker was still busy inserting a processed file while other workers were checking for its existence. Solved by introducing a two-phase reserve/complete, mutex-protected dedup with a cleanup on failure.

Commit: ab532de067f9
