# Concepts


Sherlock is a lightweight search library with prefix search built in.  It is intended to be used and configured entirely through struct tags, enabling integration into any app in just a few lines of code. 

At its core is a radix-tree inverted index. Autocomplete is first class citizen. 

It is mostly a toy engine I'm building while I work through the Introduction to Information Retrieval book, but I hope to build it in such a way that it can be embedded in small applications in golang. 

The first step with sherlock is parsing a struct, extracting its tags, and building a schema.  The schema is a reference used to tag Tokens as they're processed during indexing.  For example, if you tag a field with `sherlock:"weight=10,order=asc"`, the indexer will know to tag those hits with the weight before inserting them into the posting list.  Later in the pipeline this information is used to order queries (without interaction by the user). 

# Indexing Pipeline

document, schema, tokens, postings -> inverted index

# Query Pipeline

query, load, search, collect, sort  -> []matches
