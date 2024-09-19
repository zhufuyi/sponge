## Change log

- Modified the code directory structure for large repository services, making the `api` and `third_party` directories shared among all services, while keeping other directories unchanged.
- Added an automatic initialization command for the large repository (mono-repo repository) service mode.
- Optimized the code merging rules, supporting the merging of **gRPC template code**, **Handler code**, **Router configuration**, and **Error code**.
- Added a distributed lock library [dlock](https://github.com/zhufuyi/sponge/tree/main/pkg/dlock) with support for `Redis` and `Etcd`.
- Optimized some `pkg` libraries (**[Scheduled Task Logging](https://github.com/zhufuyi/sponge/issues/62)**, **[ID Generation Functionality](https://github.com/zhufuyi/sponge/blob/main/pkg/krand/README.md#generate-id)**).
