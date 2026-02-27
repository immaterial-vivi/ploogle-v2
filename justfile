[parallel]
dev: svelte go 

svelte: 
   cd site && npm run dev -- --host

go: 
    go run .