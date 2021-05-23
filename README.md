
# Steres

Steres is a fast linked tokenized key value store for sequences, with support for replication, multiple machines and multiple drives per machine.

It relies on LevelDB for indexing and filesystem for storage. Uses nginx as the volume server.

Its persistent, and does O(k) read/write. 
Given a sequence with some probablity of noise/error on
the character level, and word/sentence level, it tokenizes 
the sequence as per the given token size, maps that "token +/- error" to 
an identifier(id) for this sequence while linking all tokens of that identifier. 


Basically,

(text, token_size, threshold)   => [Steres]  => (text, id) 



## API

- GET /text/token_size/threshold
    - 302 redirect to nginx volume server

- PUT /text/token_size/threshold
    - 201 = written. 

- DELETE /text/token_size/threshold
     - 201 = written. 