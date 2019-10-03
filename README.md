# local-search

```bash
# load *.md and *.txt from current directory & register document
lsh -i ~/.lsh/lsh_index -d .

# load * md only
lsh -t md -d .

# search `golang` from loaded documents
lsh -q 'golang'
```
