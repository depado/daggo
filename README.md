# daggo

Trying dagger with Go and GHA. 

This README will serve as a guide and check list for things I had to figure out 
on my own, missing documentation or things I had to search in discord. 

## Initializing dagger

> [!WARNING]
> These steps are required as to not inject dagger's dependencies in the root 
> `go.mod` file while dagger/dagger#8095 isn't merged. Once merged and released,
> a simple call to `dagger init` will work as intended. 
> Target fixed version: `v0.12.5`

```sh
$ mkdir dagger
$ cd dagger
$ go mod init github.com/depado/daggo/dagger
$ cd ..
$ dagger init --sdk=go --source=dagger/
```

## Migrating from Dockerfile

> [!WARNING]
> As far as I understand, Dockerfile's `ENTRYPOINT` and `CMD` combo instructions
> don't really make sense in OCI. Use only `WithEntrypoint([]string{binPath})`
> with all the expected arguments instead. 
> - [ ] Check if there isn't something missing here.

## Thoughts & notes

> [!WARNING]
> ### Documentation
> 
> Dagger seems to have fundamentaly changed in the past few months (stopped 
> using the CUE language for instance), which mean most online resources 
> (tutorials, articles, and even some on dagger's website/blog) are outdated. 

> [!WARNING]  
> ### Cache
> 
> Dagger pushes its [caching capabilities](https://dagger.io/blog/airbyte-use-case) 
> as a way to greatly speed up build times and indeed, locally it works great. 
> But once you want to use your own CI, things get messy because the only way 
> (for now) to use dagger's caching ability is to use 
> [Dagger Cloud's paid plan](https://dagger.io/pricing) to benefit from their 
> experimental distributed cache. The community is actively trying to figure
> out how to achieve this on GitHub and Discord: dagger/dagger#6911.
> 
> That being said, initiatives are in progress to enable caching to work with 
> other distributed cache providers: dagger/dagger#8004.
> 
> It has also been suggested to take a look at 
> [buildkit-cache-dance](https://github.com/reproducible-containers/buildkit-cache-dance)
> to cache buildx caches with GHA. 

