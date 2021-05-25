
# Steres

! In development... 


Steres is a linked tokenized key value store for sequences, with support for replication, multiple machines and multiple drives per machine.

It relies on LevelDB for indexing and filesystem for storage. Uses nginx as the volume server.

Its persistent, and does O(1) read/write. 
Given a sequence with some probablity of noise/error on
the character level, and word/sentence level, it tokenizes 
the sequence as per the given token size, maps that "token +/- error" to 
an identifier(id) for this sequence while linking all tokens of that identifier. 


Basically,
```
(text, token_size, threshold)   => [Steres]  => (text, id) 
```
where, threshold is the minimum ratio of the "hypothesis(input) sequence" to the "stored sequence".

## API


### POST /
```
    data : {
        sequence: "...... .. .. .. .",   (the sequence that you are searching for)
        threshold: 0.94,                 (Similarity threshold ratio, which accounts for noise/error.)
        token_size: 8                    (Depending upon your application of this db)  
    }
```
* 302 redirect to nginx volume server
* returns { id } ( id of the sequence }

---------

### PUT /
```
     data : {
        sequence: "...... .. .. .. .",   (the sequence that you are searching for)
        threshold: 0.94,                 (Similarity threshold ratio, which accounts for noise/error.)
        token_size: 8                    (Depending upon your application of this db)  
    }
   
```
* 201 = written/updated 
* id -> if the sequence is unique (no match) it gets a new id, else return the recognized sequence's id

### DELETE /id
* Blocks. and if 204 => deleted
* deletes the entire sequence for the id.

--------
